package pkgs

import (
	"context"

	"github.com/rs/zerolog"
	appcl "github.com/vapusdata-ecosystem/vapusdata/core/app/grpcclients"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
)

type VapusArtifactStorage struct {
	Spec *models.DataSourceCredsParams `yaml:"spec"`
}

var VapusArtifactStorageManager *VapusArtifactStorage

var VapusSvcInternalClientManager *appcl.VapusSvcInternalClients

func InitVapusSvcInternalClients(hostSvc string, logger zerolog.Logger) {
	// TODO: Handle error
	res, err := appcl.SetupVapusSvcInternalClients(context.Background(), NetworkConfigManager, NetworkConfigManager.PlatformSvc.ServiceName, logger)
	if err != nil {
		logger.Err(err).Msg("error while initializing vapus svc internal clients.")
	}
	VapusSvcInternalClientManager = res
}
