package pkgs

import (
	appconfigs "github.com/vapusdata-ecosystem/vapusdata/core/app/configs"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	filetools "github.com/vapusdata-ecosystem/vapusdata/core/tools/files"
)

var ServiceConfigManager *appconfigs.VapusAISvcConfig
var NetworkConfigManager *appconfigs.NetworkConfig

func newServiceConfig(configRoot, path string) *appconfigs.VapusAISvcConfig {
	return LoadServiceConfig(configRoot, path)
}

func InitServiceConfig(configRoot, path string) {
	if ServiceConfigManager == nil {
		ServiceConfigManager = newServiceConfig(configRoot, path)
	}
}

func LoadServiceConfig(configRoot, path string) *appconfigs.VapusAISvcConfig {
	// Read the service configuration from the file
	DmLogger.Info().Msgf("Reading service configuration with path - %v ", path)

	cf, err := dmutils.ReadBasicConfig(filetools.GetConfFileType(path), path, &appconfigs.VapusAISvcConfig{})
	if err != nil {
		DmLogger.Panic().Err(err).Msg("error while loading service config")
		return nil
	}

	svcConf := cf.(*appconfigs.VapusAISvcConfig)
	svcConf.Path = configRoot
	return svcConf
}

func InitNetworkConfig(configRoot, path string) error {
	DmLogger.Info().Msgf("Reading network configuration with path - %v ", path)

	cf, err := dmutils.ReadBasicConfig(filetools.GetConfFileType(path), path, &appconfigs.NetworkConfig{})
	if err != nil {
		DmLogger.Panic().Err(err).Msg("error while loading service config")
		return err
	}

	svcnetConf, ok := cf.(*appconfigs.NetworkConfig)
	if !ok {
		DmLogger.Panic().Msg("error while loading network config")
		return apperr.ErrInvalidNetworkConfig
	} else {
		NetworkConfigManager = svcnetConf
	}
	return nil
}
