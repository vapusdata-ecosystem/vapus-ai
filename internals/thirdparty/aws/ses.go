package vaws

import (
	"context"

	"github.com/rs/zerolog"
	options "github.com/vapusdata-ecosystem/vapusdata/core/options"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	ses "github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type AWSSesStore struct {
	Client       *ses.Client
	kmsSecretKey string
	logger       zerolog.Logger
}

// NewAwsSmClient creates a new instance of AwsSmStore using the provided AWSConfig options.
// It returns a pointer to AwsSmStore and an error, if any.
func NewAwsSesClient(ctx context.Context, opts *AWSConfig, logger zerolog.Logger) (*AWSSesStore, error) {
	configCl, err := opts.getAwsCLientConfig(ctx)
	if err != nil {
		return nil, dmerrors.DMError(ErrAwsConfigLoading, nil)
	}
	// mySession := session.Must(session.NewSession())
	return &AWSSesStore{
		logger: logger,
		Client: ses.New(ses.Options{
			Credentials:      configCl.Credentials,
			Logger:           configCl.Logger,
			Region:           configCl.Region,
			RetryMaxAttempts: configCl.RetryMaxAttempts,
		}, nil),
		kmsSecretKey: opts.KMSKey,
	}, nil
}

func (a *AWSSesStore) SendEmail(ctx context.Context, opts *options.SendEmailRequest, agentId string) error {
	var err error
	email := &ses.SendEmailInput{
		Source: aws.String(opts.From),
		Destination: &types.Destination{
			ToAddresses: opts.To,
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: aws.String(opts.Subject),
			},
			Body: &types.Body{
				Text: &types.Content{
					Data: aws.String(opts.Body),
				},
			},
		},
	}
	op, err := a.Client.SendEmail(ctx, email)
	if err != nil {
		a.logger.Error().Err(err).Msgf("Error sending email for agent %s", agentId)
	} else {
		a.logger.Info().Msgf("Email sent for agent %s with message id %s", agentId, *op.MessageId)
	}
	return err
}

func (a *AWSSesStore) SendRawEmail(ctx context.Context, opts *options.SendEmailRequest, agentId string) error {
	var err error
	email := &ses.SendRawEmailInput{
		Source:       aws.String(opts.From),
		Destinations: opts.To,
		RawMessage: &types.RawMessage{
			Data: opts.RawMessage,
		},
		Tags: func() []types.MessageTag {
			var tags []types.MessageTag
			for k, v := range opts.Tags {
				tags = append(tags, types.MessageTag{
					Name:  aws.String(k),
					Value: aws.String(v),
				})
			}
			return tags

		}(),
	}
	op, err := a.Client.SendRawEmail(ctx, email)
	if err != nil {
		a.logger.Error().Err(err).Msgf("Error sending email for agent %s", agentId)
	} else {
		a.logger.Info().Msgf("Email sent for agent %s with message id %s", agentId, *op.MessageId)
	}
	return err
}
