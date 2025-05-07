package emailstore

import (
	"context"
	"encoding/base64"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	dputils "github.com/vapusdata-ecosystem/vapusai/core/data-platform/utils"
	models "github.com/vapusdata-ecosystem/vapusai/core/models"
	options "github.com/vapusdata-ecosystem/vapusai/core/options"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	tpgcp "github.com/vapusdata-ecosystem/vapusai/core/thirdparty/gcp"
)

type EmailStore interface {
	Close()
	SendEmail(ctx context.Context, opts *options.SendEmailRequest, agentId string) error
	SendRawEmail(ctx context.Context, opts *options.SendEmailRequest, agentId string) error
	// LoadMetadata(ctx context.Context, opts *options.ListEmailsRequest) ([]*models.EmailStoreMetadata, error)
	Delete(ctx context.Context, opts *options.DeleteEmailRequest) error
	Watch(ctx context.Context, opts *options.ListEmailsRequest) error
	// ListEmails(ctx context.Context, opts *options.ListEmailsRequest) ([]*models.EmailMetadata, error)
	// ReadEmail(ctx context.Context, opts *options.ReadEmailRequest) (*models.EmailMetadata, error)
}

type EmailStoreClient struct {
	EmailStore
	DataStoreParams *models.DataSourceCredsParams
	Logger          zerolog.Logger
	Debug           bool
}

type Options func(*EmailStoreClient)

func WithDebug(debug bool) Options {
	return func(d *EmailStoreClient) {
		d.Debug = debug
	}
}

func WithLogger(log zerolog.Logger) Options {
	return func(d *EmailStoreClient) {
		d.Logger = log
	}
}

func WithDataSourceCredsParams(params *models.DataSourceCredsParams) Options {
	return func(d *EmailStoreClient) {
		d.DataStoreParams = params
	}
}

func (d *EmailStoreClient) Close() {
	if d.EmailStore != nil {
		d.Close()
	}
}

func New(ctx context.Context, opts ...Options) (EmailStore, error) {
	resultCl := &EmailStoreClient{}
	for _, opt := range opts {
		opt(resultCl)
	}
	if opts == nil || resultCl.DataStoreParams.DataSourceEngine == "" || resultCl.DataStoreParams.DataSourceCreds == nil {
		return nil, dmerrors.DMError(dputils.ErrInvalidDataStoreEngine, dputils.ErrDataStoreConn)
	}
	switch resultCl.DataStoreParams.DataSourceService {
	case mpb.DataSourceServices_GMAIL.String():
		decodeData, err := base64.StdEncoding.DecodeString(resultCl.DataStoreParams.DataSourceCreds.GcpCreds.ServiceAccountKey)
		if err != nil {
			resultCl.Logger.Err(err).Msg("Error decoding gcp service account key")
			return nil, err
		}
		client, err := tpgcp.NewGmailSource(ctx, &tpgcp.GcpConfig{
			ServiceAccountKey: []byte(decodeData),
			ProjectID:         resultCl.DataStoreParams.DataSourceCreds.GcpCreds.ProjectId,
			Region:            resultCl.DataStoreParams.DataSourceCreds.GcpCreds.Region,
		}, resultCl.Logger)
		if err != nil {
			log.Err(err).Msg("Error creating gcp secret manager client")
			return nil, err
		}
		resultCl.EmailStore = client
		return resultCl, nil

	default:
		return nil, dmerrors.DMError(dputils.ErrInvalidEmailStoreSvc, dputils.ErrDataStoreConn)
	}
}
