"""Discord bot client for agent coordination."""

import asyncio
import sys

import discord
from discord import app_commands
from discord.app_commands import describe

from coordinator.core.config import get_config
from coordinator.registries.agent_registry import (
    load_agents_registry,
    load_models_registry,
    get_agent,
    get_model,
)
from coordinator.tasks.queue import TaskQueue, TaskStatus


@app_commands.command(name="agents", description="List available agents")
async def cmd_agents_list(interaction: discord.Interaction):
    await interaction.response.defer()
    agents = load_agents_registry()

    embed = discord.Embed(title="🤖 Available Agents", color=0x5865F2)
    for agent_id, agent in agents.items():
        capabilities = ", ".join(agent['capabilities'][:3])
        embed.add_field(
            name=f"{agent['name']} ({agent_id})",
            value=(
                f"{agent['description']}\n"
                f"**Default:** {agent['default_model']}\n"
                f"**Can:** {capabilities}"
            ),
            inline=False,
        )
    await interaction.followup.send(embed=embed)


@app_commands.command(name="models", description="List available models")
async def cmd_models_list(interaction: discord.Interaction):
    await interaction.response.defer()
    models = load_models_registry()
    embed = discord.Embed(title="🧠 Available Models", color=0x5865F2)
    for model_id, model in models.items():
        cost = f"${model['cost_per_mtok_input']}/${model['cost_per_mtok_output']} per MTok"
        strengths = ", ".join(model['strengths'][:3])
        embed.add_field(
            name=f"{model_id} ({model['tier']})",
            value=f"**Cost:** {cost}\n**Strengths:** {strengths}",
            inline=False,
        )
    await interaction.followup.send(embed=embed)


@app_commands.command(name="modelinfo", description="Get details about a specific model")
@describe(model="The model ID to get info for")
async def cmd_model_info(interaction: discord.Interaction, model: str):
    models = load_models_registry()
    if model not in models:
        await interaction.response.send_message(f"❌ Model '{model}' not found")
        return
    m = models[model]
    embed = discord.Embed(title=f"📊 {model}", color=0x5865F2)
    embed.add_field(name="Provider", value=m['provider'], inline=True)
    embed.add_field(name="Tier", value=m['tier'], inline=True)
    embed.add_field(
        name="Cost",
        value=f"${m['cost_per_mtok_input']}/${m['cost_per_mtok_output']} per MTok",
        inline=False,
    )
    embed.add_field(name="Context", value=f"{m['context_window']} tokens", inline=True)
    embed.add_field(name="Max Output", value=f"{m['max_tokens']} tokens", inline=True)
    embed.add_field(name="Strengths", value="\n".join(f"✅ {s}" for s in m['strengths']), inline=False)
    embed.add_field(name="Best For", value="\n".join(f"• {b}" for b in m['best_for']), inline=False)
    await interaction.response.send_message(embed=embed)


@app_commands.command(name="launch", description="Launch an agent with a task")
@describe(agent="The agent to launch", model="The AI model to use", task="The task description")
async def cmd_launch(
    interaction: discord.Interaction,
    agent: str,
    model: str,
    task: str,
):
    try:
        agent_config = get_agent(agent)
        model_config = get_model(model)
    except ValueError as e:
        await interaction.response.send_message(f"❌ {e}")
        return

    if model not in agent_config['supported_models']:
        await interaction.response.send_message(
            f"❌ Model '{model}' not supported by agent '{agent}'\n"
            f"Supported: {', '.join(agent_config['supported_models'])}"
        )
        return

    estimated_input_tokens = 5000
    estimated_output_tokens = 2000
    _ = model_config  # kept for potential future cost calc

    queue = TaskQueue()
    task_id = queue.add_task(
        agent_id=agent,
        model=model,
        task_description=task,
        requested_by=interaction.user.name,
    )

    embed = discord.Embed(title="✅ Task Queued", color=0x57F287)
    embed.add_field(name="Task ID", value=task_id, inline=False)
    embed.add_field(name="Agent", value=agent_config['name'], inline=True)
    embed.add_field(name="Model", value=model, inline=True)
    embed.add_field(name="Task", value=task, inline=False)
    await interaction.response.send_message(embed=embed)


