package pkgs

import (
	"path/filepath"

	appconfigs "github.com/vapusdata-ecosystem/vapusai/core/app/configs"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	filetools "github.com/vapusdata-ecosystem/vapusai/core/tools/files"
)

var WebAppConfigManager *appconfigs.WebAppConfig
var NetworkConfigManager *appconfigs.NetworkConfig

func InitServiceConfig(configRoot, path string) {
	DmLogger.Info().Msgf("WebAppConfigManager current value - %v", WebAppConfigManager)
	if WebAppConfigManager == nil {
		DmLogger.Info().Msg("Initializing the WebAppConfigManager")
		WebAppConfigManager = loadServiceConfig(configRoot, path)
	}
}

func loadServiceConfig(configRoot, fileName string) *appconfigs.WebAppConfig {
	// Read the service configuration from the file
	DmLogger.Info().Msgf("Reading service configuration with path - %v", filepath.Join(configRoot, fileName))
	cf, err := dmutils.ReadBasicConfig(filetools.GetConfFileType(fileName), filepath.Join(configRoot, fileName), &appconfigs.WebAppConfig{})
	if err != nil {
		DmLogger.Panic().Err(err).Msg("error while loading webapp config")
	}

	svcConf := cf.(*appconfigs.WebAppConfig)
	svcConf.Path = configRoot
	return svcConf
}
