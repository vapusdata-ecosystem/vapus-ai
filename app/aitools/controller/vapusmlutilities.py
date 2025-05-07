from utils.importer import proto_importer
from google.protobuf.json_format import MessageToDict
import json

proto_importer()

import grpc
from protos.vapus_aiutilities.v1alpha1 import vapus_aiutilities_pb2_grpc as aiutilities
import protos.vapus_aiutilities.v1alpha1.vapus_aiutilities_pb2 as pb2
from models import aiutility_models
from enum import Enum
from services.utilities import Utilities
class AIUtilityService(aiutilities.AIUtilityServicer):

    utilities = Utilities()
    
    def GenerateEmbedding(self, request, context):        
        print("Generate Embedding",request)
        try:
            aiutility_models.GenerateEmbeddingRequest.model_validate(MessageToDict(request))
        except Exception as e:
            raise e        
        
        return self.utilities.GenerateEmbedding(request,context)
        
    def SensitivityAnalyzer(self, request, context):

        try:
            aiutility_models.SensitivityAnalyzerRequest.model_validate(MessageToDict(request))
        except Exception as e:
            raise e
        
        return self.utilities.SensitivityAnalyzer(request,context)

    def Summarizer(self, request, context):

        try:
            aiutility_models.SummarizerRequest.model_validate(MessageToDict(request))
        except Exception as e:
            raise e
        
        return self.utilities.SummarizeText(request,context)





        








