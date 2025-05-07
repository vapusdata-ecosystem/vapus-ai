package plugins

import (
	"context"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	searchengine "github.com/vapusdata-ecosystem/vapusai/core/operator/search"
)

func (p *VapusPlugins) LoadSearchPlugin(ctx context.Context, plugin *models.Plugin) error {
	p.logger.Info().Msg("Loading search plugin")
	pl, err := searchengine.New(ctx, plugin.PluginService, plugin.NetworkParams, plugin.DynamicParams, p.logger)
	if err != nil {
		p.logger.Error().Msg("Error creating search plugin")
		return apperr.ErrPluginNotConnected
	}
	switch plugin.Scope {
	case mpb.ResourceScope_PLATFORM_SCOPE.String():
		p.PlatformPlugins.SearchEngine = pl
	case mpb.ResourceScope_USER_SCOPE.String():
		pool, err := p.GetUserPluginPool(ctx, plugin.CreatedBy)
		if err != nil {
			p.logger.Err(err).Msg("Error getting user plugin pool")
			return err
		}
		pool.SearchEngine = pl
		p.LoadUserPluginPool(ctx, plugin.CreatedBy, pool)
	case mpb.ResourceScope_ORGANIZATION_SCOPE.String():
		pool, err := p.GetORGANIZATIONPluginPool(ctx, plugin.Organization)
		if err != nil {
			p.logger.Err(err).Msg("Error getting ORGANIZATION plugin pool")
			return err
		}
		pool.SearchEngine = pl
		p.LoadORGANIZATIONPluginPool(ctx, plugin.Organization, pool)
	default:
		p.logger.Error().Msg("Invalid scope for ORGANIZATION email plugin")
		return apperr.ErrInvalidPluginScope
	}
	return nil
}
