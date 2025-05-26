package apppkgs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"github.com/vapusdata-ecosystem/vapusai/core/data-platform/connectors/databases"
	datasvcpkgs "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	dmlogger "github.com/vapusdata-ecosystem/vapusai/core/pkgs/logger"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

type DmStoreOpts func(*VapusStore)

func WithVapusStoreSecretPath(path string) DmStoreOpts {
	return func(dm *VapusStore) {
		dm.SecretStorePath = path
	}
}

func WithVapusStoreDBPath(path string) DmStoreOpts {
	return func(dm *VapusStore) {
		dm.DBStorePath = path
	}
}

func WithVapusStoreBlobPath(path string) DmStoreOpts {
	return func(dm *VapusStore) {
		dm.BlobStorePath = path
	}
}

func WithVapusStoreArtifactPath(path string) DmStoreOpts {
	return func(dm *VapusStore) {
		dm.ArtifactStorePath = path
	}
}

func WithVapusCacheStorePath(path string) DmStoreOpts {
	return func(dm *VapusStore) {
		dm.CacheStorePath = path
	}
}

type VapusStore struct {
	*SecretStore
	*BeDataStore
	BlobStore          *BlobStore
	Error              error
	ArtifactStoreCreds *models.DataSourceCredsParams
	SecretStorePath    string
	DBStorePath        string
	BlobStorePath      string
	ArtifactStorePath  string
	CacheStorePath     string
}

type BeDataStore struct {
	Db     *databases.DataStoreClient
	PubSub *databases.DataStoreClient
	Cacher *databases.DataStoreClient
	Error  error
}

func NewVapusStore(ctx context.Context, logger zerolog.Logger, opts ...DmStoreOpts) (*VapusStore, error) {
	dmStore := &VapusStore{}
	for _, opt := range opts {
		opt(dmStore)
	}

	dmSec, err := dmStore.NewVapusBESecretStore(ctx, dmStore.SecretStorePath, logger)
	if err != nil {
		return nil, err
	} else {
		logger.Info().Msg("Secret store created successfully")
		dmStore.SecretStore = dmSec
	}

	dmDb, err := dmStore.NewDBStore(ctx, dmSec, logger)
	if err != nil {
		return nil, err
	} else {
		logger.Info().Msg("DB store created successfully")
		dmStore.BeDataStore = dmDb
	}
	err = dmStore.ActivatePostgresExtension(ctx, logger)
	if err != nil {
		return nil, err
	}
	blobOps, err := dmStore.GetBlobStore(ctx, dmStore.BlobStorePath, dmSec, logger)
	if err != nil {
		return nil, err
	} else {
		logger.Info().Msg("Blob store created successfully")
		dmStore.BlobStore = blobOps
	}

	if dmStore.ArtifactStorePath != "" {
		artifactStoreCreds, err := dmStore.NewArtifactStoreCreds(ctx, dmStore.ArtifactStorePath, dmSec, logger)
		if err != nil {
			return nil, err
		} else {
			logger.Info().Msg("Artifact store creds created successfully")
			dmStore.ArtifactStoreCreds = artifactStoreCreds
		}
	}

	return dmStore, nil
}

func (x *VapusStore) NewDBStore(ctx context.Context, secretStore *SecretStore, logger zerolog.Logger) (*BeDataStore, error) {
	bds := &BeDataStore{}
	dbClient, err := x.initDbStores(ctx, x.DBStorePath, secretStore, logger)
	if err != nil {
		logger.Err(err).Msg("error while initializing db data store")
		return nil, err
	}
	bds.Db = dbClient
	cacheClient, err := x.initDbStores(ctx, x.CacheStorePath, secretStore, logger)
	if err != nil {
		logger.Error().Err(err).Msg("error while initializing cache data store")
		return nil, err
	}
	bds.Cacher = cacheClient
	return bds, nil
}

func (x *VapusStore) NewArtifactStoreCreds(ctx context.Context, secName string, secretStore *SecretStore, logger zerolog.Logger) (*models.DataSourceCredsParams, error) {
	artifactCreds := &models.DataSourceCredsParams{}
	origVal, err := secretStore.ReadSecret(ctx, secName)
	if err != nil {
		logger.Err(err).Msg("error while reading artifact store secret data")
		return nil, err
	}

	err = json.Unmarshal([]byte(dmutils.AnyToStr(origVal)), artifactCreds)
	if err != nil {
		logger.Err(err).Msg("error while unmarshalling artifact store secret data")
		return nil, err
	}
	return artifactCreds, nil
}

