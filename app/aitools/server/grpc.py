from concurrent import futures
from loguru import logger
import sys
# import os
# from grpc_reflection.v1alpha import reflection
from utils.importer import proto_importer
from interceptors.interceptor import JwtValidationInterceptor
from database.connector import DatabaseConnector
import argparse
proto_importer()


import grpc
import protos.vapus_aiutilities.v1alpha1.vapus_aiutilities_pb2 as pb2
from protos.vapus_aiutilities.v1alpha1 import vapus_aiutilities_pb2_grpc

from helpers.config import load_vapusaiserver_config,VapusAiConfig
# from helpers.secrets import init_vapus_backend_secrets, SecretStore
from helpers.logger import *
from helpers import settings
# from server.boot import ServerBoot
from controller.vapusmlutilities import AIUtilityService

class GrpcServer:
    """
    Represents a gRPC server for the VapusAi service.
    """

    serviceConfig: VapusAiConfig
    # secretsConfig: SecretStore
    
    @classmethod
    def configure_logger(cls,args: argparse.Namespace):
        """
        Configures the logger based on the command line arguments.

        Args:
            args (list): The command line arguments.
        """
        logger.remove()
        try:
            debug = args.debug
        except:
            debug = False
        if debug:
            logger.add(sys.stderr, format="{time} {level} {message}", level="DEBUG")
            config = {
            "handlers": [
                {"sink": sys.stdout, "level": "DEBUG"},
                ]
            }
        else:
            logger.add(sys.stderr, format="{time} {level} {message}", level="INFO")
            config = {
            "handlers": [
                {"sink": sys.stdout, "level": "INFO"},
            ]
        }
        logger.configure(**config)
        logger.info("Logger configured")
        service_logger = logger
        
    @classmethod
    def init_server(cls):
        """
        Initializes and starts the gRPC server.

        This method initializes a gRPC server, adds the VapusAiServiceServicer to the server,
        starts the server on the specified port, and waits for termination.

        Args:
            cls: The class object.

        Returns:
            None
        """
        port = cls.serviceConfig.networkConfig.aiutility.port
        service_logger.info("Starting server on port {port}", port=port)
        try:
            # Add these server options
            server_options = [
                ('grpc.so_reuseport', 1),
                ('grpc.max_connection_idle_ms', 30000),
                ('grpc.max_concurrent_streams', 100),
                ('grpc.http2.min_time_between_pings_ms', 10000),
                ('grpc.http2.max_ping_strikes', 0),
                ('grpc.http2.max_pings_without_data', 0)
            ]
            
            cls.server = grpc.server(
                futures.ThreadPoolExecutor(max_workers=10),
                interceptors=[JwtValidationInterceptor()],
                options=server_options
            )
            
            vapus_aiutilities_pb2_grpc.add_AIUtilityServicer_to_server(
                AIUtilityService(), 
                cls.server
            )
            
            for address in [f'0.0.0.0:{port}', f'[::]:{port}', f'localhost:{port}']:
                try:
                    cls.server.add_insecure_port(address)
                    service_logger.info(f"Successfully bound to {address}")
                    break
                except Exception as e:
                    service_logger.warning(f"Could not bind to {address}: {str(e)}")
            
        except Exception as e:
            service_logger.error("Error initializing grpc server: {error}", error=str(e))
            raise e
        
    def start(cls):
        """
        Starts the gRPC server.
        """
        cls.configure()
        cls.init_server()
        egnine = cls.startDBEngine()
        # ServerBoot.boot(service_logger)

        try:
            cls.server.start()
            service_logger.info("Server started on port {port}", port=str(cls.serviceConfig.networkConfig.aiutility.port))
            cls.server.wait_for_termination()
        except Exception as e:
            service_logger.error("Error starting grpc server: {error}", error=str(e))
            raise e
        
    @classmethod  
    def startDBEngine(cls):
        db_connector = DatabaseConnector()
        dbSecret = cls.serviceConfig.mainConfig.vapusBESecretStorage.secret
        try:
            engine = db_connector.NewConnection(dbSecret)
        except Exception as e:
            raise(e)
        return engine
    
    @classmethod
    def configure(cls):
        """
        Starts the gRPC server by loading the service configuration, initializing secrets, configuring the logger, and initializing the server.
        """ 
        parser = argparse.ArgumentParser(description="Configure and start the gRPC server.")
        parser.add_argument("--conf", required=True, help="Path to the service configuration file.")
        # parser.add_argument("--debug", action="store_true", help="Enable debug logging.")
        args = parser.parse_args()
        config_path = args.conf
        cls.serviceConfig = load_vapusaiserver_config(config_path)
        settings.set_service_config(cls.serviceConfig)
        service_logger.info("Loaded service config")
        #cls.secretsConfig = init_vapus_backend_secrets(cls.serviceConfig.mainConfig.vapusBEDbStore.path, cls.serviceConfig.mainConfig.vapusBESecretStore.path)
        #settings.set_secret_store(cls.secretsConfig)
        service_logger.info("Loaded secrets")
        cls.configure_logger(args)
        
    
