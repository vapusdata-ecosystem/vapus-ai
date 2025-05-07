package coderepos

import (
	"context"

	"github.com/rs/zerolog"
	// velasticsearch "github.com/vapusdata-ecosystem/vapusdata/core/data-platform/dataservices/elasticsearch"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	githubsvc "github.com/vapusdata-ecosystem/vapusdata/core/data-platform/dataservices/githubconnector"
	dputils "github.com/vapusdata-ecosystem/vapusdata/core/data-platform/utils"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
	logger "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/logger"
)

type CodeRepoStore interface{}

type CodeRepoClient struct {
	CodeRepoStore
	GithubClient    *githubsvc.GithubStore
	DataStoreParams *models.DataSourceCredsParams
	Logger          zerolog.Logger
	Debug           bool
}

var (
	log = logger.CoreLogger
)

func (d *CodeRepoClient) Close() {

	if d.GithubClient != nil {
		d.GithubClient.Close()
	}
}

type Options func(*CodeRepoClient)

func WithDebug(debug bool) Options {
	return func(d *CodeRepoClient) {
		d.Debug = debug
	}
}

func WithLogger(log zerolog.Logger) Options {
	return func(d *CodeRepoClient) {
		d.Logger = log
	}
}

func WithDataSourceCredsParams(params *models.DataSourceCredsParams) Options {
	return func(d *CodeRepoClient) {
		d.DataStoreParams = params
	}
}

// NewDataConnClient creates a new DataConnClient
func New(ctx context.Context, opts ...Options) (*CodeRepoClient, error) {
	dsc := &CodeRepoClient{}
	for _, opt := range opts {
		opt(dsc)
	}
	if dsc.DataStoreParams == nil {
		return nil, dmerrors.DMError(dputils.ErrDataStoreParams404, dputils.ErrDataStoreParams404)
	}
	log.Debug().Msgf("Creating new data source connection for %s", dsc.DataStoreParams.DataSourceEngine)
	switch dsc.DataStoreParams.DataSourceSvcProvider {
	case mpb.DataSourceServices_GITHUB_SVC.String():
		client, err := githubsvc.New(&githubsvc.GithubOpts{
			Pat: dsc.DataStoreParams.DataSourceCreds.GenericCredentialModel.ApiToken,
		}, log)
		if err != nil {
			log.Err(err).Msg("Error connecting to redis")
			return nil, err
		}
		return &CodeRepoClient{GithubClient: client, DataStoreParams: dsc.DataStoreParams, Logger: log}, nil
	default:
		return nil, dmerrors.DMError(dputils.ErrInvalidDataStoreEngine, dputils.ErrDataStoreConn)
	}
}
