package server

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"

	appBooter "github.com/vapusdata-ecosystem/vapusdata/core/app/booter"
	appconfigs "github.com/vapusdata-ecosystem/vapusdata/core/app/configs"
	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusdata/core/pkgs/authn"
	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	dmstores "github.com/vapusdata-ecosystem/vapusdata/platform/datastoreops"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/platform/pkgs"
	services "github.com/vapusdata-ecosystem/vapusdata/platform/services"
)

func packagesInit() {
	//Initialize the logger
	pkgs.InitWAPLogger(debugLogFlag)

	logger = pkgs.GetSubDMLogger(pkgs.IDEN, "vapusPlatform server init")

	logger.Info().Msg("Loading service config for VapusData server")
	// Load the service configuration, secrets inton the memory of the service. These information will be used by the service to connect to the database, vault etc connections
	pkgs.InitServiceConfig(flagconfPath, filepath.Join(flagconfPath, configName))

	logger.Info().Msg("Service config loaded successfully")

	ctx := context.Background()

	bootStores(ctx, pkgs.ServiceConfigManager)
	logger.Info().Msg("Service data stores loaded successfully")
	initStoreDependencies(ctx, pkgs.ServiceConfigManager)

	pkgs.InitNetworkConfig(flagconfPath, filepath.Join(flagconfPath, pkgs.ServiceConfigManager.NetworkConfigFile))

	// Initialize the NewVapusAuth
	//Boot the data stores
	logger.Info().Msgf("Platform Policy Config Path: %s", filepath.Join(flagconfPath, pkgs.ServiceConfigManager.GetPolicyConfPath()))
	err := pkgs.InitPlatformSvcPackages(logger, apppkgs.WithJwtParams(pkgs.JwtParams),
		apppkgs.WithAuthnParams(pkgs.AuthnParams),
		apppkgs.WithPbacConfigPath(filepath.Join(flagconfPath, pkgs.ServiceConfigManager.GetPolicyConfPath())))
	if err != nil {
		if !errors.Is(err, apppkgs.ErrPbacConfigInitFailed) {
			logger.Fatal().Err(err).Msg("error while initializing platform service packages")
		}
	}
	pkgs.InitAuthnManager(pkgs.SvcPackageParams.AuthnParams)
	bootPlatform(ctx, pkgs.ServiceConfigManager.PlatformBaseAccount)
	logger.Info().Msg("Platform booting completed successfully")
	pkgs.InitVapusSvcInternalClients("platform", logger)
	err = pkgs.InitTrinoClient()
	if err != nil {
		logger.Fatal().Err(err).Msg("error while initializing trino client")
	}
	defer ctx.Done()
}

func bootStores(ctx context.Context, conf *appconfigs.PlatformServiceConfig) {
	//Boot the stores
	logger.Info().Msg("Booting the data stores")
	dmstores.InitDMStore(conf)
	if dmstores.DMStoreManager.Error != nil {
		logger.Fatal().Err(dmstores.DMStoreManager.Error).Msg("error while initializing data stores.")
	}
	services.InitDMServices(dmstores.DMStoreManager)
	dmstores.DMStoreManager.ActivatePostgresExtension(ctx, logger)
	appBooter.BootDataTables(ctx, dmstores.DMStoreManager.VapusStore, logger)
}

func initStoreDependencies(ctx context.Context, conf *appconfigs.PlatformServiceConfig) {
	if pkgs.AuthnParams == nil {
		bootAuthn(ctx, conf.GetAuthnSecrets())
	}
	if pkgs.JwtParams == nil {
		bootJwtAuthn(ctx, conf.GetJwtAuthSecretPath())
	}
}

func bootAuthn(ctx context.Context, secName string) {
	logger.Info().Msgf("Boot Authn with secret path: %s", secName)
	secretStr, err := dmstores.DMStoreManager.SecretStore.ReadSecret(ctx, secName)
	if err != nil {
		logger.Fatal().Err(err).Msgf("error while reading authn secret %s", secName)
	}
	tmp := &authn.AuthnSecrets{}
	err = json.Unmarshal([]byte(dmutils.AnyToStr(secretStr)), tmp)
	if err != nil {
		logger.Fatal().Err(err).Msgf("error while unmarshalling authn secret %s", secName)
	}
	pkgs.AuthnParams = tmp
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

func bootPlatform(ctx context.Context, conf *appconfigs.PlatformBootConfig) {
	var err error
	if dmstores.DMStoreManager.Db == nil {
		logger.Fatal().Msg("error while booting VapusData platform")
	}
	logger.Info().Msgf("Platform Boot Config - %v", conf)
	err = appBooter.BootPlatform(ctx, conf, dmstores.DMStoreManager.VapusStore, pkgs.SvcPackageManager, pkgs.SvcPackageParams, logger)

	if err != nil {
		logger.Fatal().Msgf("error while booting platform account. error: %v", err)
	}
	dmstores.InitAccountPool(context.Background(), dmstores.DMStoreManager)
	dmstores.InitPluginPool(context.Background(), dmstores.DMStoreManager)
}
