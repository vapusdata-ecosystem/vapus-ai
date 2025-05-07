package dmstores

import (
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	appconfigs "github.com/vapusdata-ecosystem/vapusai/core/app/configs"
	aidmstore "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo/aistudio"
)

// GLobal var for DM store, it can accessed across the service
var (
	DMStoreManager *aidmstore.AIStudioDMStore
	logger         = pkgs.GetSubDMLogger(pkgs.DSTORES, "DBStore")
)

// Constructor to create new object for DMStore struct
func newDMStore(conf *appconfigs.VapusAISvcConfig) *aidmstore.AIStudioDMStore {
	logger = pkgs.GetSubDMLogger(pkgs.DSTORES, "DBStore")
	return aidmstore.InitDMStore(conf, logger)
}

func InitDMStore(conf *appconfigs.VapusAISvcConfig) {
	if DMStoreManager == nil || DMStoreManager.SecretStore == nil || DMStoreManager.BeDataStore == nil {
		DMStoreManager = newDMStore(conf)
	}
}
