package server

import (
	"bytes"
	"context"
	"reflect"

	// "encoding/json"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytedance/sonic"
	json "github.com/bytedance/sonic"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	dmstores "github.com/vapusdata-ecosystem/vapusai/aigateway/datastoreops"
	"github.com/vapusdata-ecosystem/vapusai/aigateway/middlewares"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aigateway/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/aigateway/services"
	appconfigs "github.com/vapusdata-ecosystem/vapusai/core/app/configs"
	appdrepo "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	"github.com/vapusdata-ecosystem/vapusai/core/pkgs/pbtools"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

var debugLogFlag bool
var configName = "config/aistudio-service-config.yaml"
var logger zerolog.Logger
var flagconfPath string

func loadConfPath() {
	flag.StringVar(&flagconfPath, "conf", "", "config path, eg: -conf /data/vapusdata")
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
}

type NDJSONMarshaler struct{}

// ContentType returns the content type based on the value being marshaled
func (m *NDJSONMarshaler) ContentType(v interface{}) string {
	return "application/x-ndjson"
}

func (m *NDJSONMarshaler) Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.ConfigDefault.NewEncoder(&buf)

	switch val := v.(type) {
	case proto.Message:
		// Handle single protobuf message
		data, err := pbtools.ProtoJsonMarshal(val)
		if err != nil {
			return nil, err
		}
		buf.Write(data)
		buf.WriteByte('\n')

	case []proto.Message:
		// Handle slice of protobuf messages
		for _, msg := range val {
			data, err := pbtools.ProtoJsonMarshal(msg)
			if err != nil {
				return nil, err
			}
			buf.Write(data)
			buf.WriteByte('\n')
		}

	case []interface{}:
		// Handle generic slices
		for _, item := range val {
			if err := enc.Encode(item); err != nil {
				return nil, err
			}
		}

	default:
		rv := reflect.ValueOf(v)
		log.Println("protobuf message default rv", rv)
		// Handle single value
		if err := enc.Encode(v); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (m *NDJSONMarshaler) Unmarshal(data []byte, v interface{}) error {
	switch val := v.(type) {
	case proto.Message:
		obj, ok := val.(proto.Message)
		if !ok {
			return errors.New("invalid type for proto message")
		}
		return pbtools.ProtoJsonUnMarshal(data, obj)
	case []proto.Message:
		// Handle slice of protobuf messages
		for _, msg := range val {
			obj, ok := msg.(proto.Message)
			if !ok {
				return errors.New("invalid type for proto message")
			}
			if err := pbtools.ProtoJsonUnMarshal(data, obj); err != nil {
				return err
			}
		}
	case []interface{}:
		// Handle generic slices
		for _, item := range val {
			if err := sonic.Unmarshal(data, item); err != nil {
				return err
			}
		}
	default:
		// Handle single value
		if err := sonic.Unmarshal(data, v); err != nil {
			return err
		}
	}
	return sonic.Unmarshal(data, v)
}

func (m *NDJSONMarshaler) NewDecoder(r io.Reader) runtime.Decoder {
	return json.ConfigDefault.NewDecoder(r)
}

func (m *NDJSONMarshaler) NewEncoder(w io.Writer) runtime.Encoder {
	return json.ConfigDefault.NewEncoder(w)
}

// Initialize the echo server for webapp
func init() {
	// INitialize the logger
	ctx := context.Background()
	pkgs.InitWAPLogger(debugLogFlag)

	logger = pkgs.GetSubDMLogger(pkgs.IDEN, "VapusData platform server init")

	logger.Info().Msg("Logger middleware Initialized Successfully")

	loadConfPath()
	// Initialize the webapp configuration
	pkgs.InitServiceConfig(flagconfPath, filepath.Join(flagconfPath, configName))

	pkgs.InitNetworkConfig(flagconfPath, filepath.Join(flagconfPath, pkgs.ServiceConfigManager.NetworkConfigFile))
	bootStores(ctx, pkgs.ServiceConfigManager)
	initStoreDependencies(ctx, pkgs.ServiceConfigManager)
	err := pkgs.InitPlatformSvcPackages(logger, apppkgs.WithJwtParams(pkgs.JwtParams))
	if err != nil {
		if !errors.Is(err, apppkgs.ErrPbacConfigInitFailed) {
			logger.Fatal().Err(err).Msg("error while initializing platform service packages")
		}
	}

	// CHeck if the webapp configuration is initialized
	if pkgs.ServiceConfigManager == nil {
		logger.Fatal().Msg("Failed to initialize the WebAppConfigManager")
	}

	// CHeck if the webapp authentication configuration is initialized
	if pkgs.ServiceConfigManager.GetJwtAuthSecretPath() == types.EMPTYSTR {
		logger.Fatal().Msg("Failed to initialize the AuthnSecrets")
	}

	pkgs.InitVapusSvcInternalClients()
	bootConnectionPool()
}

func GRPCGatewayProxy(handler http.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Preserve the original Accept header
		accept := c.Get("Accept")
		// Convert Fiber request to HTTP request
		req, err := http.NewRequest(
			c.Method(),
			c.OriginalURL(),
			bytes.NewReader(c.Body()),
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Copy headers
		for k, v := range c.GetReqHeaders() {
			req.Header[k] = v
		}

		// Force the Accept header if needed
		if accept != "" {
			req.Header.Set("Accept", accept)
		}
		// Create response recorder
		rec := httptest.NewRecorder()

		// Serve via gRPC gateway
		handler.ServeHTTP(rec, req)

		// Convert response to Fiber
		resp := rec.Result()
		defer resp.Body.Close()
		// Copy headers
		for k, v := range resp.Header {
			c.Set(k, strings.Join(v, ","))
		}
		// Return the response
		return c.Status(resp.StatusCode).SendStream(resp.Body)
	}
}

func NewAIGateway() *fiber.App {
	var err error
	app := fiber.New(fiber.Config{
		// Prefork:           true,
		CaseSensitive:     true,
		StrictRouting:     true,
		AppName:           "Vapus AIGateway",
		StreamRequestBody: true,
		ServerHeader:      "Vapus AIGateway",
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
		ColorScheme:       fiber.DefaultColors,
		BodyLimit:         100 * 1024 * 1024,
		Concurrency:       256 * 1024,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Aimodelnode",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS, HEAD",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: false,
	}))
	app.Use(middlewares.LoggingMiddleware())
	gwmux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(func(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
			ct := string(w.Header().Get("Transfer-Encoding"))
			if strings.Contains(ct, "chunked") {
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
			}
			return nil
		}),
		// runtime.WithForwardResponseOption(runtime.ForwardResponseStream),
		runtime.WithOutgoingHeaderMatcher(runtime.DefaultHeaderMatcher),
		runtime.WithIncomingHeaderMatcher(runtime.DefaultHeaderMatcher),
		runtime.WithMarshalerOption("application/x-ndjson", &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   false,
				EmitUnpopulated: true,
				UseEnumNumbers:  false,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
		runtime.WithMarshalerOption("application/ndjson", &NDJSONMarshaler{}),
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
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024*1024),
			grpc.MaxCallSendMsgSize(1024*1024*1024),
		),
	}
	err = pb.RegisterVapusdataServiceHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering VapusdataService handler from endpoint")
	}

	err = pb.RegisterOrganizationServiceHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering OrganizationService handler from endpoint")
	}

	err = pb.RegisterDatasourceServiceHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering DataProduct service handler from endpoint")
	}

	err = pb.RegisterPluginServiceHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering PluginService handler from endpoint")
	}

	err = pb.RegisterUserManagementServiceHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering UserManagementService handler from endpoint")
	}

	err = pb.RegisterUtilityServiceHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering UtilityService handler from endpoint")
	}

	err = pb.RegisterAIGuardrailsHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering AIAgents handler from endpoint")
	}
	err = pb.RegisterAgentServiceHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering Vapus Agents Service from endpoint")
	}
	err = pb.RegisterAgentStudioHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering Vapus Agents Studio from endpoint")
	}
	err = pb.RegisterAIStudioHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering AIStudio handler from endpoint")
	}
	err = pb.RegisterAIModelsHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering AIModels handler from endpoint")
	}
	err = pb.RegisterAIPromptsHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering AIPrompts handler from endpoint")
	}
	err = pb.RegisterAIGuardrailsHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering AIGuardrails handler from endpoint")
	}

	err = pb.RegisterSecretServiceHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering SecretService handler from endpoint")
	}

	// app.All("/*", func(c *fiber.Ctx) error {
	// 	fasthttpadaptor.NewFastHTTPHandler(gwmux)(c.Context())
	// 	return nil
	// })
	// app.All("/api/v1alpha1/*", adaptor.HTTPHandler(gwmux))
	// app.All("/api/v1alpha1/*", GRPCGatewayProxy(gwmux))
	app.Use(recover.New())
	streamHandler := fasthttpadaptor.NewFastHTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gwmux.ServeHTTP(w, r)
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	})
	app.Use(middlewares.Authentication)
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			return services.AIGatewayServicesManager.Ready()
		},
		ReadinessEndpoint: "/ready",
	}))
	app.Get("/metrics", monitor.New(monitor.Config{Title: "AI Gateway Metrics"}))
	gtw := app.Group("/gateway")
	chatRouter(gtw)
	// app.All("/api/v1alpha1/*", adaptor.HTTPHandler(gwmux))
	app.All("/api/v1alpha1/*", func(c *fiber.Ctx) error {
		// log.Println("=============================================== 1111", string(c.Request().Header.Header()))
		// // fasthttpadaptor.NewFastHTTPHandler(gwmux)(c.Context())
		// log.Println("===============================================  2222", string(c.Response().Header.Header()))
		// ct := string(c.Response().Header.Peek(fiber.HeaderContentType))
		// log.Println("=============================================== 3333", ct)
		// streamHandler(c.Context())
		// if strings.Contains(ct, "application/x-ndjson") {

		// 	c.Set("Cache-Control", "no-cache")
		// 	c.Set("Connection", "keep-alive")
		// }
		streamHandler(c.Context())
		return nil
	})

	return app
}

