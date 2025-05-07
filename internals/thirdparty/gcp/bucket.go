package gcp

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"github.com/rs/zerolog"
	options "github.com/vapusdata-ecosystem/vapusdata/core/options"
	"google.golang.org/api/iterator"
	option "google.golang.org/api/option"
)

type GcpBucketClient struct {
	ProjectID string
	Client    *storage.Client
	logger    zerolog.Logger
}

//temporary function to satusfy interface
// func(s *GcpBucketClient) GetStorageData(ctx context.Context,params *options.BlobOpsParams) (*model.BlobStoreSchema,error) {

// 	blobStoreSchema := model.BlobStoreSchema{}

// 	return &blobStoreSchema , nil
// }

func NewBucketAgent(ctx context.Context, opts *GcpConfig, logger zerolog.Logger) (*GcpBucketClient, error) {
	opts.SetGcpProjectId(logger)
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(opts.ServiceAccountKey))
	if err != nil {
		fmt.Printf("Failed to create bucket client: %v\n", err)
		return nil, ErrCreatingBucketClient
	}
	return &GcpBucketClient{
		ProjectID: opts.ProjectID,
		Client:    client,
		logger:    logger,
	}, nil
}

func (g *GcpBucketClient) Close() {
	g.Client.Close()
}

func (g *GcpBucketClient) CreateBucket(ctx context.Context, params *options.BlobOpsParams) error {
	bucketName := params.BucketName
	bucket := g.Client.Bucket(bucketName)
	if err := bucket.Create(ctx, g.ProjectID, nil); err != nil {
		g.logger.Err(err).Msgf("Failed to create bucket: %v", err)
		return ErrCreatingBucket
	}
	return nil
}

func (g *GcpBucketClient) DeleteBucket(ctx context.Context, params *options.BlobOpsParams) error {
	bucketName := params.BucketName
	bucket := g.Client.Bucket(bucketName)
	if err := bucket.Delete(ctx); err != nil {
		g.logger.Err(err).Msgf("Failed to delete bucket: %v", err)
		return ErrDeletingBucket
	}
	return nil
}

func (g *GcpBucketClient) ListBuckets(ctx context.Context) ([]string, error) {
	buckets := g.Client.Buckets(ctx, g.ProjectID)
	var bucketNames []string
	for {
		bucketAttrs, err := buckets.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("Failed to list buckets: %v\n", err)
			return nil, err
		}
		bucketNames = append(bucketNames, bucketAttrs.Name)
	}
	return bucketNames, nil
}

func (g *GcpBucketClient) GetBucket(ctx context.Context, params *options.BlobOpsParams) (string, error) {
	bucketName := params.BucketName
	bucket := g.Client.Bucket(bucketName)
	_, err := bucket.Attrs(ctx)
	if err != nil {
		g.logger.Err(err).Msgf("Failed to get bucket: %v", err)
		return "", ErrGetBucket
	}
	return bucket.BucketName(), nil
}

func (g *GcpBucketClient) GetBucketAttrs(ctx context.Context, params *options.BlobOpsParams) (*storage.BucketAttrs, error) {
	bucketName := params.BucketName
	bucket := g.Client.Bucket(bucketName)
	attrs, err := bucket.Attrs(ctx)
	if err != nil {
		g.logger.Err(err).Msgf("Failed to get bucket attributes: %v", err)
		return nil, ErrGetBucketAttrs
	}
	return attrs, nil
}

func (g *GcpBucketClient) ListObjects(ctx context.Context, params *options.BlobOpsParams) ([]string, error) {
	bucketName := params.BucketName
	bucket := g.Client.Bucket(bucketName)
	objs := bucket.Objects(ctx, nil)
	var objNames []string
	for {
		objAttrs, err := objs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("Failed to list objects: %v\n", err)
			return nil, err
		}
		objNames = append(objNames, objAttrs.Name)
	}
	return objNames, nil
}

func (g *GcpBucketClient) GetObject(ctx context.Context, params *options.BlobOpsParams) (*storage.ObjectHandle, error) {
	bucketName := params.BucketName
	objectName := params.ObjectName
	bucket := g.Client.Bucket(bucketName)
	obj := bucket.Object(objectName)
	_, err := obj.Attrs(ctx)
	if err != nil {
		g.logger.Err(err).Msgf("Failed to get object: %v", err)
		return nil, ErrGettingObject
	}
	return obj, nil
}

func (g *GcpBucketClient) UploadObject(ctx context.Context, params *options.BlobOpsParams) error {

	bucketName := params.BucketName
	objectName := params.ObjectName
	data := params.Data
	bucket := g.Client.Bucket(bucketName)
	obj := bucket.Object(objectName)
	w := obj.NewWriter(ctx)
	if _, err := w.Write(data); err != nil {
		g.logger.Err(err).Msgf("Failed to write object: %v", err)
		return ErrUploadingObject
	}
	if err := w.Close(); err != nil {
		g.logger.Err(err).Msgf("Failed to close object: %v", err)
		return ErrUploadingObject
	}
	return nil
}

func (g *GcpBucketClient) DownloadObject(ctx context.Context, params *options.BlobOpsParams) ([]byte, error) {

	bucketName := params.BucketName
	objectName := params.ObjectName
	bucket := g.Client.Bucket(bucketName)
	obj := bucket.Object(objectName)
	r, err := obj.NewReader(ctx)
	if err != nil {
		g.logger.Err(err).Msgf("Failed to get object reader: %v", err)
		return nil, ErrDownloadingObject
	}
	defer r.Close()
	data, err := io.ReadAll(r)
	if err != nil {
		g.logger.Err(err).Msgf("Failed to read object: %v", err)
		return nil, ErrDownloadingObject
	}
	return data, nil
}

func (g *GcpBucketClient) DeleteObject(ctx context.Context, params *options.BlobOpsParams) error {

	bucketName := params.BucketName
	objectName := params.ObjectName
	bucket := g.Client.Bucket(bucketName)
	obj := bucket.Object(objectName)
	if err := obj.Delete(ctx); err != nil {
		g.logger.Err(err).Msgf("Failed to delete object: %v", err)
		return ErrDeletingObject
	}
	return nil
}
