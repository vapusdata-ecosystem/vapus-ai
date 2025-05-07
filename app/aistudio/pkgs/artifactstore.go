package pkgs

import (
	models "github.com/vapusdata-ecosystem/vapusdata/core/models"
)

type VapusArtifactStorage struct {
	Spec *models.DataSourceCredsParams `yaml:"spec"`
}

var VapusArtifactStorageManager *VapusArtifactStorage
