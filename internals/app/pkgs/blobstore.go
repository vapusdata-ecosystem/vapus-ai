package apppkgs

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusdata/core/data-platform/connectors/blobs"
	models "github.com/vapusdata-ecosystem/vapusdata/core/models"
	dmlogger "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/logger"
)

type BlobStore struct {
	Path  string
	creds *models.DataSourceCredsParams
	blobs.BlobStore
}

func NewBlobStoreClient(ctx context.Context, creds *models.DataSourceCredsParams, log zerolog.Logger) (*BlobStore, error) {
	s := &BlobStore{}
	client, err := blobs.New(ctx, blobs.WithDataSourceCredsParams(creds), blobs.WithLogger(dmlogger.GetSubDMLogger(log, "datasvc", "blob_store_main")), blobs.WithDebug(true))
	if err != nil {
		return nil, err
	}
	s.creds = creds
	s.BlobStore = client
	return s, nil
}

func (s *BlobStore) GetCreds() *models.DataSourceCredsParams {
	return s.creds
}
