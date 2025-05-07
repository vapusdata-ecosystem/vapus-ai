package server

import (
	"context"
	"flag"
	"os"

	"github.com/rs/zerolog"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
	pbtools "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/pbtools"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
	dmcontrollers "github.com/vapusdata-ecosystem/vapusdata/platform/controllers"
	dmstores "github.com/vapusdata-ecosystem/vapusdata/platform/datastoreops"
	middlewares "github.com/vapusdata-ecosystem/vapusdata/platform/middlewares"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/platform/pkgs"

	interceptors "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	rpcauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	selector "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpc "google.golang.org/grpc"
	health "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	debugLogFlag bool
	flagconfPath string
	configName   = "config/platform-service-config.yaml"
	bootConfig   = "config/vapusplatform-boot-config.yaml"
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
			selector.UnaryServerInterceptor(rpcauth.UnaryServerInterceptor(middlewares.AuthnMiddleware), selector.MatchFunc(allButTheez)),
			grpc_prometheus.UnaryServerInterceptor,
			selector.UnaryServerInterceptor(middlewares.UnaryRequestValidator(), selector.MatchFunc(allButTheez)),
		),
		grpc.ChainStreamInterceptor(
			selector.StreamServerInterceptor(rpcauth.StreamServerInterceptor(middlewares.AuthnMiddleware), selector.MatchFunc(allButTheez)),
			grpc_prometheus.StreamServerInterceptor,
		),
	}

	// Create a new GRPC server
	//First step is to configure the vapusPlatform server
	grpcServer.ServerPort = pkgs.NetworkConfigManager.PlatformSvc.Port

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
	pb.RegisterVapusdataServiceServer(grpcServer.GrpcServ, dmcontrollers.NewVapusDataController())
	pb.RegisterUserManagementServiceServer(grpcServer.GrpcServ, dmcontrollers.NewVapusDataUsers())
	pb.RegisterOrganizationServiceServer(grpcServer.GrpcServ, dmcontrollers.NewOrganizationController())
	pb.RegisterPluginServiceServer(grpcServer.GrpcServ, dmcontrollers.NewPluginController())
	pb.RegisterUtilityServiceServer(grpcServer.GrpcServ, dmcontrollers.NewUtilityController())
	pb.RegisterDatasourceServiceServer(grpcServer.GrpcServ, dmcontrollers.NewDataSourcesController())
	pb.RegisterObservabilityServiceServer(grpcServer.GrpcServ, dmcontrollers.NewObservabilityController())
	pb.RegisterSecretServiceServer(grpcServer.GrpcServ, dmcontrollers.NewSecretsController())

	if grpcServer.SetUpHttp {
		grpcServer.HttpGwPort = int32(pkgs.NetworkConfigManager.PlatformSvc.HttpGwPort)
		grpcServer.MuxInstance = runtime.NewServeMux(
			runtime.WithMetadata(middlewares.HttpAuthnMiddleware),
			runtime.WithMarshalerOption("application/protobuf", &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
					UseEnumNumbers:  false,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			}),
		)
		if err := pb.RegisterUserManagementServiceHandlerServer(context.Background(), grpcServer.MuxInstance, dmcontrollers.NewVapusDataUsers()); err != nil {
			logger.Fatal().Err(err).Msg("Failed to register the Vapus Data Product Server handler")
		}
		if err := pb.RegisterOrganizationServiceHandlerServer(context.Background(), grpcServer.MuxInstance, dmcontrollers.NewOrganizationController()); err != nil {
			logger.Fatal().Err(err).Msg("Failed to register the Vapus Data Product Server handler")
		}
		logger.Info().Msgf("Http Server configured at - %v", grpcServer.Http1Address)
	}

	logger.Info().Msgf("Grpc Server configured at - %v", grpcServer.ServerAddress)

}

func GrpcServer() *pbtools.GRPCServer {
	var grpcServer *pbtools.GRPCServer
	var serverOpts []pbtools.GrpcOptions
	// Initialize the service configuration
	if pkgs.ServiceConfigManager.ServerCerts.Mtls {
		logger.Info().Msg("Configuring VapusData Platform Server with MTLS connection")
		serverOpts = append(serverOpts, pbtools.WithMtls(pkgs.ServiceConfigManager.GetMtlsCerts()))
	} else if pkgs.ServiceConfigManager.ServerCerts.PlainTls {
		logger.Info().Msg("Configuring VapusData Platform Server with PlainTLS connection")
		serverOpts = append(serverOpts, pbtools.WithPlainTls(pkgs.ServiceConfigManager.GetPlainTlsCerts()))
	} else {
		logger.Info().Msg("Configuring VapusData Platform Server with insecure connection")
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
