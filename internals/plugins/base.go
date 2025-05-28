package plugins

import (
	"context"
	"log"
	"sync"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusai/core/data-platform/connectors/filestores"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	"github.com/vapusdata-ecosystem/vapusai/core/operator/calendar"
	"github.com/vapusdata-ecosystem/vapusai/core/operator/emailer"
	searchengine "github.com/vapusdata-ecosystem/vapusai/core/operator/search"
)

type PluginsBase struct {
	Emailer      emailer.Emailer
	FileManager  filestores.FileStore
	SearchEngine searchengine.Search
	Calendar     calendar.Calendar
}

type VapusPlugins struct {
	PlatformPlugins *PluginsBase
	// ORGANIZATIONPlugins       map[string]*PluginsBase
	// UserPlugins         map[string]*PluginsBase
	OrganizationPlugins *sync.Map
	UserPlugins         *sync.Map
	logger              zerolog.Logger
	PluginsNotConencted []string
}

var PluginTypeScopeMap = map[string][]string{
	mpb.IntegrationPluginTypes_EMAIL.String(): {
		mpb.ResourceScope_PLATFORM_SCOPE.String(),
	},
	mpb.IntegrationPluginTypes_FILESTORES.String(): {
		mpb.ResourceScope_USER_SCOPE.String(),
		mpb.ResourceScope_ORGANIZATION_SCOPE.String(),
	},
	mpb.IntegrationPluginTypes_CODE_REPOSITORIES.String(): {
		mpb.ResourceScope_USER_SCOPE.String(),
		mpb.ResourceScope_ORGANIZATION_SCOPE.String(),
	},
	mpb.IntegrationPluginTypes_SMS.String(): {
		mpb.ResourceScope_PLATFORM_SCOPE.String(),
		mpb.ResourceScope_USER_SCOPE.String(),
		mpb.ResourceScope_ORGANIZATION_SCOPE.String(),
	},
	mpb.IntegrationPluginTypes_SEARCHAPI.String(): {
		mpb.ResourceScope_PLATFORM_SCOPE.String(),
		mpb.ResourceScope_USER_SCOPE.String(),
		mpb.ResourceScope_ORGANIZATION_SCOPE.String(),
	},
}

func New(ctx context.Context, pluginList []*models.Plugin, selectedPlugins []string, logger zerolog.Logger) *VapusPlugins {
	obj := &VapusPlugins{
		PlatformPlugins:     &PluginsBase{},
		OrganizationPlugins: &sync.Map{},
		UserPlugins:         &sync.Map{},
		logger:              logger,
	}
	logger.Info().Msg("Creating new vapus plugins pool")
	for _, plugin := range pluginList {
		err := obj.LoadPlugin(ctx, plugin)
		if err != nil {
			logger.Err(err).Msg("Error loading plugin")
			obj.PluginsNotConencted = append(obj.PluginsNotConencted, plugin.VapusID)
			continue
		}
	}
	return obj
}

func (p *VapusPlugins) LoadPlugin(ctx context.Context, plugin *models.Plugin) error {
	switch plugin.PluginType {
	case mpb.IntegrationPluginTypes_EMAIL.String():
		p.LoadEmailPlugin(ctx, plugin)
	case mpb.IntegrationPluginTypes_SEARCHAPI.String():
		p.LoadSearchPlugin(ctx, plugin)
	case mpb.IntegrationPluginTypes_FILESTORES.String():
		p.LoadFileStorePlugin(ctx, plugin)
	case mpb.IntegrationPluginTypes_CALENDAR.String():
		p.LoadCalendarPlugin(ctx, plugin)
	default:
		p.logger.Error().Msgf("Invalid plugin type - %s", plugin.PluginType)
		return apperr.ErrInvalidPluginType
	}
	log.Println("================================================================plugin", plugin.Scope, "Service-----", plugin.PluginService)
	log.Println(p.UserPlugins)
	return nil
}

func (p *VapusPlugins) GetORGANIZATIONPluginPool(ctx context.Context, ORGANIZATION string) (*PluginsBase, error) {
	pluginPool, ok := p.OrganizationPlugins.Load(ORGANIZATION)
	if !ok {
		return &PluginsBase{}, nil
	}
	ORGANIZATIONPlugin, ok := pluginPool.(*PluginsBase)
	if !ok {
		return nil, apperr.ErrInvalidPluginObject
	} else {
		return ORGANIZATIONPlugin, nil
	}
}

func (p *VapusPlugins) LoadORGANIZATIONPluginPool(ctx context.Context, ORGANIZATION string, pluginBase *PluginsBase) {
	p.OrganizationPlugins.Store(ORGANIZATION, pluginBase)
}

func (p *VapusPlugins) GetUserPluginPool(ctx context.Context, owner string) (*PluginsBase, error) {
	pluginPool, ok := p.UserPlugins.Load(owner)
	if !ok {
		return &PluginsBase{}, nil
	}
	userPluginPool, ok := pluginPool.(*PluginsBase)

	if !ok {
		return nil, apperr.ErrInvalidPluginObject
	} else {
		return userPluginPool, nil
	}
}

func (p *VapusPlugins) LoadUserPluginPool(ctx context.Context, owner string, pluginBase *PluginsBase) {
	p.UserPlugins.Store(owner, pluginBase)
}
