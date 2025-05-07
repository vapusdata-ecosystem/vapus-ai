package dmstores

import (
	"context"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	pluginsstore "github.com/vapusdata-ecosystem/vapusdata/core/plugins"
)

var AccountPool *models.Account

var PluginPool *pluginsstore.VapusPlugins

func InitPluginPool(ctx context.Context, ds *DMStore) {
	plugins, err := ds.ListPlugins(ctx, "deleted_at IS NULL AND scope = 'PLATFORM_SCOPE'", map[string]string{})
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while getting the list of plugins")
		return
	}
	for _, plugin := range plugins {
		if plugin.Status == mpb.CommonStatus_ACTIVE.String() {
			creds, err := apppkgs.ReadCredentialFromStore(ctx, plugin.NetworkParams.SecretName, ds.VapusStore, logger)
			if err != nil {
				logger.Err(err).Ctx(ctx).Msg("error while reading the secret from the vault")
				continue
			}
			plugin.NetworkParams.Credentials = creds
		}
	}
	PluginPool = pluginsstore.New(ctx, plugins, []string{}, logger)
	logger.Info().Ctx(ctx).Msg("PluginPool initialized")
	return
}

func NewPluginPool(ctx context.Context, ds *DMStore) {
	if PluginPool == nil {
		InitPluginPool(ctx, ds)
	}
}

func InitAccountPool(ctx context.Context, ds *DMStore) {
	result := make([]*models.Account, 0)
	query := "select * from " + apppkgs.AccountsTable
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msg("error while getting the list of accounts")
	}
	AccountPool = result[0]
	logger.Info().Ctx(ctx).Msg("AccountPool initialized")
	return
}

func GetAccountFromPool(ctx context.Context, ds *DMStore, ctxClaim map[string]string) *models.Account {

	if AccountPool != nil {
		return AccountPool
	} else {
		InitAccountPool(ctx, ds)
		return AccountPool
	}
}
