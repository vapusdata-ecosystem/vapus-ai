package plugins

import (
	"context"
	// "fmt"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	"github.com/vapusdata-ecosystem/vapusai/core/operator/calendar"
)

func (p *VapusPlugins) LoadCalendarPlugin(ctx context.Context, plugin *models.Plugin) error {
	p.logger.Info().Msg("Loading calendar plugin")

	if plugin.NetworkParams == nil || plugin.NetworkParams.Credentials == nil {
		p.logger.Error().Msg("Invalid calendar plugin credentials")
		return apperr.ErrInvalidCalenderService
	}

	pl, err := calendar.New(ctx, plugin.PluginService, plugin.NetworkParams, nil, p.logger)

	if err != nil {
		p.logger.Err(err).Msg("Error creating calendar plugin")
		return err
	}

	switch plugin.Scope {
	case mpb.ResourceScope_USER_SCOPE.String():

		pool, err := p.GetUserPluginPool(ctx, plugin.CreatedBy)
		if err != nil {
			p.logger.Err(err).Msg("Error getting user plugin pool")
			return err
		}
		pool.Calendar = pl
		p.LoadUserPluginPool(ctx, plugin.CreatedBy, pool)
		// fmt.Println("pool:", pool)
	default:
		p.logger.Error().Msg("Invalid scope for calendar plugin. Only USER_SCOPE is supported.")
		return apperr.ErrInvalidPluginScope
	}

	return nil
}
