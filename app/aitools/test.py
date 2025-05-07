import grpc
import unittest
import sys
from utils.importer import proto_importer
proto_importer()

from protos.vapus_aiutilities.v1alpha1 import vapus_aiutilities_pb2_grpc as pb2_grpc
import protos.vapus_aiutilities.v1alpha1.vapus_aiutilities_pb2 as pb2

import argparse


SERVER_ADDRESS = "localhost:9025"


class TestGenerateEmbedding(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        """Set up the gRPC channel and stub."""
        cls.channel = grpc.insecure_channel(SERVER_ADDRESS)
        cls.stub = pb2_grpc.AIUtilityStub(cls.channel)

    @classmethod
    def tearDownClass(cls):
        """Close the gRPC channel."""
        cls.channel.close()

    def test_generate_embedding(self):
        """Test the GenerateEmbedding method."""
        
        request = pb2.GenerateEmbeddingRequest(
            text=["My name is John Doe.", "This is a test for embedding generation."]
        )

        
        metadata = [("authorization", f"Bearer {JWT_TOKEN}")]

       
        response = self.stub.GenerateEmbedding(request, metadata=metadata)

        
        self.assertTrue(response.embeddings, "Response embeddings are empty.")
        
        for embedding in response.embeddings:
            print(embedding.index)
            self.assertIsInstance(embedding.index, int, "Index is not an integer.")
            self.assertIsInstance(list(embedding.embedding), list, "Embedding is not a list.")
            self.assertTrue(
                all(isinstance(x, float) for x in embedding.embedding),
                "Embedding contains non-float values.",
            )
            self.assertGreater(len(embedding.embedding), 0, "Embedding is empty.")


if __name__ == "__main__":
    
    parser = argparse.ArgumentParser(description="Test GenerateEmbedding gRPC service.")
    parser.add_argument("--jwt", required=True, help="JWT token for authorization.")
    args, remaining_args = parser.parse_known_args()

    
    JWT_TOKEN = args.jwt

    
    sys.argv = [sys.argv[0]] + remaining_args

    
    unittest.main()