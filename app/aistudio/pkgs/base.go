package pkgs

import (
	"github.com/rs/zerolog"
	appdrepo "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	"github.com/vapusdata-ecosystem/vapusai/core/plugins"
	sqlops "github.com/vapusdata-ecosystem/vapusai/core/tools/sqlops"
)

var SvcPackageManager *apppkgs.VapusSvcPackages
var SvcPackageParams *apppkgs.VapusSvcPackageParams
var JwtParams *encryption.JWTAuthn
var GuardrailPoolManager *appdrepo.GuardrailPool
var AIModelNodeConnectionPoolManager *appdrepo.AIModelNodeConnectionPool
var PluginServiceManager *plugins.VapusPlugins
var SqlOps *sqlops.SQLOperator

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

func InitSqlOps() error {
	var err error
	if SqlOps == nil {
		SqlOps, err = sqlops.New(DmLogger)
		if err != nil {
			return err
		}

	}
	return nil
}
