package artifacts

import (
	"context"
	"encoding/base64"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	dputils "github.com/vapusdata-ecosystem/vapusdata/core/data-platform/utils"
	models "github.com/vapusdata-ecosystem/vapusdata/core/models"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
	dmhttp "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/http"
	awstp "github.com/vapusdata-ecosystem/vapusdata/core/thirdparty/aws"
	gcptp "github.com/vapusdata-ecosystem/vapusdata/core/thirdparty/gcp"
)

type ArtifactStore interface {
}

type ArtifactStoreClient struct {
	*gcptp.GcpArManager
	*awstp.ECRManager
	DataStoreParams *models.DataSourceCredsParams
	Logger          zerolog.Logger
	Debug           bool
}

type Options func(*ArtifactStoreClient)

func WithDebug(debug bool) Options {
	return func(d *ArtifactStoreClient) {
		d.Debug = debug
	}
}

func WithLogger(log zerolog.Logger) Options {
	return func(d *ArtifactStoreClient) {
		d.Logger = log
	}
}

func WithDataSourceCredsParams(params *models.DataSourceCredsParams) Options {
	return func(d *ArtifactStoreClient) {
		d.DataStoreParams = params
	}
}

// NewDataConnClient creates a new DataConnClient
func New(ctx context.Context, opts ...Options) (ArtifactStore, error) {
	resultCl := &ArtifactStoreClient{}
	for _, opt := range opts {
		opt(resultCl)
	}
	if opts == nil || resultCl.DataStoreParams.DataSourceEngine == "" || resultCl.DataStoreParams.DataSourceCreds == nil {
		return nil, dmerrors.DMError(dputils.ErrInvalidDataStoreEngine, dputils.ErrDataStoreConn)
	}
	switch resultCl.DataStoreParams.DataSourceService {
	case mpb.DataSourceServices_GAR.String():
		// if opts.DataSourceCreds.GcpCreds.Base64Encoded {}
		decodeData, err := base64.StdEncoding.DecodeString(resultCl.DataStoreParams.DataSourceCreds.GcpCreds.ServiceAccountKey)
		if err != nil {
			resultCl.Logger.Err(err).Msg("Error decoding gcp service account key")
			return nil, err
		}
		client, err := gcptp.NewGcpArManager(ctx, &gcptp.GcpConfig{
			ServiceAccountKey: []byte(decodeData),
			ProjectID:         resultCl.DataStoreParams.DataSourceCreds.GcpCreds.ProjectId,
			Region:            resultCl.DataStoreParams.DataSourceCreds.GcpCreds.Region,
		}, dmhttp.GetOrAddUrlScheme(resultCl.DataStoreParams.DataSourceCreds.URL), resultCl.Logger)
		if err != nil {
			log.Err(err).Msg("Error creating gcp secret manager client")
			return nil, err
		}
		resultCl.GcpArManager = client
		return resultCl, nil
	case mpb.DataSourceServices_ECR.String():
		client, err := awstp.NewECRManager(ctx, &awstp.AWSConfig{
			SecretAccessKey: resultCl.DataStoreParams.DataSourceCreds.AwsCreds.SecretAccessKey,
			AccessKeyId:     resultCl.DataStoreParams.DataSourceCreds.AwsCreds.AccessKeyId,
			Region:          resultCl.DataStoreParams.DataSourceCreds.AwsCreds.Region,
		}, dmhttp.GetOrAddUrlScheme(resultCl.DataStoreParams.DataSourceCreds.URL), resultCl.Logger)
		if err != nil {
			log.Err(err).Msg("Error creating gcp secret manager client")
			return nil, err
		}
		resultCl.ECRManager = client
		return resultCl, nil
	default:
		return nil, dmerrors.DMError(dputils.ErrInvalidArtifactStoreSvc, dputils.ErrDataStoreConn)
	}
}
