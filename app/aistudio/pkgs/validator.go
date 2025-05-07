package pkgs

import dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"

var platformRequestValidator *dmutils.DMValidator

func initPlatformRequestValidator() {
	validator, err := dmutils.NewDMValidator()
	if err != nil {
		pkgLogger.Panic().Err(err).Msg("Error while loading validator")
	}
	platformRequestValidator = validator
}

func GetPlatformRequestValidator() *dmutils.DMValidator {
	if platformRequestValidator == nil {
		initPlatformRequestValidator()
	}
	return platformRequestValidator
}
