package aidmstore

import (
	"context"

	"github.com/rs/zerolog"
	appconfigs "github.com/vapusdata-ecosystem/vapusdata/core/app/configs"
	appdrepo "github.com/vapusdata-ecosystem/vapusdata/core/app/datarepo"
	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	models "github.com/vapusdata-ecosystem/vapusdata/core/models"
	dmlogger "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/logger"
)

type AIStudioDMStore struct {
	*apppkgs.VapusStore
	logger  zerolog.Logger
	Account *models.Account
}

// Constructor to create new object for DMStore struct
func newDMStore(conf *appconfigs.VapusAISvcConfig, logger zerolog.Logger) *AIStudioDMStore {
	ctx := context.Background()
	vapusStore, err := apppkgs.NewVapusStore(ctx,
		logger,
		apppkgs.WithVapusStoreSecretPath(conf.GetSecretStoragePath()),
		apppkgs.WithVapusStoreDBPath(conf.GetDBStoragePath()),
		apppkgs.WithVapusStoreBlobPath(conf.GetFileStorePath()),
		apppkgs.WithVapusCacheStorePath(conf.GetCachStoragePath()),
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while creating VapusStore")
	}
	var account *models.Account
	res := appdrepo.BootAccountCache(vapusStore, logger)
	if res == nil {
		logger.Fatal().Msg("error while booting account cache")
	}
	for _, acc := range res {
		account = acc
		break
	}
	return &AIStudioDMStore{
		VapusStore: vapusStore,
		logger:     logger,
		Account:    account,
	}
}

func InitDMStore(conf *appconfigs.VapusAISvcConfig, logger zerolog.Logger) *AIStudioDMStore {
	logger = dmlogger.GetSubDMLogger(logger, "DBStore", "AIStudio")
	return newDMStore(conf, logger)
}

func (ds *AIStudioDMStore) GetDbStoreParams() *models.DataSourceCredsParams {
	return ds.BeDataStore.Db.DataStoreParams
}
