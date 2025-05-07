package aidmstore

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	models "github.com/vapusdata-ecosystem/vapusai/core/models"
)

func (ds *AIStudioDMStore) DBStoreIntoCredentialModel() *models.GenericCredentialModel {
	if ds.BeDataStore.Db.DataStoreParams.DataSourceCreds == nil {
		return nil
	}

	return &models.GenericCredentialModel{
		ApiToken:     ds.BeDataStore.Db.DataStoreParams.DataSourceCreds.ApiToken,
		ApiTokenType: ds.BeDataStore.Db.DataStoreParams.DataSourceCreds.ApiTokenType,
		Username:     ds.BeDataStore.Db.DataStoreParams.DataSourceCreds.Username,
		Password:     ds.BeDataStore.Db.DataStoreParams.DataSourceCreds.Password,
		AwsCreds:     ds.BeDataStore.Db.DataStoreParams.DataSourceCreds.AwsCreds,
		AzureCreds:   ds.BeDataStore.Db.DataStoreParams.DataSourceCreds.AzureCreds,
		GcpCreds:     ds.BeDataStore.Db.DataStoreParams.DataSourceCreds.GcpCreds,
	}
}

func AlterColumn(ctx context.Context, tablename string, store *apppkgs.VapusStore, logger zerolog.Logger) {
	ddl := fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS scope varchar DEFAULT '';", tablename)
	logger.Debug().Msgf("DDL: %v", ddl)
	err := store.Db.RunDDLs(ctx, &ddl)
	if err != nil {
		logger.Err(err).Msg("error while running DDLs")
		return
	}
	ddl = fmt.Sprintf("UPDATE %s SET scope='PLATFORM_SCOPE' where scope='';", tablename)
	logger.Debug().Msgf("DDL: %v", ddl)
	err = store.Db.RunDDLs(ctx, &ddl)
	if err != nil {
		logger.Err(err).Msg("error while running DDLs")
	}
}
