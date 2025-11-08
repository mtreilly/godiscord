"""Configuration loader for agent coordinator."""

import os
from pathlib import Path
from typing import Any

import yaml
from dotenv import load_dotenv


class Config:
    """Configuration manager for Discord agent coordinator."""

    def __init__(self, config_path: str | None = None):
        """Load configuration from YAML file and environment variables.

        Args:
            config_path: Path to config.yaml. If None, searches for it in:
                         1. Current directory
                         2. Parent of coordinator module
        """
        # Load environment variables
        load_dotenv()

        # Find config file
        if config_path is None:
            config_path = self._find_config()

        # Load YAML config
        with open(config_path) as f:
            self._config = yaml.safe_load(f)

        # Override with environment variables
        self._apply_env_overrides()

    def _find_config(self) -> str:
        """Find config.yaml in standard locations."""
        # Try current working directory
        cwd_config = Path.cwd() / "config.yaml"
        if cwd_config.exists():
            return str(cwd_config)

        # Try parent of coordinator module (discord-bot/)
        module_dir = Path(__file__).parent.parent.parent
        module_config = module_dir / "config.yaml"
        if module_config.exists():
            return str(module_config)

        raise FileNotFoundError(
            "config.yaml not found. Searched:\n"
            f"  - {cwd_config}\n"
            f"  - {module_config}"
        )

    def _apply_env_overrides(self):
        """Override config values with environment variables."""
        # Discord token (required)
        bot_token = os.getenv("DISCORD_BOT_TOKEN")
        if not bot_token:
            raise ValueError(
                "DISCORD_BOT_TOKEN not set. Copy .env.example to .env and fill it in."
            )
        # Strip quotes if present
        bot_token = bot_token.strip("'\"")
        self._config["discord"]["bot_token"] = bot_token

        # Channel ID (can be set via env)
        channel_id = os.getenv("DISCORD_CHANNEL_ID")
        if channel_id:
            self._config["discord"]["channel_id"] = channel_id

        # Agent name (can be overridden)
        agent_name = os.getenv("AGENT_NAME")
        if agent_name:
            self._config["agent"]["name"] = agent_name

    def get(self, key: str, default: Any = None) -> Any:
        """Get config value using dot notation (e.g., 'discord.channel_id')."""
        keys = key.split(".")
        value = self._config
        for k in keys:
            if isinstance(value, dict):
                value = value.get(k)
            else:
                return default
        return value if value is not None else default

    @property
    def bot_token(self) -> str:
        """Discord bot token."""
        return self.get("discord.bot_token")

    @property
    def channel_id(self) -> str:
        """Discord channel ID for announcements."""
        channel_id = self.get("discord.channel_id")
        if not channel_id:
            raise ValueError(
                "discord.channel_id not set in config.yaml. "
                "Add the channel ID after creating the Discord channel."
            )
        return str(channel_id)

    @property
    def agent_name(self) -> str:
        """Default agent name."""
        return self.get("agent.name", "Agent")

    @property
    def emojis(self) -> dict[str, str]:
        """Emoji mapping for message types."""
        return self.get("formatting.emojis", {})


# Singleton instance
_config: Config | None = None


def get_config() -> Config:
    """Get or create global config instance."""
    global _config
    if _config is None:
        _config = Config()
    return _config

