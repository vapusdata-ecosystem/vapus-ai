package sendgrid

import (
	"context"
	"encoding/base64"
	"log"

	"github.com/rs/zerolog"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	options "github.com/vapusdata-ecosystem/vapusdata/core/options"
	dmlogger "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/logger"
)

type SendgridClient struct {
	client *sendgrid.Client
	logger zerolog.Logger
	Params []*models.Mapper
	host   string
	creds  *models.PluginNetworkParams
}

type Options func(*SendgridClient)

func WithParams(params []*models.Mapper) Options {
	return func(c *SendgridClient) {
		c.Params = params
	}
}

func WithCreds(creds *models.PluginNetworkParams) Options {
	return func(c *SendgridClient) {
		c.creds = creds
	}
}

func WithLogger(logger zerolog.Logger) Options {
	return func(c *SendgridClient) {
		c.logger = dmlogger.GetSubDMLogger(logger, "emailer", "sendgrid")
	}
}

func New(ctx context.Context, opts ...Options) (*SendgridClient, error) {
	sndg := &SendgridClient{}
	for _, opt := range opts {
		opt(sndg)
	}
	if sndg.creds == nil {
		sndg.logger.Error().Msg("Missing credentials")
		return nil, apperr.ErrMissingCredentials
	}
	if sndg.creds.Credentials.GetApiToken() == "" {
		sndg.logger.Error().Msg("Missing API key")
		return nil, apperr.ErrMissingCredentials
	}
	cl := sendgrid.NewSendClient(sndg.creds.Credentials.GetApiToken())
	if sndg.host != "" {
		cl.Request.BaseURL = sndg.creds.URL
	} else {
		cl.Request.BaseURL = "https://api.sendgrid.com"
	}
	sndg.client = cl
	return sndg, nil
}

func (s *SendgridClient) GetParams(pType string) string {
	for _, param := range s.Params {
		log.Println(param.Key, "====================+++++++++++++????", param.Value, "====================+++++++++++++>>")
		if param.Key == pType {
			log.Println(param.Value, "====================+++++++++++++")
			return param.Value
		}
	}
	return ""
}

func (s *SendgridClient) SendEmail(ctx context.Context, opts *options.SendEmailRequest, agentId string) error {
	var err error
	message := mail.NewSingleEmail(mail.NewEmail(opts.SenderName, opts.From),
		opts.Subject,
		mail.NewEmail("", opts.To[0]),
		opts.Body,
		opts.HtmlTemplateBody)
	client := sendgrid.NewSendClient(s.creds.Credentials.GetApiToken())
	_, err = client.Send(message)
	if err != nil {
		s.logger.Error().Err(err).Msgf("Error sending email for agent %s", agentId)
	} else {
		s.logger.Info().Msgf("Email sent for agent %s", agentId)
	}
	return err
}

func (s *SendgridClient) SendRawEmail(ctx context.Context, opts *options.SendEmailRequest, agentId string) error {
	var err error
	log.Println(s.GetParams(mpb.EmailSettings_SENDER_NAME.String()), s.GetParams(mpb.EmailSettings_SENDER_EMAIL.String()), "====================+++++++++++++")
	message := mail.NewV3Mail()
	message.SetFrom(mail.NewEmail(s.GetParams(mpb.EmailSettings_SENDER_NAME.String()), s.GetParams(mpb.EmailSettings_SENDER_EMAIL.String())))
	message.Subject = opts.Subject
	personalization := mail.NewPersonalization()
	ttos := func() []*mail.Email {
		var tos []*mail.Email
		for _, to := range opts.To {
			tos = append(tos, mail.NewEmail("", to))
		}
		return tos
	}
	personalization.AddTos(ttos()...)
	personalization.Subject = opts.Subject
	log.Println(personalization, "====================+++++++++++++", opts.Attachments, "====================+++++++++++++")
	message.AddPersonalizations(personalization)
	if len(opts.Attachments) > 0 {
		if opts.HtmlTemplateBody != "" {
			content := mail.NewContent("text/html", opts.HtmlTemplateBody)
			message.AddContent(content)
		} else {
			content := mail.NewContent("text/plain", opts.Body)
			message.AddContent(content)
		}
		for _, attach := range opts.Attachments {
			attachmentData := base64.StdEncoding.EncodeToString(attach.Data)
			log.Println(attachmentData, "====================+++++++++++++")
			att := mail.NewAttachment()
			att.SetContent(attachmentData)
			att.Type = attach.ContentType
			att.Filename = attach.FileName
			att.SetDisposition("attachment")
			message.AddAttachment(att)
		}
	}
	request := sendgrid.GetRequest(s.creds.Credentials.ApiToken, "/v3/mail/send", s.host)
	request.Method = "POST"
	request.Body = mail.GetRequestBody(message)
	response, err := sendgrid.API(request)
	if response != nil {
		log.Println(response.StatusCode, response.Body, response.Headers, "====================+++++++++++++")
	}
	if err != nil {
		s.logger.Error().Err(err).Msgf("Error sending email for agent %s", agentId)
	} else {
		s.logger.Info().Msgf("Email sent for agent %s with message id %s", agentId, response.Body)
	}
	return err
}
