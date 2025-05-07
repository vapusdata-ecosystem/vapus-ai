import jwt
from models.aiutility_models import scope
from pydantic import ValidationError


def validateJWT(token :str):

    
    try:
        payload = jwt.decode(token,options={"verify_signature":False,"verify_exp":True},algorithms=["ES512"])
    except Exception as e :
        raise Exception("Invalid JWT")
    
    
    if payload.get("exp") is None or payload.get("scope") is None:
        raise Exception("Missing fields")
    
    try:
        scope.model_validate(payload.get("scope"))
    except ValidationError as e:
        raise e
    
    return payload.get("scope")



