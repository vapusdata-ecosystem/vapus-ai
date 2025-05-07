package aidmstore

import (
	"context"
	"fmt"

	// mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	"github.com/databricks/databricks-sql-go/logger"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

// AddOrganization adds a new organization to the data organization
// It will create a new organization and attach it to the data organization
func (ds *AIStudioDMStore) ConfigureOrganization(ctx context.Context, organization *models.Organization, ctxClaim map[string]string) error {
	organization.SetAccountId(ctxClaim[encryption.ClaimAccountKey])
	if ds.Cacher != nil {
		_, err := ds.BeDataStore.Cacher.RedisClient.WrtiteData(ctx, organization.VapusID, types.EMPTYSTR, organization)
		if err != nil {
			logger.Err(err).Ctx(ctx).Msg(apperr.ErrOrganizationInitialization.Error())
			return err
		}
	}
	_, err := ds.Db.PostgresClient.DB.NewInsert().ModelTableExpr(apppkgs.OrganizationsTable).Model(organization).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving datamarketplace in datastore")
		return err
	}
	go func() {
		mCtx := context.TODO()
		_ = apppkgs.AddResourceArn(mCtx, ds.Db, &models.VapusResourceArn{
			ResourceId:   organization.VapusID,
			ResourceName: mpb.Resources_ORGANIZATIONS.String(),
			VapusBase: models.VapusBase{
				Editors: organization.Editors,
			},
		}, ds.logger, ctxClaim)
	}()
	return nil
}

func (ds *AIStudioDMStore) PatchOrganization(ctx context.Context, data, conditions map[string]interface{}, ctxClaim map[string]string) error {
	pq := ds.Db.PostgresClient.DB.NewUpdate().Model(&data).ModelTableExpr(apppkgs.OrganizationsTable)

	for key, value := range conditions {
		pq = pq.Where(key, value)
	}
	_, err := pq.Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving organization info in datastore")
		return err
	}
	return nil
}

func (ds *AIStudioDMStore) PutOrganization(ctx context.Context, obj *models.Organization, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewUpdate().Model(obj).ModelTableExpr(apppkgs.OrganizationsTable).Where(apppkgs.VapusIdFilter(), obj.VapusID).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while updating organization in datastore")
		return err
	}
	return nil
}

// GetOrganization gets the organization object from the data organization store based on the key identifier i.e. organizationid
// GetMarketplaceOrganizations retrieves the organizations associated with a given key from the DMStore.
// It returns a slice of *models.Organization, a map[string]interface{} for custom messages, and an error.
// The key parameter specifies the key to query the DMStore.
// If the retrieval is successful, the function returns the organizations, custom messages, and a nil error.
// If an error occurs during retrieval, the function returns nil for organizations, custom messages, and an error indicating the cause.
func (ds *AIStudioDMStore) ListOrganizations(ctx context.Context, condition string, ctxClaim map[string]string) ([]*models.Organization, error) {
	result := []*models.Organization{}
	condition = apppkgs.GetAccountFilter(ctxClaim, condition)
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s ORDER BY created_at DESC ", apppkgs.OrganizationsTable, condition)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	// err := ds.Db.Select(ctx, &datasvcpkgs.QueryOpts{
	// 	DataCollection: apppkgs.OrganizationsTable,
	// 	RawQuery:       query,
	// }, &result, logger)
	// resp, err := ds.Db.ElasticSearchClient.TClient.Search().Index(datasvc.VAPUS_DOMAIN_INDEX).Request(query).Do(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msgf("error while getting organization for the request")
		return nil, err
	}
	return result, nil
}

func (ds *AIStudioDMStore) GetOrganization(ctx context.Context, iden string, ctxClaim map[string]string) (*models.Organization, error) {
	result := []*models.Organization{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", apppkgs.OrganizationsTable, apppkgs.GetByIdFilter("", iden, ctxClaim))
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msgf("error while getting organization for the request")
		return nil, dmerrors.DMError(apperr.ErrInvalidOrganizationRequested, err)
	}
	return result[0], nil
}
