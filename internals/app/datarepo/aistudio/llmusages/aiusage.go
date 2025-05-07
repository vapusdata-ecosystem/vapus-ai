package llm_price

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
)

func LLMPriceList(ctx context.Context, dmstore *apppkgs.VapusStore, logger zerolog.Logger) ([]*models.AIModelPriceList, error) {
	result := []*models.AIModelPriceList{}
	query := fmt.Sprintf(`SELECT * FROM %s`, apppkgs.AIModelPriceListTable)
	err := dmstore.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while getting llm price list")
		return nil, err
	}
	return result, err
}

func AddPriceList(ctx context.Context, dmstore *apppkgs.VapusStore, obj *models.AIModelPriceList, logger zerolog.Logger) error {
	_, err := dmstore.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(apppkgs.AIModelPriceListTable).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving price list in datastore")
		return err
	}
	return nil
}
func UpdatePriceList(ctx context.Context, dmstore *apppkgs.VapusStore, obj *models.AIModelPriceList, logger zerolog.Logger) error {
	_, err := dmstore.Db.PostgresClient.DB.NewUpdate().Model(obj).ModelTableExpr(apppkgs.AIModelPriceListTable).Where(`"llm_service_provider_name" = ?`, obj.LLMServiceProviderName).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while updating ai model prompt to datastore")
		return err
	}
	return nil
}
