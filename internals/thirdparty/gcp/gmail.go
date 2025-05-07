package gcp

import (
	"context"
	"encoding/base64"
	"log"
	"strings"

	"github.com/rs/zerolog"
	options "github.com/vapusdata-ecosystem/vapusdata/core/options"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	"google.golang.org/api/gmail/v1"
	gapi "google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type GmailClient struct {
	Client *gmail.Service
	logger zerolog.Logger
}

func NewGmailService(ctx context.Context, opts *GcpConfig, logger zerolog.Logger) (*GmailClient, error) {
	client, err := gmail.NewService(ctx, option.WithCredentialsJSON(opts.ServiceAccountKey))
	if err != nil {
		logger.Err(err).Msgf("Error while creating credentials from json for GMAIL -- %v", err)
		return nil, err
	}
	return &GmailClient{
		Client: client,
		logger: logger,
	}, nil
}

func NewGmailSource(ctx context.Context, opts *GcpConfig, logger zerolog.Logger) (*GmailClient, error) {
	client, err := gmail.NewService(ctx, option.WithCredentialsJSON(opts.ServiceAccountKey))
	if err != nil {
		logger.Err(err).Msgf("Error while creating credentials from json for GMAIL -- %v", err)
		return nil, err
	}
	return &GmailClient{
		Client: client,
		logger: logger,
	}, nil
}

func (x *GmailClient) SendEmail(ctx context.Context, opts *options.SendEmailRequest, agentId string) error {
	// Create the email message
	msg := &gmail.Message{
		ThreadId: "test_thread", // optional, depending on your use case
		Payload: &gmail.MessagePart{
			Body: &gmail.MessagePartBody{},
		},
	}
	bo := ""
	if opts.HtmlTemplateBody != "" {
		bo = opts.HtmlTemplateBody
	} else {
		bo = opts.Body
	}
	// Build the email content
	email := []byte("MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"From: " + opts.From + "\r\n" +
		"To: " + strings.Join(opts.To, ",") + "\r\n" +
		"Subject: " + opts.Subject + "\r\n\r\n" +
		bo)

	// Encode the email content as base64
	encodedMessage := base64.URLEncoding.EncodeToString(email)

	// Set the raw email content
	msg.Raw = encodedMessage

	// Send the email
	_, err := x.Client.Users.Messages.Send("me", msg).Do()
	return err
}

func (x *GmailClient) SendRawEmail(ctx context.Context, opts *options.SendEmailRequest, agentId string) error {
	msg := &gmail.Message{
		Payload: &gmail.MessagePart{
			Body: &gmail.MessagePartBody{},
		},
	}
	bo := ""
	if opts.HtmlTemplateBody != "" {
		bo = opts.HtmlTemplateBody
	} else {
		bo = opts.Body
	}
	// Build the email content
	email := []byte("MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"From: " + opts.From + "\r\n" +
		"To: " + strings.Join(opts.To, ",") + "\r\n" +
		"Subject: " + opts.Subject + "\r\n\r\n" +
		bo)

	// Encode the email content as base64
	encodedMessage := base64.StdEncoding.EncodeToString(email)

	// Set the raw email content
	msg.Raw = encodedMessage

	// Send the email
	att := opts.Attachments[0]
	decodedAttachment, err := base64.StdEncoding.DecodeString(string(att.Data))
	if err != nil {
		x.logger.Err(err).Msgf("Error while decoding attachment -- %v", err)
		return err
	}

	_, err = x.Client.Users.Messages.Send("me", msg).Media(strings.NewReader(string(decodedAttachment)),
		gapi.ContentType(att.ContentType),
	).Do()
	return err
}

func (x *GmailClient) Delete(ctx context.Context, opts *options.DeleteEmailRequest) error {
	_, err := x.Client.Users.Messages.Trash("me", opts.ID).Do()
	if err != nil {
		x.logger.Err(err).Msgf("Error while reading email -- %v", err)
		return err
	}
	return nil
}

func (x *GmailClient) Watch(ctx context.Context, opts *options.ListEmailsRequest) error {
	response, err := x.Client.Users.Watch("me", &gmail.WatchRequest{}).Do()
	if err != nil {
		x.logger.Err(err).Msgf("Error while reading email -- %v", err)
		return err
	}
	log.Println(response)
	return nil
}

func (x *GmailClient) Close() {
	if x.Client != nil {
		x.Client = nil
		dmutils.CleanPointers(x.Client)
	}
}