@app_commands.command(name="taskstatus", description="Check status of a task")
@describe(task_id="The task ID to check")
async def cmd_task_status(interaction: discord.Interaction, task_id: str):
    queue = TaskQueue()
    task = queue.get_task(task_id)
    if not task:
        await interaction.response.send_message(f"❌ Task '{task_id}' not found")
        return

    status_emoji = {
        "queued": "⏳",
        "in_progress": "🚧",
        "completed": "✅",
        "failed": "❌",
        "cancelled": "🚫",
    }

    embed = discord.Embed(
        title=f"{status_emoji.get(task['status'], '❓')} Task: {task_id}",
        color=0x5865F2,
    )
    embed.add_field(name="Status", value=task['status'].upper(), inline=True)
    embed.add_field(name="Agent", value=task['agent'], inline=True)
    embed.add_field(name="Model", value=task['model'], inline=True)
    embed.add_field(name="Task", value=task['task'], inline=False)
    embed.add_field(name="Created", value=task['created_at'], inline=True)

    if task['started_at']:
        embed.add_field(name="Started", value=task['started_at'], inline=True)
    if task['completed_at']:
        embed.add_field(name="Completed", value=task['completed_at'], inline=True)
    if task['result']:
        embed.add_field(name="Result", value=task['result'], inline=False)
    if task['error']:
        embed.add_field(name="Error", value=task['error'], inline=False)
    await interaction.response.send_message(embed=embed)


@app_commands.command(name="tasklist", description="List recent tasks")
@describe(status="Filter by status (queued, in_progress, completed, failed, cancelled, all)")
async def cmd_task_list(interaction: discord.Interaction, status: str = "all"):
    queue = TaskQueue()
    status_filter = None
    if status != "all":
        try:
            status_filter = TaskStatus(status)
        except ValueError:
            await interaction.response.send_message(
                "❌ Invalid status. Use: queued, in_progress, completed, failed, cancelled, all"
            )
            return

    tasks = queue.list_tasks(status=status_filter, limit=10)
    if not tasks:
        await interaction.response.send_message("No tasks found")
        return

    embed = discord.Embed(title=f"📋 Recent Tasks ({status})", color=0x5865F2)
    for task in tasks[:5]:
        status_emoji = {
            "queued": "⏳",
            "in_progress": "🚧",
            "completed": "✅",
            "failed": "❌",
            "cancelled": "🚫",
        }
        value = (
            f"**Agent:** {task['agent']}\n"
            f"**Model:** {task['model']}\n"
            f"**Status:** {task['status']}"
        )
        embed.add_field(
            name=f"{status_emoji.get(task['status'], '❓')} {task['id']}",
            value=value,
            inline=False,
        )
    await interaction.response.send_message(embed=embed)


class CoordinatorBot(discord.Client):
    def __init__(self):
        intents = discord.Intents.default()
        intents.message_content = True
        intents.guilds = True
        super().__init__(intents=intents)

        self.config = get_config()
        self.tree = app_commands.CommandTree(self)
        self.tree.add_command(cmd_agents_list)
        self.tree.add_command(cmd_models_list)
        self.tree.add_command(cmd_model_info)
        self.tree.add_command(cmd_launch)
        self.tree.add_command(cmd_task_status)
        self.tree.add_command(cmd_task_list)

    async def on_ready(self):
        print(f"✅ Connected as {self.user} (ID: {self.user.id})")
        print(f"📡 Monitoring channel: {self.config.channel_id}")
        for guild in self.guilds:
            print(f"   - {guild.name} (ID: {guild.id})")
        await self.tree.sync()
        print(f"✅ Synced {len(self.tree.get_commands())} slash commands")

    async def send_message(self, content: str) -> bool:
        try:
            channel = self.get_channel(int(self.config.channel_id))
            if channel is None:
                channel = await self.fetch_channel(int(self.config.channel_id))
            if not hasattr(channel, 'send'):
                return False
            await channel.send(content)
            return True
        except Exception:
            return False


_bot: CoordinatorBot | None = None


def get_bot() -> CoordinatorBot:
    global _bot
    if _bot is None:
        _bot = CoordinatorBot()
    return _bot


async def send_announcement(message: str) -> bool:
    bot = CoordinatorBot()
    config = get_config()

    async def send_and_close():
        try:
            await bot.start(config.bot_token)
        except asyncio.CancelledError:
            pass

    async def wait_and_send():
        try:
            await bot.wait_until_ready()
            success = await bot.send_message(message)
            await bot.close()
            return success
        except Exception:
            await bot.close()
            return False

    bot_task = asyncio.create_task(send_and_close())
    send_task = asyncio.create_task(wait_and_send())
    try:
        success = await asyncio.wait_for(send_task, timeout=10.0)
        return success
    except asyncio.TimeoutError:
        return False
    finally:
        bot_task.cancel()
        try:
            await bot_task
        except asyncio.CancelledError:
            pass


def run_bot():
    bot = get_bot()
    config = get_config()
    try:
        print("🤖 Starting Discord Agent Coordinator Bot...")
        print(f"📡 Channel ID: {config.channel_id}")
        bot.run(config.bot_token)
    except KeyboardInterrupt:
        print("\n👋 Shutting down bot...")
        sys.exit(0)
    except Exception as e:
        print(f"❌ Bot error: {e}")
        sys.exit(1)


if __name__ == "__main__":
    run_bot()

