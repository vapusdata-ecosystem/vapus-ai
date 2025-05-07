package plugins

import (
	"context"
	"log"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusai/core/data-platform/connectors/filestores"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
)

func (p *VapusPlugins) LoadFileStorePlugin(ctx context.Context, plugin *models.Plugin) error {
	p.logger.Info().Msg("Loading filestore plugin")
	params := &models.DataSourceCredsParams{
		DataSourceCreds: &models.DataSourceSecrets{
			GenericCredentialModel: plugin.NetworkParams.Credentials,
			URL:                    plugin.NetworkParams.URL,
			Port:                   plugin.NetworkParams.Port,
		},
		DataSourceService: plugin.PluginService,
	}
	pl, err := filestores.New(ctx, filestores.WithLogger(p.logger),
		filestores.WithDataSourceCredsParams(params))
	if err != nil {
		p.logger.Err(err).Msg("Error creating filestore plugin")
		return apperr.ErrPluginNotConnected
	}
	p.logger.Info().Msgf("Plugin created successfully - %s", pl)
	switch plugin.Scope {
	case mpb.ResourceScope_USER_SCOPE.String():
		pool, err := p.GetUserPluginPool(ctx, plugin.CreatedBy)
		if err != nil {
			p.logger.Err(err).Msg("Error getting user plugin pool")
			return err
		}
		log.Println("User Plugin pool", pool)
		pool.FileManager = pl
		p.LoadUserPluginPool(ctx, plugin.CreatedBy, pool)
	case mpb.ResourceScope_ORGANIZATION_SCOPE.String():
		pool, err := p.GetORGANIZATIONPluginPool(ctx, plugin.Organization)
		if err != nil {
			p.logger.Err(err).Msg("Error getting ORGANIZATION plugin pool")
			return err
		}
		pool.FileManager = pl
		p.LoadORGANIZATIONPluginPool(ctx, plugin.Organization, pool)
	default:
		p.logger.Error().Msg("Invalid scope for filestore plugin")
		return apperr.ErrInvalidPluginScope
	}
	return nil
}
