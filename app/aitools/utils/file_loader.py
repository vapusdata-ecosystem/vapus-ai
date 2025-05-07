from typing import Any
from pathlib import Path
import yaml
import json

def load_basic_config(fileName: str) -> Any:
    path = Path(fileName)
    with path.open() as f:
        return yaml.safe_load(f)