func Shutdown() {
	dmstores.DMStoreManager.CloseConnection()
}

func initStoreDependencies(ctx context.Context, conf *appconfigs.VapusAISvcConfig) {
	if pkgs.JwtParams == nil {
		bootJwtAuthn(ctx, conf.GetJwtAuthSecretPath())
	}
}

func bootStores(ctx context.Context, conf *appconfigs.VapusAISvcConfig) {
	//Boot the stores
	logger.Info().Msg("Booting the data stores")
	dmstores.InitDMStore(conf)
	if dmstores.DMStoreManager.Error != nil {
		logger.Fatal().Err(dmstores.DMStoreManager.Error).Msg("error while initializing data stores.")
	}
	services.NewAIGatewayServices(dmstores.DMStoreManager)
}

func bootConnectionPool() {
	// Boot the connection pool
	logger.Info().Msg("Booting the connection pool")
	resp := appdrepo.InitAIModelNodeConnectionPool(pkgs.AIModelNodeConnectionPoolManager, appdrepo.WithMpLogger(logger),
		appdrepo.WithMpDMStore(dmstores.DMStoreManager.VapusStore))
	if resp != nil {
		pkgs.AIModelNodeConnectionPoolManager = resp
	}
	logger.Info().Msg("Connection pool booted successfully")
	gdResp, err := appdrepo.InitGuardrailPool(context.Background(), pkgs.AIModelNodeConnectionPoolManager, appdrepo.WithGpLogger(logger),
		appdrepo.WithGpStore(dmstores.DMStoreManager.VapusStore))
	if err != nil {
		logger.Fatal().Err(err).Msg("error while booting guardrail pool")
	}
	if gdResp != nil {
		pkgs.GuardrailPoolManager = gdResp
	}
	logger.Info().Msg("Guardrail pool booted successfully")
}

func bootJwtAuthn(ctx context.Context, secName string) {
	logger.Info().Msgf("Boot Jwt Authn with secret path: %s", secName)
	secretStr, err := dmstores.DMStoreManager.SecretStore.ReadSecret(ctx, secName)
	if err != nil {
		logger.Fatal().Err(err).Msgf("error while reading Jwt secret %s", secName)
	}
	tmp := &encryption.JWTAuthn{}
	err = json.Unmarshal([]byte(dmutils.AnyToStr(secretStr)), tmp)
	if err != nil {
		logger.Fatal().Err(err).Msgf("error while unmarshalling Jwt secret %s", secName)
	}
	pkgs.JwtParams = tmp
}
