import json
import uuid
from pathlib import Path
from datetime import datetime
from typing import Dict, List, Optional
from enum import Enum

QUEUE_FILE = Path(__file__).parent.parent / "tasks.json"


class TaskStatus(Enum):
    QUEUED = "queued"
    IN_PROGRESS = "in_progress"
    COMPLETED = "completed"
    FAILED = "failed"
    CANCELLED = "cancelled"


class TaskQueue:
    def __init__(self):
        self.queue_file = QUEUE_FILE
        self._ensure_queue_file()

    def _ensure_queue_file(self):
        if not self.queue_file.exists():
            self._save_queue({"tasks": []})

    def _load_queue(self) -> Dict:
        with open(self.queue_file) as f:
            return json.load(f)

    def _save_queue(self, data: Dict):
        with open(self.queue_file, 'w') as f:
            json.dump(data, f, indent=2, default=str)

    def add_task(
        self,
        agent_id: str,
        model: str,
        task_description: str,
        requested_by: str,
    ) -> str:
        task_id = f"task-{uuid.uuid4().hex[:8]}"
        task = {
            "id": task_id,
            "agent": agent_id,
            "model": model,
            "task": task_description,
            "status": TaskStatus.QUEUED.value,
            "requested_by": requested_by,
            "created_at": datetime.now().isoformat(),
            "started_at": None,
            "completed_at": None,
            "result": None,
            "error": None,
        }
        data = self._load_queue()
        data['tasks'].append(task)
        self._save_queue(data)
        return task_id

    def get_task(self, task_id: str) -> Optional[Dict]:
        data = self._load_queue()
        for task in data['tasks']:
            if task['id'] == task_id:
                return task
        return None

    def update_task_status(
        self,
        task_id: str,
        status: TaskStatus,
        **kwargs,
    ):
        data = self._load_queue()
        for task in data['tasks']:
            if task['id'] == task_id:
                task['status'] = status.value
                task.update(kwargs)
                break
        self._save_queue(data)

    def get_next_queued_task(self) -> Optional[Dict]:
        data = self._load_queue()
        for task in data['tasks']:
            if task['status'] == TaskStatus.QUEUED.value:
                return task
        return None

    def list_tasks(
        self,
        status: Optional[TaskStatus] = None,
        limit: int = 10,
    ) -> List[Dict]:
        data = self._load_queue()
        tasks = data['tasks']
        if status:
            tasks = [t for t in tasks if t['status'] == status.value]
        tasks.sort(key=lambda t: t['created_at'], reverse=True)
        return tasks[:limit]

