package services

import (
	"github.com/rs/zerolog"
	nabclients "github.com/vapusdata-ecosystem/vapusdata/core/nabrunners/clients"
	dmstores "github.com/vapusdata-ecosystem/vapusdata/platform/datastoreops"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/platform/pkgs"
)

type DMServices struct {
	DMStore *dmstores.DMStore
	*VapusDataServices
	logger zerolog.Logger
	*nabclients.NabRunnerClient
}

var DMServicesManager *DMServices
var helperLogger zerolog.Logger

func newDMServices(dmstore *dmstores.DMStore) *DMServices {
	return &DMServices{
		DMStore: dmstore,
	}
}

func InitDMServices(dmstore *dmstores.DMStore) {
	InitVapusDataServices(dmstore)
	helperLogger = pkgs.GetSubDMLogger(pkgs.SVCS, "helpers")
	if DMServicesManager == nil {
		DMServicesManager = newDMServices(dmstore)
		DMServicesManager.VapusDataServices = VapusDataServiceManager
		DMServicesManager.logger = pkgs.GetSubDMLogger(pkgs.SVCS, "DMServices")
		if dmstores.DMStoreManager.GetCacheStoreParams() != nil {
			DMServicesManager.NabRunnerClient = nabclients.NewAsynqNabRunnerClient(false, nabclients.WithClientRedisOpt(dmstores.DMStoreManager.GetCacheStoreParams().DataSourceCreds))
		} else {
			DMServicesManager.logger.Fatal().Msg("Cache store params not found")
		}
	}
}

// VapusDataServices is a struct that contains the DMStore.
type VapusDataServices struct {
	DMStore *dmstores.DMStore
	Logger  zerolog.Logger
}

// VapusDataServicesManager is the global variable for VapusDataServices struct.
var VapusDataServiceManager *VapusDataServices

// newVapusDataServices creates a new object for VapusDataServices struct.
func newVapusDataServices(dmstore *dmstores.DMStore) *VapusDataServices {
	return &VapusDataServices{
		DMStore: dmstore,
		Logger:  pkgs.GetSubDMLogger(pkgs.SVCS, "VapusDataServices"),
	}
}

// InitVapusDataServices initializes the data marketplace services.
func InitVapusDataServices(dmstore *dmstores.DMStore) {
	if VapusDataServiceManager == nil {
		VapusDataServiceManager = newVapusDataServices(dmstore)
	}
}
