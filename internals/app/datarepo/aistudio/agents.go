package aidmstore

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
)

// CreateDataSource creates a new data source in the data Organization store
func (ds *AIStudioDMStore) CreateVapusAgent(ctx context.Context, obj *models.VapusAgents, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(apppkgs.VapusAgentsTable).Exec(ctx)
	if err != nil {
		ds.logger.Err(err).Ctx(ctx).Msg("error while saving fabric agent in datastore")
		return err
	}
	return nil
}

func (ds *AIStudioDMStore) PutVapusAgent(ctx context.Context, obj *models.VapusAgents, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewUpdate().Model(obj).ModelTableExpr(apppkgs.VapusAgentsTable).Where("vapus_id = ?", obj.VapusID).Exec(ctx)
	if err != nil {
		ds.logger.Err(err).Ctx(ctx).Msg("error while fabric agent in datastore")
		return err
	}
	return nil
}
func (ds *AIStudioDMStore) GetVapusAgent(ctx context.Context, agentId string, ctxClaim map[string]string) (*models.VapusAgents, error) {
	result := []*models.VapusAgents{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE vapus_id='%s' AND created_by='%s' AND Organization='%s' ORDER BY ID DESC",
		apppkgs.VapusAgentsTable, agentId, ctxClaim[encryption.ClaimUserIdKey], ctxClaim[encryption.ClaimOrganizationKey])
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		ds.logger.Err(err).Ctx(ctx).Msg("error while getting fabric agent from datastore")
		return nil, err
	}
	return result[0], err
}

func (ds *AIStudioDMStore) ListVapusAgents(ctx context.Context, condition string, ctxClaim map[string]string) ([]*models.VapusAgents, error) {
	result := []*models.VapusAgents{}
	query := ""
	if condition != "" {
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s",
			apppkgs.VapusAgentsTable, condition)
	} else {
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s AND created_by='%s' AND Organization='%s' ORDER BY ID DESC",
			apppkgs.VapusAgentsTable, condition, ctxClaim[encryption.ClaimUserIdKey], ctxClaim[encryption.ClaimOrganizationKey])
	}

	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		ds.logger.Err(err).Ctx(ctx).Msg("error while getting fabric agent from datastore")
		return nil, err
	}
	return result, err
}

func CreateVapusAgentLog(ctx context.Context, dmstore *apppkgs.VapusStore, obj *models.VapusAgentLog, logger zerolog.Logger, ctxClaim map[string]string) error {
	_, err := dmstore.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(apppkgs.VapusAgentLogTable).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving fabric agent log in datastore")
		return err
	}
	return nil
}

func AddVapusAgent(ctx context.Context, dmstore *apppkgs.VapusStore, obj *models.VapusAgents, logger zerolog.Logger, ctxClaim map[string]string) error {
	_, err := dmstore.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(apppkgs.VapusAgentsTable).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving fabric agent in datastore")
		return err
	}
	return nil
}

func PutVapusAgent(ctx context.Context, dmstore *apppkgs.VapusStore, obj *models.VapusAgents, logger zerolog.Logger, ctxClaim map[string]string) error {
	_, err := dmstore.Db.PostgresClient.DB.NewUpdate().Model(obj).ModelTableExpr(apppkgs.VapusAgentsTable).Where("vapus_id = ?", obj.VapusID).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while fabric agent in datastore")
		return err
	}
	return nil
}
func GetVapusAgent(ctx context.Context, dmstore *apppkgs.VapusStore, agentId string, logger zerolog.Logger, ctxClaim map[string]string) (*models.VapusAgents, error) {
	result := []*models.VapusAgents{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE vapus_id='%s' AND created_by='%s' AND Organization='%s' ORDER BY ID DESC",
		apppkgs.VapusAgentsTable, agentId, ctxClaim[encryption.ClaimUserIdKey], ctxClaim[encryption.ClaimOrganizationKey])
	err := dmstore.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msg("error while getting fabric agent from datastore")
		return nil, err
	}
	return result[0], err
}
