package datarepo

import (
	"context"
	"path/filepath"
	"slices"

	"github.com/bytedance/sonic"
	"github.com/databricks/databricks-sql-go/logger"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	options "github.com/vapusdata-ecosystem/vapusai/core/options"
)

func GetOrCreateBucket(ctx context.Context, dmstores *apppkgs.VapusStore, createIfNotExist bool, bucketName string) (string, error) {
	name, err := dmstores.BlobStore.GetBucket(ctx, &options.BlobOpsParams{
		BucketName: bucketName,
	})
	if err != nil || name == "" {
		if !createIfNotExist {
			return "", err
		} else {
			err = dmstores.BlobStore.CreateBucket(ctx, &options.BlobOpsParams{
				BucketName: bucketName,
			})
			if err != nil {
				return "", err
			}
		}
	}
	return name, nil
}

func LogFileStoreLog(ctx context.Context, dmstores *apppkgs.VapusStore, fileStoreLog *models.FileStoreLog, ctxClaim map[string]string) error {
	fileStoreLog.PreSaveCreate(ctxClaim)
	if fileStoreLog == nil {
		return nil
	}
	_, err := dmstores.Db.PostgresClient.DB.NewInsert().Model(fileStoreLog).ModelTableExpr(apppkgs.FileStoreLogTable).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving file store log in datastore")
		return err
	}
	return nil
}

func LogFileCacher(ctx context.Context, dmstores *apppkgs.VapusStore, counter int64, checksum, name, path string) error {
	obj := &models.FileStoreCache{
		Checksums: []string{checksum},
		Name:      name,
		Path:      path,
		Counter:   counter,
	}
	bbytes, err := sonic.Marshal(obj)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while marshalling file cache object")
		return err
	}
	_, err = dmstores.Cacher.RedisClient.WriteKV(ctx, filepath.Join(path, name), string(bbytes))
	return err
}

func ValidateFileCache(ctx context.Context, dmstores *apppkgs.VapusStore, checksum, name, path string) (int, bool, error) {
	val, err := dmstores.Cacher.RedisClient.ReadKV(ctx, filepath.Join(path, name))
	if err != nil || val == nil {
		logger.Err(err).Ctx(ctx).Msg("error while getting file cache from redis, key not found")
		return 0, false, LogFileCacher(ctx, dmstores, 0, checksum, name, path)
	}
	nVal, ok := val.(string)
	if !ok {
		logger.Err(err).Ctx(ctx).Msg("error while getting file cache from redis, value not a string")
		return 0, false, LogFileCacher(ctx, dmstores, 0, checksum, name, path)
	}
	obj := &models.FileStoreCache{}
	err = sonic.Unmarshal([]byte(nVal), obj)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while unmarshalling file cache from redis")
		return 0, false, err
	}
	if obj.Checksums != nil {
		if slices.Contains(obj.Checksums, checksum) {
			return int(obj.Counter), true, nil
		} else {
			obj.Counter += 1
			LogFileCacher(ctx, dmstores, obj.Counter, checksum, name, path)
			return int(obj.Counter), false, nil
		}
	} else {
		return 0, false, LogFileCacher(ctx, dmstores, 0, checksum, name, path)
	}
}
