package connectors

import (
	"context"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusai/core/data-platform/connectors/artifacts"
	"github.com/vapusdata-ecosystem/vapusai/core/data-platform/connectors/blobs"
	coderepos "github.com/vapusdata-ecosystem/vapusai/core/data-platform/connectors/coderepo"
	"github.com/vapusdata-ecosystem/vapusai/core/data-platform/connectors/databases"
	emailstore "github.com/vapusdata-ecosystem/vapusai/core/data-platform/connectors/emailstores"
	"github.com/vapusdata-ecosystem/vapusai/core/data-platform/connectors/filestores"
	secretstore "github.com/vapusdata-ecosystem/vapusai/core/data-platform/connectors/secrets-stores"
	dputils "github.com/vapusdata-ecosystem/vapusai/core/data-platform/utils"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
)

type DataSourceClient struct {
	blobs.BlobStore
	filestores.FileStore
	coderepos.CodeRepoStore
	*secretstore.SecretStoreClient
	emailstore.EmailStore
	*databases.DataStoreClient
	artifacts.ArtifactStore
	Opts *DataSourceOpts
}

type DataSourceOpts struct {
	Debug  bool
	InApp  bool
	Logger zerolog.Logger
	Params *models.DataSourceCredsParams
}

type Options func(*DataSourceOpts)

func WithDebug(debug bool) Options {
	return func(d *DataSourceOpts) {
		d.Debug = debug
	}
}

func WithInApp(inApp bool) Options {
	return func(d *DataSourceOpts) {
		d.InApp = inApp
	}
}

func WithLogger(log zerolog.Logger) Options {
	return func(d *DataSourceOpts) {
		d.Logger = log
	}
}

func WithDataSourceCredsParams(params *models.DataSourceCredsParams) Options {
	return func(d *DataSourceOpts) {
		d.Params = params
	}
}

func New(ctx context.Context, options ...Options) (*DataSourceClient, error) {
	opts := &DataSourceOpts{}
	for _, option := range options {
		option(opts)
	}
	if opts.Params == nil {
		return nil, dmerrors.DMError(dputils.ErrInvalidDataStoreEngine, dputils.ErrDataStoreConn)
	}
	dataSourceClient := &DataSourceClient{
		Opts: opts,
	}
	switch opts.Params.DataSourceType {
	case mpb.DataSourceType_BLOB_STORE.String():
		client, err := blobs.New(ctx, blobs.WithDataSourceCredsParams(opts.Params), blobs.WithLogger(opts.Logger), blobs.WithDebug(opts.Debug))
		if err != nil {
			opts.Logger.Err(err).Msg("Error creating blob store client")
			return nil, dmerrors.DMError(dputils.ErrInvalidDataStoreCredentials, err)
		}
		dataSourceClient.BlobStore = client
	case mpb.DataSourceType_DATABASE.String():
		client, err := databases.New(ctx, databases.WithDataSourceCredsParams(opts.Params), databases.WithLogger(opts.Logger), databases.WithDebug(opts.Debug))
		if err != nil {
			opts.Logger.Err(err).Msg("Error creating database client")
			return nil, dmerrors.DMError(dputils.ErrInvalidDataStoreCredentials, err)
		}
		dataSourceClient.DataStoreClient = client
	case mpb.DataSourceType_ARTIFACT.String():
		client, err := artifacts.New(ctx, artifacts.WithDataSourceCredsParams(opts.Params), artifacts.WithLogger(opts.Logger), artifacts.WithDebug(opts.Debug))
		if err != nil {
			opts.Logger.Err(err).Msg("Error creating code repository client")
			return nil, dmerrors.DMError(dputils.ErrInvalidDataStoreCredentials, err)
		}
		dataSourceClient.ArtifactStore = client
	case mpb.DataSourceType_CODE_REPOSITORY.String():
		client, err := coderepos.New(ctx, coderepos.WithDataSourceCredsParams(opts.Params), coderepos.WithLogger(opts.Logger), coderepos.WithDebug(opts.Debug))
		if err != nil {
			opts.Logger.Err(err).Msg("Error creating code repository client")
			return nil, dmerrors.DMError(dputils.ErrInvalidDataStoreCredentials, err)
		}
		dataSourceClient.CodeRepoStore = client
	case mpb.DataSourceType_FILE_STORE.String():
		client, err := filestores.New(ctx, filestores.WithDataSourceCredsParams(opts.Params), filestores.WithLogger(opts.Logger), filestores.WithDebug(opts.Debug))
		if err != nil {
			opts.Logger.Err(err).Msg("Error creating file store client")
			return nil, dmerrors.DMError(dputils.ErrInvalidDataStoreCredentials, err)
		}
		dataSourceClient.FileStore = client
	case mpb.DataSourceType_EMAIL_STORE.String():
		client, err := emailstore.New(ctx, emailstore.WithDataSourceCredsParams(opts.Params), emailstore.WithLogger(opts.Logger), emailstore.WithDebug(opts.Debug))
		if err != nil {
			opts.Logger.Err(err).Msg("Error creating email store client")
			return nil, dmerrors.DMError(dputils.ErrInvalidDataStoreCredentials, err)
		}
		dataSourceClient.EmailStore = client
	case mpb.DataSourceType_SECRET_STORE.String():
		client, err := secretstore.New(ctx, secretstore.WithDataSourceCredsParams(opts.Params), secretstore.WithLogger(opts.Logger), secretstore.WithDebug(opts.Debug))
		if err != nil {
			opts.Logger.Err(err).Msg("Error creating secret store client")
			return nil, dmerrors.DMError(dputils.ErrInvalidDataStoreCredentials, err)
		}
		dataSourceClient.SecretStore = client
	}
	return dataSourceClient, nil
}
