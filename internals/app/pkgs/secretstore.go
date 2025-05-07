package apppkgs

import (
	"context"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	secretstore "github.com/vapusdata-ecosystem/vapusdata/core/data-platform/connectors/secrets-stores"
	models "github.com/vapusdata-ecosystem/vapusdata/core/models"
	dmlogger "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/logger"
	"gopkg.in/yaml.v2"
)

type SecretStore struct {
	Path  string
	creds *models.DataSourceCredsParams
	secretstore.SecretStore
}

type secretOpts func(*SecretStore)

func WithSecretStorePath(r string) secretOpts {
	return func(s *SecretStore) {
		s.Path = r
	}
}

func NewSecretStoreCreds(filePath string, log zerolog.Logger) (*models.DataSourceCredsParams, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal().Msgf("Error reading file: %v", err)
		return nil, err
	}
	conf := &models.DataSourceCredsParams{}
	log.Info().Msgf("File Content after reading : %s", string(bytes))
	err = yaml.Unmarshal(bytes, conf)
	if err != nil {
		log.Fatal().Msgf("Error unmarshalling file: %v", err)
		return nil, err
	}
	err = validator.New().Struct(conf)

	if err != nil {
		return nil, err
	}
	return conf, nil
}

func NewSecretStoreClient(ctx context.Context, log zerolog.Logger, opts ...secretOpts) (*SecretStore, error) {
	s := &SecretStore{}
	for _, opt := range opts {
		opt(s)
	}
	creds, err := NewSecretStoreCreds(s.Path, log)
	if err != nil {
		return nil, err
	}
	client, err := secretstore.New(ctx, secretstore.WithDataSourceCredsParams(creds), secretstore.WithLogger(dmlogger.GetSubDMLogger(log, "datasvc", "secret_store_main")), secretstore.WithDebug(true))
	if err != nil {
		return nil, err
	}
	s.creds = creds
	s.SecretStore = client
	return s, nil
}

func (s *SecretStore) GetCreds() *models.DataSourceCredsParams {
	return s.creds
}

func (s *SecretStore) Close() {
	s.SecretStore.Close()
}
