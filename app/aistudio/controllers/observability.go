package dmcontrollers

import (
	"github.com/rs/zerolog"

	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"

	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	dmsvc "github.com/vapusdata-ecosystem/vapusai/aistudio/services"
)

type ObservabilityController struct {
	pb.UnimplementedObservabilityServiceServer
	validator  *dmutils.DMValidator
	DMServices *dmsvc.AIStudioServices
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
		DMServices: dmsvc.AIStudioServiceManager,
	}
}

func InitObservabilityController() {
	if ObservabilityControllerManager == nil {
		ObservabilityControllerManager = NewObservabilityController()
	}
}
