package apppkgs

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusdata/core/data-platform/connectors/databases"
	models "github.com/vapusdata-ecosystem/vapusdata/core/models"
	dmlogger "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/logger"
)

type DataStoreOptions struct {
	creds *models.DataSourceCredsParams
}

type dataStoreOpts func(*DataStoreOptions)

func WithDataStoreCreds(r *models.DataSourceCredsParams) dataStoreOpts {
	return func(s *DataStoreOptions) {
		s.creds = r
	}
}

func NewDataStore(ctx context.Context, log zerolog.Logger, opts ...dataStoreOpts) (*databases.DataStoreClient, error) {
	s := &DataStoreOptions{}
	for _, opt := range opts {
		opt(s)
	}
	return databases.New(ctx, databases.WithDataSourceCredsParams(s.creds), databases.WithLogger(dmlogger.GetSubDMLogger(log, "vapusloader Package", "DataStoreClient")))
}
