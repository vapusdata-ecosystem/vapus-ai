package blobs

import (
	"context"
	"encoding/base64"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	dputils "github.com/vapusdata-ecosystem/vapusdata/core/data-platform/utils"
	models "github.com/vapusdata-ecosystem/vapusdata/core/models"
	options "github.com/vapusdata-ecosystem/vapusdata/core/options"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
	tpaws "github.com/vapusdata-ecosystem/vapusdata/core/thirdparty/aws"
	tpgcp "github.com/vapusdata-ecosystem/vapusdata/core/thirdparty/gcp"
)

type BlobStore interface {
	Close()
	CreateBucket(ctx context.Context, params *options.BlobOpsParams) error
	DeleteBucket(ctx context.Context, params *options.BlobOpsParams) error
	ListBuckets(ctx context.Context) ([]string, error)
	GetBucket(ctx context.Context, params *options.BlobOpsParams) (string, error)
	ListObjects(ctx context.Context, params *options.BlobOpsParams) ([]string, error)
	UploadObject(ctx context.Context, params *options.BlobOpsParams) error
	DeleteObject(ctx context.Context, params *options.BlobOpsParams) error
	DownloadObject(ctx context.Context, params *options.BlobOpsParams) ([]byte, error)
	// GetStorageData(ctx context.Context,params *options.BlobOpsParams) (*models.BlobStoreSchema,error)
}

type BlobStoreClient struct {
	BlobStore
	DataStoreParams *models.DataSourceCredsParams
	Logger          zerolog.Logger
	Debug           bool
}

func (d *BlobStoreClient) Close() {
	if d.BlobStore != nil {
		d.Close()
	}
}

type Options func(*BlobStoreClient)

func WithDebug(debug bool) Options {
	return func(d *BlobStoreClient) {
		d.Debug = debug
	}
}

func WithLogger(log zerolog.Logger) Options {
	return func(d *BlobStoreClient) {
		d.Logger = log
	}
}

func WithDataSourceCredsParams(params *models.DataSourceCredsParams) Options {
	return func(d *BlobStoreClient) {
		d.DataStoreParams = params
	}
}

// NewDataConnClient creates a new DataConnClient
func New(ctx context.Context, opts ...Options) (BlobStore, error) {
	resultCl := &BlobStoreClient{}
	for _, opt := range opts {
		opt(resultCl)
	}
	if opts == nil || resultCl.DataStoreParams.DataSourceEngine == "" || resultCl.DataStoreParams.DataSourceCreds == nil {
		return nil, dmerrors.DMError(dputils.ErrInvalidDataStoreEngine, dputils.ErrDataStoreConn)
	}
	switch resultCl.DataStoreParams.DataSourceService {
	case mpb.DataSourceServices_GCP_CLOUD_STORAGE.String():
		// if opts.DataSourceCreds.GcpCreds.Base64Encoded {}
		decodeData, err := base64.StdEncoding.DecodeString(resultCl.DataStoreParams.DataSourceCreds.GcpCreds.ServiceAccountKey)
		if err != nil {
			resultCl.Logger.Err(err).Msg("Error decoding gcp service account key")
			return nil, err
		}
		client, err := tpgcp.NewBucketAgent(ctx, &tpgcp.GcpConfig{
			ServiceAccountKey: []byte(decodeData),
			ProjectID:         resultCl.DataStoreParams.DataSourceCreds.GcpCreds.ProjectId,
			Region:            resultCl.DataStoreParams.DataSourceCreds.GcpCreds.Region,
		}, resultCl.Logger)
		if err != nil {
			log.Err(err).Msg("Error creating gcp secret manager client")
			return nil, err
		}
		resultCl.BlobStore = client
		return resultCl, nil
	case mpb.DataSourceServices_AWS_S3.String():
		client, err := tpaws.NewBucketAgent(ctx, &tpaws.AWSConfig{
			SecretAccessKey: resultCl.DataStoreParams.DataSourceCreds.AwsCreds.SecretAccessKey,
			AccessKeyId:     resultCl.DataStoreParams.DataSourceCreds.AwsCreds.AccessKeyId,
			Region:          resultCl.DataStoreParams.DataSourceCreds.AwsCreds.Region,
		}, resultCl.Logger)
		if err != nil {
			log.Err(err).Msg("Error creating gcp secret manager client")
			return nil, err
		}
		resultCl.BlobStore = client
		return resultCl, nil
	default:
		return nil, dmerrors.DMError(dputils.ErrInvalidBlobStoreSvc, dputils.ErrDataStoreConn)
	}
}
