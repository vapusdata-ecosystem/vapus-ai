from pydantic import BaseModel
from enum import Enum
from typing import List

class GenerateEmbeddingRequest(BaseModel):
    text: List[str]

class actionEnum(str,Enum):
    INVALID_ACTION = "INVALID_ACTION"
    ANALYZE = "ANALYZE"
    ACT = "ACT"

class postDetectActionsEnum(str,Enum):
    INVALID_PDA = "INVALID_PDA"
    REDACT = "REDACT"
    FAKE = "FAKE"
    EMPTY = "EMPTY"

class SensitivityAnalyzerRequest(BaseModel):
    text: List[str]
    entities: List[str]
    action: actionEnum
    postDetectAction: postDetectActionsEnum

class SummarizerRequest(BaseModel):
    text: bytes
    
class scope(BaseModel):
    userId: str
    accountId: str
    domainId: str
    domainRole: str = None
    roleScope: str = None
    platformRole: str = None
