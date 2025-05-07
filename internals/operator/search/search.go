package searchengine

import (
	"context"

	"github.com/rs/zerolog"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	options "github.com/vapusdata-ecosystem/vapusdata/core/options"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
	"github.com/vapusdata-ecosystem/vapusdata/core/tools/serp"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
)

type Search interface {
	SearchFormatted(opts *options.SearchInput) (*options.SearchResult, error)
	SearchRaw(opts *options.SearchInput) (map[string]any, error)
}

func New(ctx context.Context,
	service string,
	netOps *models.PluginNetworkParams,
	ops []*models.Mapper,
	logger zerolog.Logger) (Search, error) {
	if netOps == nil || netOps.Credentials == nil {
		logger.Error().Msg("Invalid search plugin credentials")
		return nil, dmerrors.DMError(apperr.ErrInvalidPluginCredentials, nil)
	}
	switch service {
	case types.SERPSEARCH.String():
		return serp.New(serp.WithCreds(netOps)), nil
	default:
		return nil, dmerrors.DMError(apperr.ErrInvalidPluginService, nil)
	}
}
