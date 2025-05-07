package pkgs

import (
	"github.com/rs/zerolog"
	appdrepo "github.com/vapusdata-ecosystem/vapusdata/core/app/datarepo"
	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
	"github.com/vapusdata-ecosystem/vapusdata/core/plugins"
)

var SvcPackageManager *apppkgs.VapusSvcPackages
var SvcPackageParams *apppkgs.VapusSvcPackageParams
var JwtParams *encryption.JWTAuthn
var GuardrailPoolManager *appdrepo.GuardrailPool
var AIModelNodeConnectionPoolManager *appdrepo.AIModelNodeConnectionPool
var PluginServiceManager *plugins.VapusPlugins

func InitSvcPackageParams() {
	SvcPackageParams = &apppkgs.VapusSvcPackageParams{}
}

func InitPlatformSvcPackages(logger zerolog.Logger, opts ...apppkgs.VapusSvcPkgOpts) error {
	if SvcPackageParams == nil {
		SvcPackageParams = &apppkgs.VapusSvcPackageParams{}
	}
	for _, opt := range opts {
		opt(SvcPackageParams)
	}

	var err error
	SvcPackageParams, SvcPackageManager, err = apppkgs.InitSvcPackages(SvcPackageParams, SvcPackageManager, logger, opts...)
	if err != nil {
		return err
	}
	return nil
}
