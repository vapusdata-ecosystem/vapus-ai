package pkgs

import (
	"github.com/rs/zerolog"
	appconfigs "github.com/vapusdata-ecosystem/vapusai/core/app/configs"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
)

var SvcPackageManager *apppkgs.VapusSvcPackages
var SvcPackageParams *apppkgs.VapusSvcPackageParams

func InitSvcPackageParams() {
	SvcPackageParams = &apppkgs.VapusSvcPackageParams{}
}

func InitPlatformSvcPackages(logger zerolog.Logger, opts ...apppkgs.VapusSvcPkgOpts) error {
	if SvcPackageParams == nil {
		SvcPackageParams = &apppkgs.VapusSvcPackageParams{}
	}
	if SvcPackageManager == nil {
		SvcPackageManager = &apppkgs.VapusSvcPackages{}
	}
	for _, opt := range opts {
		opt(SvcPackageParams)
	}
	SvcPackageManager.ValidEnums = appconfigs.GetValidEnums()
	return nil
}
