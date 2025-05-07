package dmstores

import (
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	appconfigs "github.com/vapusdata-ecosystem/vapusai/core/app/configs"
	aigwdmstore "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo/aistudio"
)

// GLobal var for DM store, it can accessed across the service
var (
	DMStoreManager *aigwdmstore.AIStudioDMStore
	logger         = pkgs.GetSubDMLogger(pkgs.DSTORES, "DBStore")
)

// Constructor to create new object for DMStore struct
func newDMStore(conf *appconfigs.VapusAISvcConfig) *aigwdmstore.AIStudioDMStore {
	logger = pkgs.GetSubDMLogger(pkgs.DSTORES, "DBStore")
	return aigwdmstore.InitDMStore(conf, logger)
}

func InitDMStore(conf *appconfigs.VapusAISvcConfig) {
	if DMStoreManager == nil || DMStoreManager.SecretStore == nil || DMStoreManager.BeDataStore == nil {
		DMStoreManager = newDMStore(conf)
	}
}
