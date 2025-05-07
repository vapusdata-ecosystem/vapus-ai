package server

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"

	dmstores "github.com/vapusdata-ecosystem/vapusai/aistudio/datastoreops"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	services "github.com/vapusdata-ecosystem/vapusai/aistudio/services"
	appconfigs "github.com/vapusdata-ecosystem/vapusai/core/app/configs"
	appdrepo "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

func packagesInit() {
	//Initialize the logger
	pkgs.InitWAPLogger(debugLogFlag)

	logger = pkgs.GetSubDMLogger(pkgs.IDEN, "AIstudio server init")

	logger.Info().Msg("Loading service config for VapusData server")
	// Load the service configuration, secrets inton the memory of the service. These information will be used by the service to connect to the database, vault etc connections
	pkgs.InitServiceConfig(flagconfPath, filepath.Join(flagconfPath, configName))

	pkgs.InitNetworkConfig(flagconfPath, filepath.Join(flagconfPath, pkgs.ServiceConfigManager.NetworkConfigFile))

	logger.Info().Msg("Service config loaded successfully")

	ctx := context.Background()
	bootStores(ctx, pkgs.ServiceConfigManager)
	initStoreDependencies(ctx, pkgs.ServiceConfigManager)
	logger.Info().Msg("Service data stores loaded successfully")

	// dmstores.InitStoreDependencies(ctx, pkgs.ServiceConfigManager)
	logger.Info().Msg("Service store dependencies loaded successfully")

	logger.Info().Msg("Service config loaded successfully")
	// Initialize the jwt authn validator
	logger.Info().Msgf("Loading JWT authn with secret path: %s", pkgs.ServiceConfigManager.GetJwtAuthSecretPath())

	err := pkgs.InitPlatformSvcPackages(logger, apppkgs.WithJwtParams(pkgs.JwtParams))
	if err != nil {
		if !errors.Is(err, apppkgs.ErrPbacConfigInitFailed) {
			logger.Fatal().Err(err).Msg("error while initializing platform service packages")
		}
	}
	pkgs.InitVapusSvcInternalClients()
	err = pkgs.InitTrinoClient()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize Trino client")
	}
	err = pkgs.InitSqlOps()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize SQL operator")
	}
	bootConnectionPool()
	defer ctx.Done()
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
	if pkgs.PluginServiceManager == nil {
		pkgs.PluginServiceManager, _ = appdrepo.NewPluginPool(context.Background(), dmstores.DMStoreManager.VapusStore, logger)
	}

	services.InitAIStudioServices(dmstores.DMStoreManager)
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
