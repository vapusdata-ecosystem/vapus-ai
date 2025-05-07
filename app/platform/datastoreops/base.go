package dmstores

import (
	"context"

	"github.com/rs/zerolog"
	appconfigs "github.com/vapusdata-ecosystem/vapusdata/core/app/configs"
	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/platform/pkgs"
)

type DMStore struct {
	*apppkgs.VapusStore
}

// GLobal var for DM store, it can accessed across the service
var (
	DMStoreManager *DMStore
	logger         zerolog.Logger
)

// Constructor to create new object for DMStore struct
func newDMStore(conf *appconfigs.PlatformServiceConfig) *DMStore {
	ctx := context.Background()
	logger = pkgs.GetSubDMLogger(pkgs.DSTORES, "DBStore")
	vapusStore, err := apppkgs.NewVapusStore(ctx,
		logger,
		apppkgs.WithVapusStoreSecretPath(conf.GetSecretStoragePath()),
		apppkgs.WithVapusStoreDBPath(conf.GetDBStoragePath()),
		apppkgs.WithVapusStoreBlobPath(conf.GetFileStorePath()),
		apppkgs.WithVapusStoreArtifactPath(conf.GetArtifactStorePath()),
		apppkgs.WithVapusCacheStorePath(conf.GetCachStoragePath()),
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while creating VapusStore")
	}
	if vapusStore.ArtifactStoreCreds == nil {
		logger.Fatal().Msg("error while creating VapusStore, artifact store creds is nil")
	} else {
		pkgs.VapusArtifactStorageManager = &pkgs.VapusArtifactStorage{
			Spec: vapusStore.ArtifactStoreCreds,
		}
	}
	return &DMStore{
		VapusStore: vapusStore,
	}
}

// Initializing DMStore struct with object and global var
func InitDMStore(conf *appconfigs.PlatformServiceConfig) {
	logger = pkgs.GetSubDMLogger(pkgs.DSTORES, "DBStore")
	if DMStoreManager == nil || DMStoreManager.SecretStore == nil || DMStoreManager.BeDataStore == nil {
		DMStoreManager = newDMStore(conf)
	}

}

func (ds *DMStore) GetCacheStoreParams() *models.DataSourceCredsParams {
	return ds.BeDataStore.Cacher.DataStoreParams
}
