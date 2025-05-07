package pkgs

import (
	pbac "github.com/vapusdata-ecosystem/vapusai/core/pkgs/pbac"
)

var PlatformRBACManager *pbac.PbacConfig

func InitPolicyLib(configPath string) {
	if PlatformRBACManager == nil {
		val, err := pbac.LoadPbacConfig(configPath)
		if err != nil {
			pkgLogger.Err(err).Msg("error initializing and loading of pbac config")
			panic(err)
		}
		PlatformRBACManager = val
	}
}
