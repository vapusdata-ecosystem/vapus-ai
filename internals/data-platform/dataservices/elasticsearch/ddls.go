package elasticsearch

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/pkgs"
)

func (ds *ElasticSearchStore) CreateIndex(ctx context.Context, index string) error {
	if exists, err := ds.TClient.Indices.Exists(index).Do(ctx); exists || err != nil {
		ds.logger.Debug().Msgf("Index %v already exists.", index)
		return err
	}
	_, err := ds.TClient.Indices.Create(index).Do(ctx)
	if err != nil {
		ds.logger.Fatal().Err(err).Ctx(ctx).Msg("error while creating index in elastic search")
		return err
	}
	ds.logger.Debug().Msgf("Index %v created successfully.", index)
	return nil
}

func (ds *ElasticSearchStore) CreateIndexWithMapping(ctx context.Context, opts *pkgs.DataTablesOpts) error {
	if exists, err := ds.TClient.Indices.Exists(opts.Name).Do(ctx); exists || err != nil {
		ds.logger.Debug().Msgf("Index %v already exists.", opts.Name)
		return err
	}
	mappingJSON, err := json.Marshal(opts.MapsScheme)
	if err != nil {
		return err
	}
	createReq := esapi.IndicesCreateRequest{
		Index: opts.Name,
		Body:  strings.NewReader(string(mappingJSON)),
	}
	res, err := createReq.Do(context.Background(), ds.TClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return err
	}
	return nil
}
