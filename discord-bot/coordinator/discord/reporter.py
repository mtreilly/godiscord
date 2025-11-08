"""Post task updates to Discord channel."""

import discord
import asyncio
from coordinator.core.config import get_config


class DiscordReporter:
    """Posts task progress to Discord using one-shot connections."""

    def __init__(self):
        self.config = get_config()

    async def _send_message(self, content=None, embed=None):
        """Helper to send a message to the configured channel."""
        intents = discord.Intents.default()
        intents.guilds = True
        client = discord.Client(intents=intents)

        success = False

        @client.event
        async def on_ready():
            nonlocal success
            try:
                channel = client.get_channel(int(self.config.channel_id))
                if channel is None:
                    channel = await client.fetch_channel(int(self.config.channel_id))
                if content:
                    await channel.send(content=content)
                if embed:
                    await channel.send(embed=embed)
                success = True
            finally:
                await client.close()

        try:
            await asyncio.wait_for(client.start(self.config.bot_token), timeout=10.0)
        except asyncio.TimeoutError:
            pass
        except Exception:
            pass

        return success

