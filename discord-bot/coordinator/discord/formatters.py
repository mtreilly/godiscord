"""Message formatting utilities for Discord announcements."""

from datetime import datetime


def format_task_announcement(
    agent_name: str,
    title: str,
    files: list[str] | None = None,
    notes: str | None = None,
    emoji: str = "🚧",
    show_timestamp: bool = True,
    code_blocks: bool = True,
    max_files: int = 10,
) -> str:
    """Format a task announcement message."""
    parts = [f"{emoji} **{agent_name}** started: {title}"]

    if files:
        if len(files) > max_files:
            display_files = files[:max_files]
            remaining = len(files) - max_files
            files_text = "\n".join(display_files)
            if code_blocks:
                parts.append(f"```\n{files_text}\n... and {remaining} more files\n```")
            else:
                parts.append(f"{files_text}\n... and {remaining} more files")
        else:
            files_text = "\n".join(files)
            if code_blocks:
                parts.append(f"```\n{files_text}\n```")
            else:
                parts.append(files_text)

    if notes:
        parts.append(f"_{notes}_")

    if show_timestamp:
        timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        parts.append(f"_Started at {timestamp}_")

    return "\n".join(parts)


def format_task_completion(
    agent_name: str,
    title: str,
    duration: str | None = None,
    commit: str | None = None,
    notes: str | None = None,
    emoji: str = "✅",
    show_timestamp: bool = True,
) -> str:
    """Format a task completion message."""
    parts = [f"{emoji} **{agent_name}** completed: {title}"]

    if duration:
        parts.append(f"Duration: {duration}")
    if commit:
        parts.append(f"Commit: `{commit}`")
    if notes:
        parts.append(f"_{notes}_")

    if show_timestamp:
        timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        parts.append(f"_Completed at {timestamp}_")

    return "\n".join(parts)


def format_task_update(
    agent_name: str,
    message: str,
    emoji: str = "🔄",
    show_timestamp: bool = True,
) -> str:
    """Format a task update message."""
    parts = [f"{emoji} **{agent_name}**: {message}"]
    if show_timestamp:
        timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        parts.append(f"_{timestamp}_")
    return "\n".join(parts)


def format_question(
    agent_name: str,
    question: str,
    context: str | None = None,
    emoji: str = "❓",
    show_timestamp: bool = True,
) -> str:
    """Format a question/approval request message."""
    parts = [f"{emoji} **{agent_name}** asks: {question}"]
    if context:
        parts.append(f"```\n{context}\n```")
    if show_timestamp:
        timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        parts.append(f"_{timestamp}_")
    return "\n".join(parts)


def format_error(
    agent_name: str,
    error_message: str,
    emoji: str = "❌",
    show_timestamp: bool = True,
) -> str:
    """Format an error message."""
    parts = [f"{emoji} **{agent_name}** error: {error_message}"]
    if show_timestamp:
        timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        parts.append(f"_{timestamp}_")
    return "\n".join(parts)

