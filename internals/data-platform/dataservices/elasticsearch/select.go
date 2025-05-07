package elasticsearch

import (
	"context"

	datasvcpkgs "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/pkgs"
)

func (ds *ElasticSearchStore) Count(ctx context.Context, queryOpts *datasvcpkgs.QueryOpts) (int64, error) {
	count, err := ds.TClient.Count().Index(queryOpts.DataCollection).Do(ctx)
	if err != nil {
		return 0, err
	}
	return count.Count, nil
}