func (x *VapusStore) initDbStores(ctx context.Context, secName string, secretStore *SecretStore, logger zerolog.Logger) (*databases.DataStoreClient, error) {
	log.Println("Creating db store client with secName: ++++++++++++++++++++++++++++++++++++++ ", secName)
	origVal, err := secretStore.ReadSecret(ctx, secName)
	if err != nil {
		logger.Err(err).Msg("error while reading secret data for data store")
		return nil, err
	}
	log.Println("initDbStores value ------------------->>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>.", dmutils.AnyToStr(origVal))
	creds := &models.DataSourceCredsParams{}
	err = json.Unmarshal([]byte(dmutils.AnyToStr(origVal)), creds)
	if err != nil {
		logger.Err(err).Msg("error while unmarshalling secret data")
		return nil, err
	}
	log.Println("Data source creds", creds)
	return databases.New(ctx, databases.WithInApp(true), databases.WithDataSourceCredsParams(creds), databases.WithLogger(dmlogger.GetSubDMLogger(logger, "PlatformSvc Dbstore", "DataStoreClient")))
}

func (x *VapusStore) GetBlobStore(ctx context.Context, secName string, dmSec *SecretStore, logger zerolog.Logger) (*BlobStore, error) {
	log.Println("Creating GetBlobStore store client with secName: ++++++++++++++++++++++++++++++++++++++ ", secName)
	creds := &models.DataSourceCredsParams{}
	origVal, err := dmSec.ReadSecret(ctx, secName)
	if err != nil {
		logger.Err(err).Msg("error while reading BlobStore store secret data")
		return nil, err
	}
	log.Println("initDbStores value ------------------->>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>.", origVal)
	err = json.Unmarshal([]byte(dmutils.AnyToStr(origVal)), creds)
	if err != nil {
		logger.Err(err).Msg("error while unmarshalling BlobStore store secret data")
		return nil, err
	}
	return NewBlobStoreClient(ctx, creds, logger)
}

func (x *VapusStore) NewVapusBESecretStore(ctx context.Context, secName string, logger zerolog.Logger) (*SecretStore, error) {
	logger.Debug().Msgf("Creating secret store client with secName: %s", secName)
	client, err := NewSecretStoreClient(ctx, logger, WithSecretStorePath(secName))
	if err != nil {
		logger.Info().Msgf("Error while creating secret store client: %v", err)
		return nil, err
	}
	return client, nil
}

func (ds *VapusStore) GetDbStoreParams() *models.DataSourceCredsParams {
	return ds.BeDataStore.Db.DataStoreParams
}

func (ds *VapusStore) ActivatePostgresExtension(ctx context.Context, logger zerolog.Logger) error {
	ddl := "CREATE EXTENSION IF NOT EXISTS vector;"
	logger.Debug().Msgf("DDL: %v", ddl)
	err := ds.Db.RunDDLs(ctx, &ddl)
	if err != nil {
		logger.Err(err).Msg("error while running DDLs for vector extension")
		return err
	}

	// btree_gin used for array, JSONB, and full-text search
	ddl = "CREATE EXTENSION IF NOT EXISTS btree_gin;"
	logger.Debug().Msgf("DDL: %v", ddl)
	err = ds.Db.RunDDLs(ctx, &ddl)
	if err != nil {
		logger.Err(err).Msg("error while running DDLs for btree_gin extension")
		return err
	}
	return nil
}

func (ds *VapusStore) GetPostgresIndexQuery(q *bun.CreateIndexQuery, param *datasvcpkgs.PostgresIndexOpts, logger zerolog.Logger) error {
	mctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if param.IndexAlgo != "" {
		param.IndexAlgo = " " + param.IndexAlgo
	}
	_, err := q.IfNotExists().TableExpr(param.TableName).
		Index(param.Indexname).
		ColumnExpr(fmt.Sprintf("%s%s", param.FieldName, param.IndexAlgo)).
		Using(param.IndexType).
		Exec(mctx)
	if err != nil {
		return err
	}
	return nil
}

func (x *VapusStore) CloseConnection() {
	if x.BeDataStore.Db != nil {
		x.BeDataStore.Db.Close()
	}
	if x.BeDataStore.PubSub != nil {
		x.BeDataStore.PubSub.Close()
	}
	if x.BeDataStore.Cacher != nil {
		x.BeDataStore.Cacher.Close()
	}
	if x.BlobStore != nil {
		x.BlobStore.Close()
	}
	if x.SecretStore != nil {
		x.SecretStore.Close()
	}
}
