package datarepo

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/bytedance/sonic"
	"github.com/databricks/databricks-sql-go/logger"
	"github.com/rs/zerolog"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	options "github.com/vapusdata-ecosystem/vapusai/core/options"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
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
		fmt.Println("obj.Checksums: ", obj.Checksums)
		fmt.Println("Checksums: ", checksum)
		if slices.Contains(obj.Checksums, checksum) {
			fmt.Println("Checksum exsits ========")
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

func GetFile(ctx context.Context, dmstores *apppkgs.VapusStore, filePath string, ctxClaim map[string]string) (*models.FileStoreLog, error) {
	if filePath == "" {
		return nil, apperr.ErrAIModelNode404
	}
	result := &models.FileStoreLog{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE path = '%s' AND created_by = '%s'", apppkgs.FileStoreLogTable, filePath, ctxClaim[encryption.ClaimUserIdKey])
	err := dmstores.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while getting the file from datastore")
		return nil, err
	}
	return result, nil
}

func DeleteRedisKey(ctx context.Context, keyPath string, dmstores *apppkgs.VapusStore, logger zerolog.Logger) error {
	fmt.Println("Error while uploading the data....", keyPath)
	_, err := dmstores.Cacher.RedisClient.DeleteKey(ctx, keyPath) // Deleting the Key from the redis
	if err != nil {
		logger.Err(err).Msgf("Error while deleting the Redis chache")
		return err
	}
	return nil
}
