package pkgs

import (
	"github.com/rs/zerolog"
	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
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
