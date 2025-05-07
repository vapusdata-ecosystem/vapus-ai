package services

import (
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusdata/aigateway/pkgs"
	aidmstore "github.com/vapusdata-ecosystem/vapusdata/core/app/datarepo/aistudio"
	"google.golang.org/protobuf/encoding/protojson"
)

type (
	ValidationErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}
)

type AIGatewayServices struct {
	dmstores    *aidmstore.AIStudioDMStore
	logger      zerolog.Logger
	validator   *validator.Validate
	unProtojson *protojson.UnmarshalOptions
	mProtojson  *protojson.MarshalOptions
}

var AIGatewayServicesManager *AIGatewayServices

func NewAIGatewayServices(dmstores *aidmstore.AIStudioDMStore) *AIGatewayServices {
	AIGatewayServicesManager = &AIGatewayServices{
		dmstores:  dmstores,
		logger:    pkgs.DmLogger,
		validator: validator.New(),
		unProtojson: &protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
		mProtojson: &protojson.MarshalOptions{
			UseProtoNames:  true,
			UseEnumNumbers: false,
		},
	}
	return AIGatewayServicesManager
}

func (a *AIGatewayServices) Ready() bool {
	return a.dmstores != nil
}

func (v AIGatewayServices) Validate(data interface{}) []ValidationErrorResponse {
	validationErrors := []ValidationErrorResponse{}

	errs := v.validator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ValidationErrorResponse
			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true
			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}
