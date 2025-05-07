package gcp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/rs/zerolog"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
	sheets "google.golang.org/api/sheets/v4"
)

type SheetsOpts struct {
	Client *sheets.Service
	logger zerolog.Logger
}

func NewGSheets(ctx context.Context, opts *GcpConfig, userEmailAddr string, logger zerolog.Logger) (*SheetsOpts, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(string(opts.ServiceAccountKey))
	if err != nil {
		logger.Err(err).Msgf("Error while decoding the GCP KEY -- %v", err)
		return nil, err
	}
	creds, err := google.CredentialsFromJSON(ctx, decodedKey)
	if err != nil || creds == nil {
		logger.Err(err).Msgf("Error while creating credentials from json for GCP drive plugin-- %v", err)
		return nil, err
	}
	keyJson := map[string]any{}
	err = json.Unmarshal(creds.JSON, &keyJson)
	if err != nil {
		logger.Err(err).Msgf("Error while unmarshalling the GCP KEY json -- %v", err)
		return nil, err
	}
	clEmail, ok := keyJson["client_email"].(string)
	if !ok {
		logger.Err(err).Msgf("Error while getting the client_email from the GCP KEY json -- %v", err)
		return nil, err
	}
	log.Println("gcp-SVC-KEY-EMAIL", string(decodedKey))
	log.Println("gcp-SVC-KEY-EMAIL", clEmail)
	log.Println("gcp-SVC-KEY-EMAIL", userEmailAddr)
	log.Println("gcp-SVC-KEY-TokenSource", creds.TokenSource)

	// Now impersonate the ORGANIZATION user
	tokenSource, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
		TargetPrincipal: clEmail,
		Subject:         userEmailAddr,
		Scopes:          []string{"https://www.googleapis.com/auth/drive"},
	}, option.WithCredentialsJSON(decodedKey))
	if err != nil {
		logger.Err(err).Msgf("Error while impersonating the user -- %v", err)
		return nil, err
	}
	client, err := sheets.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		logger.Err(err).Msgf("Error while creating credentials from json -- %v", err)
		return nil, err
	}
	return &SheetsOpts{
		Client: client,
		logger: logger,
	}, nil
}
