package plugins

import (
	"context"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	"github.com/vapusdata-ecosystem/vapusdata/core/operator/emailer"
)

func (p *VapusPlugins) LoadEmailPlugin(ctx context.Context, plugin *models.Plugin) error {
	p.logger.Info().Msg("Loading filestore plugin")
	pl, err := emailer.New(ctx, plugin.PluginService, plugin.NetworkParams, plugin.DynamicParams, p.logger)
	if err != nil {
		p.logger.Err(err).Msg("Error creating email plugin")
		return err
	}
	switch plugin.Scope {
	case mpb.ResourceScope_PLATFORM_SCOPE.String():
		p.PlatformPlugins.Emailer = pl
	case mpb.ResourceScope_USER_SCOPE.String():
		pool, err := p.GetUserPluginPool(ctx, plugin.CreatedBy)
		if err != nil {
			p.logger.Err(err).Msg("Error getting user plugin pool")
			return err
		}
		pool.Emailer = pl
		p.LoadUserPluginPool(ctx, plugin.CreatedBy, pool)
	case mpb.ResourceScope_ORGANIZATION_SCOPE.String():
		pool, err := p.GetORGANIZATIONPluginPool(ctx, plugin.Organization)
		if err != nil {
			p.logger.Err(err).Msg("Error getting user plugin pool")
			return err
		}
		pool.Emailer = pl
		p.LoadORGANIZATIONPluginPool(ctx, plugin.Organization, pool)
	default:
		p.logger.Error().Msg("Invalid scope for email plugin")
		return apperr.ErrInvalidPluginScope
	}
	return nil
}
