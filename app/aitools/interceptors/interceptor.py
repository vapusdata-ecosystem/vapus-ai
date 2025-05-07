# import grpc
# from typing import Callable, Any
# from helpers.jwtvalidate import validateJWT
# import json
# from grpc_interceptor import ServerInterceptor
# from grpc_interceptor.exceptions import GrpcException
import grpc
import re
import jwt
from typing import Any
from helpers import settings, logger

_AUTH_HEADER_KEY = "authorization"

# class Interceptor(ServerInterceptor):
#     """
#     A gRPC interceptor that processes requests through multiple functions before handling the actual RPC.
#     """
    
#     def intercept(self, method: Callable, request: Any, context: grpc.ServicerContext, method_name: str) -> Any:
#         try:
#             print("Interceptor: Intercepting request")
#             metadata = context.invocation_metadata()
#             print("Interceptor: Metadata received:", metadata)
#             jwt = self.extract_jwt(metadata)
#             print("Interceptor: JWT extracted:", jwt)
#             scope = validateJWT(jwt)
#             print("Interceptor: JWT validated, scope:", scope)
#             context.scope = scope  
#         except grpc.RpcError as e:
#             context.abort(e.code(), e.details())
#         except Exception:
#             context.abort(grpc.StatusCode.UNAUTHENTICATED, "Unable to validate JWT")
#         print("Interceptor: Calling the actual method")
#         return method(request, context)

#     def extract_jwt(self, metadata):
#         """
#         Extracts the JWT token from the metadata.
#         """
        
#         for key, value in metadata:
#             if key.lower() == "authorization":
#                 if value.startswith("Bearer "):
#                     return value[len("Bearer "):]
#                 else:
#                     raise grpc.RpcError(
#                         grpc.StatusCode.UNAUTHENTICATED,
#                         "Invalid Authorization header format",
#                     )
#         raise grpc.RpcError(
#             grpc.StatusCode.UNAUTHENTICATED,
#             "Authorization metadata is missing",
#         )

class JwtValidationInterceptor(grpc.ServerInterceptor):
    def __init__(self):
        def abort(ignored_request, context):
            context.abort(grpc.StatusCode.UNAUTHENTICATED, self._abort_handler_message)

        self._abort_handler = grpc.unary_unary_rpc_method_handler(abort)
        self._abort_handler_message: str = "Invalid auth token"
    
    def extract_jwt_bearer_token(self, metadata) -> str:
        token = None
        for key, value in metadata:
            print(f"Key: {key}, Value: {value}")
            if key == _AUTH_HEADER_KEY:
                pattern = r"Bearer\s+([a-zA-Z0-9\-._~+/]+[=]*)"
                match = re.search(pattern, value)
                if match:
                    token = match.group(1)
                    return token
        return token
    
    def validate_jwt_token(self, token: str) -> Any:
        # validate and decode the token
        try:
            # decoded = jwt.decode(token, settings.JWT_PRIV_KEY, algorithms=["ES512"])
            decoded = jwt.decode(token, options={"verify_signature": False})
        except Exception as error:
            logger.error("Token is invalid with error - {error}", error=error)
            return None
        return decoded

    def intercept_service(self, continuation, handler_call_details):
        # Example _HandlerCallDetails(method='/vapusdata.aiplane.v1.VapusAiService/GetAvailableLlms',
        #  invocation_metadata=(_Metadatum(key='user-agent', value='grpc-node-js/1.9.14-postman.1'),
        #  _Metadatum(key='authorization', value='Bearer likyugkijyhg'),
        #  _Metadatum(key='accept-encoding', value='identity')))
        meta_datas = dict(handler_call_details.invocation_metadata)
        token = self.extract_jwt_bearer_token(meta_datas.items())
        if token is not None:
            valid_decoded = self.validate_jwt_token(token)
            if valid_decoded is not None:
                context = handler_call_details.invocation_metadata + (('token_claim', valid_decoded),)
                handler_call_details = handler_call_details._replace(invocation_metadata=context)
                return continuation(handler_call_details)
            else:
                self._abort_handler_message = "Invalid auth token"
                return self._abort_handler
        else:
            self._abort_handler_message = "No auth token provided"
            return self._abort_handler