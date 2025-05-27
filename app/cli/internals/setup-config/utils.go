package setupconfig

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

func GetValidator() *validator.Validate {
	return validator.New()
}

func HandleValiationError(err error, logger zerolog.Logger) error {
	log.Println(err.Error(), "====================================")
	if _, ok := err.(*validator.InvalidValidationError); ok {
		logger.Error().Err(err).Msg("Invalid validation error")
		return err
	}
	if _, ok := err.(*validator.ValidationErrors); ok {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			logger.Error().Err(err).Msgf("Validation error on field: %s, condition: %s", fieldErr.Namespace(), fieldErr.ActualTag())
		}
		return err
	}
	return nil
}
