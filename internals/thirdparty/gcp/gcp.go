package gcp

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
	dmlogger "github.com/vapusdata-ecosystem/vapusai/core/pkgs/logger"
)

type GcpConfig struct {
	ProjectID, Region, Zone string
	ServiceAccountKey       []byte
}

var (
	GCP_RES_TEMP = "projects/%s"
	GCP_SM_RES   = "secrets/%s"
)

var logger = dmlogger.CoreLogger

func (x *GcpConfig) GetGcpResource(location, resourceName, resourceValue string) string {
	return fmt.Sprintf("projects/%s/locations/%s/%s/%s", x.ProjectID, location, resourceName, resourceValue)
}

func (x *GcpConfig) SetGcpProjectId(logger zerolog.Logger) {
	var res map[string]any
	err := json.Unmarshal(x.ServiceAccountKey, &res)
	if err != nil {
		logger.Err(err).Msg("Error while unmarshalling service account key")
	}
	val, ok := res["project_id"]
	if !ok {
		logger.Err(err).Msg("Error while fetching project id from service account key")
	}
	x.ProjectID, ok = val.(string)
	if !ok {
		logger.Err(err).Msg("Error while type casting project id to string")
	}
}
