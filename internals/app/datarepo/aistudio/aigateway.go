package aidmstore

import (
	"context"
	"fmt"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	models "github.com/vapusdata-ecosystem/vapusdata/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
)

func (ds *AIStudioDMStore) CreateAIStudioChat(ctx context.Context, obj *models.AIStudioChat, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(apppkgs.AIStudioChatsTable).Exec(ctx)
	if err != nil {
		ds.logger.Err(err).Ctx(ctx).Msg("error while saving ai studio chat metadata in datastore")
		return err
	}
	go func() {
		mCtx := context.TODO()
		_ = apppkgs.AddResourceArn(mCtx, ds.Db, &models.VapusResourceArn{
			ResourceId:   obj.VapusID,
			ResourceName: mpb.Resources_AIMODELS.String(),
			VapusBase: models.VapusBase{
				Editors: obj.Editors,
			},
		}, ds.logger, ctxClaim)
	}()
	return nil
}

func (ds *AIStudioDMStore) ListAIStudioChats(ctx context.Context, limit int, ctxClaim map[string]string) ([]*models.AIStudioChat, error) {
	resp := []*models.AIStudioChat{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE created_by='%s' AND Organization='%s' AND owner_account='%s' ORDER BY id DESC",
		apppkgs.AIStudioChatsTable,
		ctxClaim[encryption.ClaimUserIdKey],
		ctxClaim[encryption.ClaimOrganizationKey],
		ctxClaim[encryption.ClaimAccountKey],
	)
	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, limit)
	}
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &resp)
	if err != nil {
		ds.logger.Err(err).Ctx(ctx).Msg("error while fetching ai studio chat from datastore")
		return nil, err
	}
	return resp, nil
}

func (ds *AIStudioDMStore) PutAIStudioChat(ctx context.Context, obj *models.AIStudioChat) error {
	_, err := ds.Db.PostgresClient.DB.NewUpdate().Model(obj).ModelTableExpr(apppkgs.AIStudioChatsTable).Where(apppkgs.VapusIdFilter(), obj.VapusID).Exec(ctx)
	if err != nil {
		ds.logger.Err(err).Ctx(ctx).Msg("error while updating ai studio chat metadata in datastore")
		return err
	}
	return nil
}

func (ds *AIStudioDMStore) GetAIStudioChat(ctx context.Context, iden string, ctxClaim map[string]string) (*models.AIStudioChat, error) {
	if iden == "" {
		return nil, apperr.ErrAIModelNode404
	}
	result := &models.AIStudioChat{}
	err := ds.Db.PostgresClient.DB.NewSelect().
		Model(result).
		Relation("Messages").
		Where("chat_id = ?", iden).
		Where("created_by = ?", ctxClaim[encryption.ClaimUserIdKey]).
		Where("Organization = ?", ctxClaim[encryption.ClaimOrganizationKey]).
		Where("owner_account = ?", ctxClaim[encryption.ClaimAccountKey]).OrderExpr("created_at DESC").
		Scan(ctx)
	return result, err
}
