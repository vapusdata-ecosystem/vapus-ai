package dmcontrollers

import (
	"github.com/rs/zerolog"

	pb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"

	pkgs "github.com/vapusdata-ecosystem/vapusdata/platform/pkgs"
	dmsvc "github.com/vapusdata-ecosystem/vapusdata/platform/services"
)

type ObservabilityController struct {
	pb.UnimplementedObservabilityServiceServer
	validator  *dmutils.DMValidator
	DMServices *dmsvc.DMServices
	Logger     zerolog.Logger
}

var ObservabilityControllerManager *ObservabilityController

func NewObservabilityController() *ObservabilityController {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "ObservabilityController")
	validator, err := dmutils.NewDMValidator()
	if err != nil {
		l.Panic().Err(err).Msg("Error while loading validator")
	}

	l.Info().Msg("ObservabilityController Controller initialized")
	return &ObservabilityController{
		validator:  validator,
		Logger:     l,
		DMServices: dmsvc.DMServicesManager,
	}
}

func InitObservabilityController() {
	if ObservabilityControllerManager == nil {
		ObservabilityControllerManager = NewObservabilityController()
	}
}
