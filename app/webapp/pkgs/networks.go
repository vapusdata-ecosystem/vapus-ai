package pkgs

import (
	"path/filepath"

	appconfigs "github.com/vapusdata-ecosystem/vapusai/core/app/configs"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	filetools "github.com/vapusdata-ecosystem/vapusai/core/tools/files"
	utils "github.com/vapusdata-ecosystem/vapusai/webapp/utils"
)

type VapusArtifactStorage struct {
	Spec *models.DataSourceCredsParams `yaml:"spec"`
}

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
