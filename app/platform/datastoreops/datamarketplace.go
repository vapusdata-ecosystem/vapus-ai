package dmstores

import (
	"context"
	"fmt"

	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
)

// CreateAccount creates a new account in the data DataMarketplace store, one setup can have only one account
func (ds *DMStore) CreateAccount(ctx context.Context, obj *models.Account) (*models.Account, error) {
	logger.Info().Msgf("Creating account : %v", obj)
	_, err := ds.Db.PostgresClient.DB.NewInsert().ModelTableExpr(apppkgs.AccountsTable).Model(obj).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving account in datastore")
		return nil, err
	}

	return obj, nil
}

// GetAccount gets the account object from the data  store based on the key identifier i.e. accountid
func (ds *DMStore) GetAccount(ctx context.Context, ctxClaim map[string]string) (*models.Account, error) {
	result := make([]*models.Account, 0)
	query := fmt.Sprintf("SELECT * FROM %s WHERE vapus_id = '%s'", apppkgs.AccountsTable, ctxClaim[encryption.ClaimAccountKey])
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msg("error while getting account from datastore")
		return nil, dmerrors.DMError(apperr.ErrListingDataCatalog, err)
	}
	return result[0], err
}

func (ds *DMStore) PutAccount(ctx context.Context, obj *models.Account, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewUpdate().Model(obj).ModelTableExpr(apppkgs.AccountsTable).Where(apppkgs.VapusIdFilter(), obj.VapusID).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while updating account in datastore")
		return err
	}
	return nil
}
