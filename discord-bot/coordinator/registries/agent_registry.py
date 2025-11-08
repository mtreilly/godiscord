import yaml
from pathlib import Path
from typing import Dict, Any

# Go up two levels from coordinator/registries/ to project root, then to registries/
REGISTRIES_DIR = Path(__file__).parent.parent.parent / "registries"


def load_agents_registry() -> Dict[str, Any]:
    with open(REGISTRIES_DIR / "agents.yaml") as f:
        data = yaml.safe_load(f)
    return data['agents']


def load_models_registry() -> Dict[str, Any]:
    with open(REGISTRIES_DIR / "models.yaml") as f:
        data = yaml.safe_load(f)
    return data['models']


def get_agent(agent_id: str) -> Dict[str, Any]:
    agents = load_agents_registry()
    if agent_id not in agents:
        available = "\n  • ".join(agents.keys())
        raise ValueError(
            f"Agent '{agent_id}' not found\n\n"
            f"Available agents:\n  • {available}"
        )
    return agents[agent_id]


def get_model(model_id: str) -> Dict[str, Any]:
    models = load_models_registry()
    if model_id not in models:
        available = "\n  • ".join(models.keys())
        raise ValueError(
            f"Model '{model_id}' not found\n\n"
            f"Available models:\n  • {available}"
        )
    return models[model_id]

