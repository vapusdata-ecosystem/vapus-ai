package server

import (
	"context"
	"flag"
	"os"

	"github.com/rs/zerolog"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	dmcontrollers "github.com/vapusdata-ecosystem/vapusai/aistudio/controllers"
	dmstores "github.com/vapusdata-ecosystem/vapusai/aistudio/datastoreops"
	middlewares "github.com/vapusdata-ecosystem/vapusai/aistudio/middlewares"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	pbtools "github.com/vapusdata-ecosystem/vapusai/core/pkgs/pbtools"
	"github.com/vapusdata-ecosystem/vapusai/core/types"

	interceptors "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	rpcauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	selector "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	grpc "google.golang.org/grpc"
	health "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	debugLogFlag bool
	flagconfPath string
	configName   = "config/aistudio-service-config.yaml"
	networkDir   = "network"
	logger       zerolog.Logger
)

func init() {
	flag.StringVar(&flagconfPath, "conf", "", "config path, eg: --conf=/data/vapusdata")
	flag.BoolVar(&debugLogFlag, "debug", false, "debug loggin, set it to true to enable the debug logs")
	flag.Parse()
	if flagconfPath == "" {
		var ok bool
		flagconfPath, ok = os.LookupEnv(types.SVC_MOUNT_PATH)
		if !ok {
			logger.Fatal().Msgf("SVC_MOUNT_PATH env not found, please set env variable '%v' with dataproduct config to run the product service", types.SVC_MOUNT_PATH)
		}
	}
	logger.Info().Msgf("Config root Path: %s", flagconfPath)
	packagesInit()
}

func initServer(grpcServer *pbtools.GRPCServer) {

	// Setup auth matcher.
	allButTheez := func(ctx context.Context, callMeta interceptors.CallMeta) bool {
		return healthpb.Health_ServiceDesc.ServiceName != callMeta.Service
	}

	// Add unary and stream interceptors for prometheus
	var opts = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			grpc_prometheus.UnaryServerInterceptor,
			selector.UnaryServerInterceptor(rpcauth.UnaryServerInterceptor(middlewares.AuthnMiddleware), selector.MatchFunc(allButTheez)),
			selector.UnaryServerInterceptor(middlewares.UnaryRequestValidator(), selector.MatchFunc(allButTheez)),
		),
		grpc.ChainStreamInterceptor(
			grpc_prometheus.StreamServerInterceptor,
			selector.StreamServerInterceptor(rpcauth.StreamServerInterceptor(middlewares.AuthnMiddleware), selector.MatchFunc(allButTheez)),
		),
	}

	// Create a new GRPC server
	//First step is to configure the vapusPlatform server
	grpcServer.ServerPort = pkgs.NetworkConfigManager.AIStudioSvc.Port

	// Initialize the server
	logger.Info().Msg("Configuring VapusData Platform Server")

	// Initialize the grpc server ops and net listner for the server
	grpcServer.ConfigureGrpcServer(opts, debugLogFlag)
	logger.Info().Msg("VapusData Platform Server configured successfully.")
	if grpcServer.GrpcServ == nil {
		logger.Info().Msg("Failed to initialize VapusData Platform Server")
	}
	healthcheck := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer.GrpcServ, healthcheck)
	// Register the VapusData Platform and Node controller
	reflection.Register(grpcServer.GrpcServ)
	pb.RegisterAIGuardrailsServer(grpcServer.GrpcServ, dmcontrollers.NewVapusAIGuardrails())
	pb.RegisterAIModelsServer(grpcServer.GrpcServ, dmcontrollers.NewAIModels())
	pb.RegisterAIPromptsServer(grpcServer.GrpcServ, dmcontrollers.NewAIPrompts())
	pb.RegisterAIStudioServer(grpcServer.GrpcServ, dmcontrollers.NewAIStudio())
	pb.RegisterDatasourceServiceServer(grpcServer.GrpcServ, dmcontrollers.NewDataSourcesController())
	pb.RegisterObservabilityServiceServer(grpcServer.GrpcServ, dmcontrollers.NewObservabilityController())
	pb.RegisterOrganizationServiceServer(grpcServer.GrpcServ, dmcontrollers.NewOrganizationController())
	pb.RegisterPluginServiceServer(grpcServer.GrpcServ, dmcontrollers.NewPluginController())
	pb.RegisterSecretServiceServer(grpcServer.GrpcServ, dmcontrollers.NewSecretsController())
	pb.RegisterUserManagementServiceServer(grpcServer.GrpcServ, dmcontrollers.NewVapusDataUsers())
	pb.RegisterUtilityServiceServer(grpcServer.GrpcServ, dmcontrollers.NewUtilityController())
	pb.RegisterAgentServiceServer(grpcServer.GrpcServ, dmcontrollers.NewVapusAgentServer())
	pb.RegisterAgentStudioServer(grpcServer.GrpcServ, dmcontrollers.NewVapusAgentStudio())
	pb.RegisterVapusdataServiceServer(grpcServer.GrpcServ, dmcontrollers.NewVapusDataController())
	logger.Info().Msgf("Grpc Server configured at - %v", grpcServer.ServerAddress)
}

func GrpcServer() *pbtools.GRPCServer {
	var grpcServer *pbtools.GRPCServer
	var serverOpts []pbtools.GrpcOptions
	// Initialize the service configuration
	if pkgs.ServiceConfigManager.ServerCerts.Mtls {
		logger.Info().Msg("Configuring VapusData AIStudio Server with MTLS connection")
		serverOpts = append(serverOpts, pbtools.WithMtls(pkgs.ServiceConfigManager.GetMtlsCerts()))
	} else if pkgs.ServiceConfigManager.ServerCerts.PlainTls {
		logger.Info().Msg("Configuring VapusData AIStudio Server with PlainTLS connection")
		serverOpts = append(serverOpts, pbtools.WithPlainTls(pkgs.ServiceConfigManager.GetPlainTlsCerts()))
	} else {
		logger.Info().Msg("Configuring VapusData AIStudio Server with insecure connection")
		serverOpts = append(serverOpts, pbtools.WithInsecure(true))
	}
	grpcServer = pbtools.NewGRPCServer(serverOpts...)
	grpcServer.Logger = logger
	grpcServer.SetUpHttp = false
	initServer(grpcServer)
	return grpcServer
}

func Shutdown(server *pbtools.GRPCServer) {
	server.Close()
	dmstores.DMStoreManager.CloseConnection()
}
