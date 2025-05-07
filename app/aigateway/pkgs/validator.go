package pkgs

import dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"

var RequestValidator *dmutils.DMValidator

func initPlatformRequestValidator() {
	validator, err := dmutils.NewDMValidator()
	if err != nil {
		pkgLogger.Panic().Err(err).Msg("Error while loading validator")
	}
	RequestValidator = validator
}

func GetPlatformRequestValidator() *dmutils.DMValidator {
	if RequestValidator == nil {
		initPlatformRequestValidator()
	}
	return RequestValidator
}
