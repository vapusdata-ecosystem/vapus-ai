package appfuncs

import (
	"context"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	options "github.com/vapusdata-ecosystem/vapusdata/core/options"
	filetools "github.com/vapusdata-ecosystem/vapusdata/core/tools/files"
)

func GetMessageFileDataSet(ctx context.Context, file *mpb.FileData, bucket string, dmstores *apppkgs.BlobStore, logger zerolog.Logger) (*options.DataSetSummary, error) {
	data, err := dmstores.BlobStore.DownloadObject(ctx, &options.BlobOpsParams{
		BucketName: bucket,
		ObjectName: file.Name,
	})
	if err != nil {
		logger.Err(err).Msg("error while downloading data file")
		return nil, err
	}
	fType := filetools.GetConfFileType(file.Name)
	return filetools.FileDatasetLoader(data, fType, false)
}
