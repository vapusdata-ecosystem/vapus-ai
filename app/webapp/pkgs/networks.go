package pkgs

import (
	"context"
	"log"
	"path/filepath"

	"github.com/rs/zerolog"
	appconfigs "github.com/vapusdata-ecosystem/vapusai/core/app/configs"
	appcl "github.com/vapusdata-ecosystem/vapusai/core/app/grpcclients"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	filetools "github.com/vapusdata-ecosystem/vapusai/core/tools/files"
	utils "github.com/vapusdata-ecosystem/vapusai/webapp/utils"
)

type VapusArtifactStorage struct {
	Spec *models.DataSourceCredsParams `yaml:"spec"`
}

var VapusSvcInternalClientManager *appcl.VapusSvcInternalClients

func InitNetworkConfig(configRoot, fileName string) error {
	DmLogger.Info().Msgf("Reading network configuration with path - %v ", filepath.Join(configRoot, fileName))

	cf, err := dmutils.ReadBasicConfig(filetools.GetConfFileType(fileName), filepath.Join(configRoot, fileName), &appconfigs.NetworkConfig{})
	if err != nil {
		DmLogger.Panic().Err(err).Msg("error while loading service config")
		return err
	}

	svcnetConf, ok := cf.(*appconfigs.NetworkConfig)
	if !ok {
		DmLogger.Panic().Msg("error while loading network config")
		return utils.ErrInvalidNetworkConfig
	} else {
		NetworkConfigManager = svcnetConf
	}
	return nil
}

func InitVapusSvcInternalClients(hostSvc string, logger zerolog.Logger) {
	// TODO: Handle error
	log.Println("here--------------------------------------------")
	err := appcl.SvcUpTimeCheck(context.Background(), NetworkConfigManager, "", logger, 0)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while checking service uptime.")
	} else {
		logger.Info().Msg("service is up and running.")
	}
	log.Println("Services are up & running")
	res, err := appcl.SetupVapusSvcInternalClients(context.Background(), NetworkConfigManager, "", logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while initializing vapus svc internal clients.")
	}
	VapusSvcInternalClientManager = res
}
