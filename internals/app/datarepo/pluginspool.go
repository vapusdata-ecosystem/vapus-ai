package datarepo

import (
	"context"
	"fmt"
	"log"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	pluginsstore "github.com/vapusdata-ecosystem/vapusai/core/plugins"
)

func NewPluginPool(ctx context.Context, store *apppkgs.VapusStore, logger zerolog.Logger) (*pluginsstore.VapusPlugins, error) {
	var result []*models.Plugin
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted_at IS NULL", apppkgs.PluginsTable)
	err := store.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while fetching plugin by query from datastore")
		return nil, err
	}

	for _, plugin := range result {
		if plugin.Status == mpb.CommonStatus_ACTIVE.String() {
			log.Println("Plugin is active", plugin.Name, plugin.PluginService, plugin.Scope, "------------------------------------------")
			creds, err := apppkgs.ReadCredentialFromStore(ctx, plugin.NetworkParams.SecretName, store, logger)
			if err != nil {
				logger.Err(err).Ctx(ctx).Msgf("error while reading the secret for plugin - %s with secretname - %s of type %s", plugin.Name, plugin.NetworkParams.SecretName, plugin.PluginService)
				continue
			}
			log.Println(creds.GetGcpCreds(), "------------------------------------------creds")
			plugin.NetworkParams.Credentials = creds
		}
	}
	return pluginsstore.New(ctx, result, []string{}, logger), nil
}
