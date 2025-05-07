package dmstores

import (
	"context"
	"fmt"

	models "github.com/vapusdata-ecosystem/vapusdata/core/models"
)

func (ds *DMStore) GetDbStoreParams() *models.DataSourceCredsParams {
	return ds.BeDataStore.Db.DataStoreParams
}

func (ds *DMStore) DBStoreIntoCredentialModel() *models.GenericCredentialModel {
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

// func (ds *DMStore) CreateIndex(ctx context.Context, index string) error {
// 	if exists, err := ds.BeDataStore.Db.ElasticSearchClient.TClient.Indices.Exists(index).Do(ctx); exists || err != nil {
// 		logger.Debug().Msgf("Index %v already exists.", index)
// 		return err
// 	}
// 	_, err := ds.BeDataStore.Db.ElasticSearchClient.TClient.Indices.Create(index).Do(ctx)
// 	if err != nil {
// 		logger.Fatal().Err(err).Ctx(ctx).Msg("error while creating index in elastic search")
// 		return err
// 	}
// 	logger.Debug().Msgf("Index %v created successfully.", index)
// 	return nil
// }

// func (ds *DMStore) CreateIndexWithMapping(ctx context.Context, index string, mapping *types.TypeMapping) error {
// 	if exists, err := ds.BeDataStore.Db.ElasticSearchClient.TClient.Indices.Exists(index).Do(ctx); exists || err != nil {
// 		logger.Debug().Msgf("Index %v already exists.", index)
// 		return err
// 	}
// 	_, err := ds.BeDataStore.Db.ElasticSearchClient.TClient.Indices.Create(index).Mappings(mapping).Do(ctx)
// 	if err != nil {
// 		logger.Fatal().Err(err).Msg("error while updating index with mapping in elastic search, call error")
// 		return err
// 	}
// 	logger.Debug().Msgf("Index %v created successfully.", index)
// 	return nil
// }

func AlterColumn(ctx context.Context, tablename string) {
	ddl := fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS scope varchar DEFAULT '';", tablename)
	logger.Debug().Msgf("DDL: %v", ddl)
	err := DMStoreManager.Db.RunDDLs(ctx, &ddl)
	if err != nil {
		logger.Err(err).Msg("error while running DDLs")
		return
	}
	ddl = fmt.Sprintf("UPDATE %s SET scope='PLATFORM_SCOPE' where scope='';", tablename)
	logger.Debug().Msgf("DDL: %v", ddl)
	err = DMStoreManager.Db.RunDDLs(ctx, &ddl)
	if err != nil {
		logger.Err(err).Msg("error while running DDLs")
	}
}

// func (ds *DMStore) UpdateArns(ctx context.Context) error {
// 	resMap := map[string]string{
// 		vapussvcs.DataProductsTable:           mpb.Resources_DATAPRODUCTS.String(),
// 		vapussvcs.DataSourcesTable:            mpb.Resources_DATASOURCES.String(),
// 		vapussvcs.DataWorkerDeploymentsTable:  mpb.Resources_DATA_WORKER_DEPLOYMENTS.String(),
// 		vapussvcs.DataWorkersTable:            mpb.Resources_DATAWORKERS.String(),
// 		vapussvcs.DataProductDeploymentsTable: mpb.Resources_DATA_CONTAINER_DEPLOYMENTS.String(),
// 		vapussvcs.OrganizationsTable:                mpb.Resources_DOMAINS.String(),
// 		vapussvcs.AIModelPromptTable:          mpb.Resources_AIPROMPTS.String(),
// 		vapussvcs.AIModelsNodesTable:          mpb.Resources_AIMODEL_NODES.String(),
// 		vapussvcs.VapusAIAgentsTable:          mpb.Resources_AIAGENTS.String(),
// 		vapussvcs.AccountsTable:               mpb.Resources_ACCOUNT.String(),
// 		vapussvcs.UsersTable:                  "USER",
// 		vapussvcs.PluginsTable:                mpb.Resources_PLUGINS.String(),
// 	}
// 	for _, table := range []string{
// 		vapussvcs.DataProductsTable,
// 		vapussvcs.DataSourcesTable,
// 		vapussvcs.DataWorkerDeploymentsTable,
// 		vapussvcs.DataWorkersTable,
// 		vapussvcs.DataProductDeploymentsTable,
// 		vapussvcs.OrganizationsTable,
// 		vapussvcs.AIModelPromptTable,
// 		vapussvcs.AIModelsNodesTable,
// 		vapussvcs.VapusAIAgentsTable,
// 		vapussvcs.AccountsTable,
// 		vapussvcs.UsersTable,
// 		vapussvcs.PluginsTable,
// 	} {
// 		objs := []map[string]interface{}{}
// 		// _, err := ds.BeDataStore.Db.PostgresClient.DB.NewSelect().
// 		// 	Model(&objs).
// 		// 	// ModelTableExpr(vapussvcs.DataProductsTable).
// 		// 	ModelTableExpr("rental_team_catalog").
// 		// 	Where("deleted_at IS NULL").
// 		// 	Exec(ctx)
// 		query := fmt.Sprintf("SELECT * FROM %s WHERE deleted_at IS NULL", table)
// 		log.Println("query", query)
// 		err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &objs)
// 		if err != nil {
// 			logger.Err(err).Msgf("error while loading all data from datastore for table %v", table)
// 			return err
// 		}
// 		for _, obj := range objs {
// 			log.Println(obj["editors"])
// 			ww := obj["editors"].(string)
// 			ww = strings.ReplaceAll(ww, "{", "")
// 			ww = strings.ReplaceAll(ww, "}", "")
// 			wls := strings.Split(ww, ",")
// 			dmn := ""
// 			organizationAb, ok := obj["organization"].(string)
// 			if !ok {
// 				logger.Err(err).Msg("error while getting organization from object")
// 				dmn = ""
// 			} else {
// 				dmn = organizationAb
// 			}

// 			err = apppkgss.AddResourceArn(ctx, ds.BeDataStore.Db, &models.VapusResourceArn{
// 				ResourceName: resMap[table],
// 				ResourceId:   obj["vapus_id"].(string),
// 				VapusBase: models.VapusBase{
// 					Editors: wls,
// 				},
// 			}, logger, map[string]string{
// 				encryption.ClaimOrganizationKey:  dmn,
// 				encryption.ClaimUserIdKey:  obj["created_by"].(string),
// 				encryption.ClaimAccountKey: obj["owner_account"].(string),
// 			})
// 			if err != nil {
// 				logger.Err(err).Msgf("error while updating arn for table %v", table)
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }
