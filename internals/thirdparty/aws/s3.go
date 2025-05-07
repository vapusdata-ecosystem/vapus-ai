package vaws

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/rs/zerolog"
	options "github.com/vapusdata-ecosystem/vapusdata/core/options"
)

type AwsS3BucketClient struct {
	S3Client  *s3.Client
	S3Manager *manager.Uploader
	logger    zerolog.Logger
}

func NewBucketAgent(ctx context.Context, opts *AWSConfig, logger zerolog.Logger) (*AwsS3BucketClient, error) {

	cfg, err := opts.getAwsCLientConfig(ctx)

	if err != nil {
		logger.Error().Err(err).Msg("Unable to load SDK config")
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	s3Manager := manager.NewUploader(s3Client)

	return &AwsS3BucketClient{
		S3Client:  s3Client,
		S3Manager: s3Manager,
		logger:    logger,
	}, nil
}

func (s *AwsS3BucketClient) CreateBucket(ctx context.Context, params *options.BlobOpsParams) error {

	bucketName := params.BucketName
	region := params.Region

	_, err := s.S3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})
	if err != nil {
		var owned *types.BucketAlreadyOwnedByYou
		var exists *types.BucketAlreadyExists
		if errors.As(err, &owned) {
			log.Printf("You already own bucket %s.\n", bucketName)
			err = owned
		} else if errors.As(err, &exists) {
			log.Printf("Bucket %s already exists.\n", bucketName)
			err = exists
		}
	} else {
		err = s3.NewBucketExistsWaiter(s.S3Client).Wait(
			ctx, &s3.HeadBucketInput{Bucket: aws.String(bucketName)}, time.Minute)
		if err != nil {
			log.Printf("Failed attempt to wait for bucket %s to exist.\n", bucketName)
		}
	}
	return err
}

func (s *AwsS3BucketClient) DeleteBucket(ctx context.Context, params *options.BlobOpsParams) error {

	bucketName := params.BucketName
	_, err := s.S3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName)})
	if err != nil {
		var noBucket *types.NoSuchBucket
		if errors.As(err, &noBucket) {
			log.Printf("Bucket %s does not exist.\n", bucketName)
			err = noBucket
		} else {
			log.Printf("Couldn't delete bucket %v. Here's why: %v\n", bucketName, err)
		}
	} else {
		err = s3.NewBucketNotExistsWaiter(s.S3Client).Wait(
			ctx, &s3.HeadBucketInput{Bucket: aws.String(bucketName)}, time.Minute)
		if err != nil {
			log.Printf("Failed attempt to wait for bucket %s to be deleted.\n", bucketName)
		} else {
			log.Printf("Deleted %s.\n", bucketName)
		}
	}
	return err
}

func (s *AwsS3BucketClient) ListBuckets(ctx context.Context) ([]string, error) {
	var err error
	var output *s3.ListBucketsOutput
	var bucketNames []string

	bucketPaginator := s3.NewListBucketsPaginator(s.S3Client, &s3.ListBucketsInput{})
	for bucketPaginator.HasMorePages() {
		output, err = bucketPaginator.NextPage(ctx)
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "AccessDenied" {
				fmt.Println("You don't have permission to list buckets for this account.")
				err = apiErr
			} else {
				log.Printf("Couldn't list buckets for your account. Here's why: %v\n", err)
			}
			break
		}
		for _, bucket := range output.Buckets {
			if bucket.Name != nil {
				bucketNames = append(bucketNames, *bucket.Name)
			}
		}
	}

	return bucketNames, err
}

func (s *AwsS3BucketClient) GetBucket(ctx context.Context, params *options.BlobOpsParams) (string, error) {

	bucketName := params.BucketName
	_, err := s.S3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	exists := true
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				log.Printf("Bucket %v is available.\n", bucketName)
				exists = false
				err = nil
			default:
				log.Printf("Either you don't have access to bucket %v or another error occurred. "+
					"Here's what happened: %v\n", bucketName, err)
			}
		}
	} else {
		log.Printf("Bucket %v exists and you already own it.", bucketName)
	}

	if exists {
		return bucketName, nil
	} else {
		return "", err
	}

	// return bucketName, err
}

func (s *AwsS3BucketClient) ListObjects(ctx context.Context, params *options.BlobOpsParams) ([]string, error) {

	bucketName := params.BucketName
	var err error
	var output *s3.ListObjectsV2Output
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	}
	var objectKeys []string
	objectPaginator := s3.NewListObjectsV2Paginator(s.S3Client, input)
	for objectPaginator.HasMorePages() {
		output, err = objectPaginator.NextPage(ctx)
		for _, prefixes := range output.CommonPrefixes {

			fmt.Println(*prefixes.Prefix)
		}
		if err != nil {
			var noBucket *types.NoSuchBucket
			if errors.As(err, &noBucket) {
				log.Printf("Bucket %s does not exist.\n", bucketName)
				err = noBucket
			}
			break
		}
		// Extract object keys from the output
		for _, object := range output.Contents {
			if object.Key != nil {
				objectKeys = append(objectKeys, *object.Key)
			}
		}
	}
	return objectKeys, err
}

func (s *AwsS3BucketClient) UploadObject(ctx context.Context, params *options.BlobOpsParams) error {

	bucketName := params.BucketName
	objectName := params.ObjectName
	data := params.Data
	body := bytes.NewReader(data)

	_, err := s.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
		Body:   body,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
			log.Printf("Error while uploading object to %s. The object is too large.\n"+
				"To upload objects larger than 5GB, use the S3 console (160GB max)\n"+
				"or the multipart upload API (5TB max).", bucketName)
		} else {
			log.Printf("Couldn't upload object to %v:%v. Here's why: %v\n", bucketName, objectName, err)
		}
		return err
	}

	err = s3.NewObjectExistsWaiter(s.S3Client).Wait(
		ctx, &s3.HeadObjectInput{Bucket: aws.String(bucketName), Key: aws.String(objectName)}, time.Minute)
	if err != nil {
		log.Printf("Failed attempt to wait for object %s to exist.\n", objectName)
	}

	return nil
}

func (s *AwsS3BucketClient) DownloadObject(ctx context.Context, params *options.BlobOpsParams) ([]byte, error) {

	bucketName := params.BucketName
	objectName := params.ObjectName
	result, err := s.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		var noKey *types.NoSuchKey
		if errors.As(err, &noKey) {
			log.Printf("Can't get object %s from bucket %s. No such key exists.\n", objectName, bucketName)
			err = noKey
		} else {
			log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectName, err)
		}
		return nil, err
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", objectName, err)
		return nil, err
	}

	return body, nil
}
func (s *AwsS3BucketClient) DeleteObject(ctx context.Context, params *options.BlobOpsParams) error {

	bucketName := params.BucketName
	objectName := params.ObjectName
	versionId := params.ObjectVersionId
	bypassGovernance := params.ByPassGovernance

	deleted := false
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}
	if versionId != "" {
		input.VersionId = aws.String(versionId)
	}
	if bypassGovernance {
		input.BypassGovernanceRetention = aws.Bool(true)
	}
	_, err := s.S3Client.DeleteObject(ctx, input)
	if err != nil {
		var noKey *types.NoSuchKey
		var apiErr *smithy.GenericAPIError
		if errors.As(err, &noKey) {
			log.Printf("Object %s does not exist in %s.\n", objectName, bucketName)
			err = noKey
		} else if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case "AccessDenied":
				log.Printf("Access denied: cannot delete object %s from %s.\n", objectName, bucketName)
				err = nil
			case "InvalidArgument":
				if bypassGovernance {
					log.Printf("You cannot specify bypass governance on a bucket without lock enabled.")
					err = nil
				}
			}
		}
	} else {
		err = s3.NewObjectNotExistsWaiter(s.S3Client).Wait(
			ctx, &s3.HeadObjectInput{Bucket: aws.String(bucketName), Key: aws.String(objectName)}, time.Minute)
		if err != nil {
			log.Printf("Failed attempt to wait for object %s in bucket %s to be deleted.\n", objectName, bucketName)
		} else {
			deleted = true
		}
	}

	if deleted {
		return nil
	}
	return err
}

func (s *AwsS3BucketClient) Close() {
}
