package dmstores

import (
	"context"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	aidmstore "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo/aistudio"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	pluginsstore "github.com/vapusdata-ecosystem/vapusai/core/plugins"
)

var AccountPool *models.Account

var PluginPool *pluginsstore.VapusPlugins

func InitPluginPool(ctx context.Context, ds *aidmstore.AIStudioDMStore) {
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
}

func NewPluginPool(ctx context.Context, ds *aidmstore.AIStudioDMStore) {
	if PluginPool == nil {
		InitPluginPool(ctx, ds)
	}
}
