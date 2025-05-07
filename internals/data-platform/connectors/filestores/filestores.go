package filestores

import (
	"context"

	"github.com/rs/zerolog"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	options "github.com/vapusdata-ecosystem/vapusai/core/options"
	gcp "github.com/vapusdata-ecosystem/vapusai/core/thirdparty/gcp"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

type FileStore interface {
	UploadFiles(ctx context.Context, opts *options.FileStoreUploadRequest) (*options.FileStoreUploadResponse, error)
	ListFiles(ctx context.Context, opts *options.FileStoreListRequest) (*options.FileStoreListResponse, error)
	DownloadFiles(ctx context.Context, request *options.FileStoreDownloadRequest) (*options.FileStoreDownloadResponse, error)
	DeleteFiles(ctx context.Context, opts *options.FileStoreDeleteRequest) error
	Share(ctx context.Context, opts *options.FileStoreShareRequest) (*options.FileStoreShareResponse, error)
	CreateFolder(ctx context.Context, opts *options.FileStoreCreateDirectoryRequest) (*options.FileStoreCreateDirectoryResponse, error)
	MoveFiles(ctx context.Context, opts *options.FileStoreMoveFilesRequest) error
	DownloadFolder(ctx context.Context, request *options.FileStoreDownloadRequest) (*options.FileStoreDownloadResponse, error)
	ShareFileLink(ctx context.Context) error
	DeleteFolder(ctx context.Context, opts *options.FileStoreDeleteRequest) error
}

type FileStoreClient struct {
	FileStore
	DataStoreParams *models.DataSourceCredsParams
	Logger          zerolog.Logger
	Debug           bool
}

type Options func(*FileStoreClient)

func WithDebug(debug bool) Options {
	return func(d *FileStoreClient) {
		d.Debug = debug
	}
}

func WithLogger(log zerolog.Logger) Options {
	return func(d *FileStoreClient) {
		d.Logger = log
	}
}

func WithDataSourceCredsParams(params *models.DataSourceCredsParams) Options {
	return func(d *FileStoreClient) {
		d.DataStoreParams = params
	}
}

func New(ctx context.Context, opts ...Options) (FileStore, error) {
	client := &FileStoreClient{}
	for _, opt := range opts {
		opt(client)
	}
	switch client.DataStoreParams.DataSourceService {
	case types.GOOGLE_DRIVE.String():
		client.Logger.Debug().Msgf("Creating Google Drive client - %s", client.DataStoreParams.DataSourceCreds.GetGcpCreds())
		return gcp.NewGDrive(ctx, &gcp.GcpConfig{
			ServiceAccountKey: []byte(client.DataStoreParams.DataSourceCreds.GetGcpCreds().ServiceAccountKey),
			ProjectID:         client.DataStoreParams.DataSourceCreds.GetGcpCreds().ProjectId,
			Region:            client.DataStoreParams.DataSourceCreds.GetGcpCreds().Region,
			Zone:              client.DataStoreParams.DataSourceCreds.GetGcpCreds().Zone,
		}, client.DataStoreParams.DataSourceCreds.Username, client.Logger)
	default:
		return nil, apperr.ErrInvalidFileStoreService
	}
}
