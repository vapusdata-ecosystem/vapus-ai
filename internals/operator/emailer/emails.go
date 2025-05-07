package emailer

import (
	"context"
	"encoding/base64"

	"github.com/rs/zerolog"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	options "github.com/vapusdata-ecosystem/vapusai/core/options"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	aws "github.com/vapusdata-ecosystem/vapusai/core/thirdparty/aws"
	gcp "github.com/vapusdata-ecosystem/vapusai/core/thirdparty/gcp"
	sendgrid "github.com/vapusdata-ecosystem/vapusai/core/tools/sendgrid"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

type Emailer interface {
	SendEmail(ctx context.Context, opts *options.SendEmailRequest, agentId string) error
	SendRawEmail(ctx context.Context, opts *options.SendEmailRequest, agentId string) error
}

type EmailerClient struct {
	logger zerolog.Logger
	Emailer
}

func New(ctx context.Context,
	service string,
	netOps *models.PluginNetworkParams,
	ops []*models.Mapper,
	logger zerolog.Logger) (Emailer, error) {
	if netOps == nil || netOps.Credentials == nil {
		logger.Error().Msg("Invalid emailer credentials")
		return nil, dmerrors.DMError(apperr.ErrInvalidEmailerCredentials, apperr.ErrEmailerConn)
	}
	emailer := &EmailerClient{
		logger: logger,
	}
	switch service {
	case types.AWS_SES.String():
		client, err := aws.NewAwsSesClient(ctx, &aws.AWSConfig{
			Region:          netOps.Credentials.AwsCreds.Region,
			AccessKeyId:     netOps.Credentials.AwsCreds.AccessKeyId,
			SecretAccessKey: netOps.Credentials.AwsCreds.SecretAccessKey,
		}, logger)
		if err != nil {
			logger.Err(err).Msg("Error creating aws ses client")
			return nil, err
		}
		emailer.Emailer = client
	case types.GMAIL.String():
		decodeData, err := base64.StdEncoding.DecodeString(netOps.Credentials.GcpCreds.ServiceAccountKey)
		if err != nil {
			logger.Err(err).Msg("Error decoding gcp service account key")
			return nil, err
		}
		client, err := gcp.NewGmailService(ctx, &gcp.GcpConfig{
			ServiceAccountKey: []byte(decodeData),
			ProjectID:         netOps.Credentials.GcpCreds.ProjectId,
			Region:            netOps.Credentials.GcpCreds.Region,
		}, logger)
		if err != nil {
			logger.Err(err).Msg("Error creating gcp secret manager client")
			return nil, err
		}
		emailer.Emailer = client
	case types.SENDGRID.String():
		client, err := sendgrid.New(ctx, sendgrid.WithLogger(logger), sendgrid.WithCreds(netOps), sendgrid.WithParams(ops))
		if err != nil {
			logger.Err(err).Msg("Error creating sendgrid client")
			return nil, err
		}
		emailer.Emailer = client
	default:
		logger.Error().Msg("Invalid emailer engine")
		return nil, dmerrors.DMError(apperr.ErrInvalidEmailerEngine, apperr.ErrEmailerConn)
	}
	return emailer, nil
}
