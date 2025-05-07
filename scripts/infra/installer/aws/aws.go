package create_aws

import (
	"context"
	"fmt"
	"log"

	"github.com/rs/zerolog"
	aws_config "github.com/vapusdata-ecosystem/vapusai/core/thirdparty/aws"
	vapusmodel "github.com/vapusdata-ecosystem/vapusai/scripts/goscripts/installer"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	// vapusmodel "github.com/vapusdata-ecosystem/vapusai/infra/installer"
)

type AWSInstaller struct {
	Config    aws.Config
	RDSClient *rds.Client
	S3Client  *s3.Client
}

func AWSConfigClient(ctx context.Context, opts *vapusmodel.AWSConfig, logger zerolog.Logger) (*AWSInstaller, error) {
	awsCfg, err := aws_config.GetAwsCLientConfig(ctx, &aws_config.AWSConfig{
		Region:          opts.Region,
		AccessKeyId:     opts.AccessKey,
		SecretAccessKey: opts.SecretKey,
	})
	if err != nil {
		logger.Err(err).Msg("Error creating AWS client")
		return nil, err
	}
	return &AWSInstaller{
		Config: awsCfg,
	}, nil
}

func (la *AWSInstaller) GetRDSClient(ctx context.Context) *AWSInstaller {
	rdsClient := rds.NewFromConfig(la.Config)
	return &AWSInstaller{
		RDSClient: rdsClient,
	}
}

func (la *AWSInstaller) AWSCreateDb(ctx context.Context, creds *vapusmodel.SetupPostgressInstance, logger zerolog.Logger) (*AWSInstaller, error) {

	input := &rds.CreateDBInstanceInput{
		DBInstanceIdentifier: aws.String(creds.InstanceID),
		DBName:               aws.String(creds.DBName),
		Engine:               aws.String("postgres"),
		DBInstanceClass:      aws.String(creds.DBInstanceClass), // example instance class "db.t3.micro"
		AllocatedStorage:     aws.Int32(creds.AllocatedStorage), // 20 GB
		MasterUsername:       aws.String(creds.Username),
		MasterUserPassword:   aws.String(creds.Password),
		// DBSubnetGroupName:    aws.String(creds.DBSubnetGroupName),
		PubliclyAccessible: aws.Bool(creds.PubliclyAccessible),
	}

	out, err := la.RDSClient.CreateDBInstance(ctx, input)
	if err != nil {
		logger.Err(err).Msg("failed to create RDS instance:")
		return nil, err
	}

	log.Printf("RDS Postgres creation initiated: %v (status: %v)\n",
		aws.ToString(out.DBInstance.DBInstanceIdentifier),
		out.DBInstance.DBInstanceStatus)
	return &AWSInstaller{
		RDSClient: la.RDSClient,
	}, nil
}

func (la *AWSInstaller) DeleteInstance(ctx context.Context, instanceName string) error {
	// fmt.Println("RDS Client", la.RDSClient)
	// rdsClient := rds.NewFromConfig(la.Config)
	deleteOutput, err := la.RDSClient.DeleteDBInstance(ctx, &rds.DeleteDBInstanceInput{
		DBInstanceIdentifier:   aws.String(instanceName),
		SkipFinalSnapshot:      aws.Bool(true),
		DeleteAutomatedBackups: aws.Bool(true),
	})
	if err != nil {
		log.Printf("Couldn't delete instance %v: %v\n", instanceName, err)
		return err
	}
	fmt.Println("Deletion Metadata: ", deleteOutput.ResultMetadata)
	return nil
}

func (la *AWSInstaller) GetDBInstance(ctx context.Context, instanceName string) error {
	output, err := la.RDSClient.DescribeDBInstances(ctx,
		&rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: aws.String(instanceName),
		})
	if err != nil {
		log.Printf("Couldn't get instance %v: %v\n", instanceName, err)
		return err
	} else {
		for _, val := range output.DBInstances {
			fmt.Println("*val.DBInstanceStatus: ", *val.DBInstanceStatus)
			fmt.Println("*val.DBInstanceClass: ", *val.DBInstanceClass)
			fmt.Println("*val.DBName: ", *val.DBName)
			fmt.Println("")
		}
		return nil
	}
}

func (la *AWSInstaller) GetBucketClient(ctx context.Context) *AWSInstaller {
	s3Client := s3.NewFromConfig(la.Config)
	return &AWSInstaller{
		S3Client: s3Client,
	}
}

func (la *AWSInstaller) CreateBucket(ctx context.Context, params *vapusmodel.SetupBucketInput, logger zerolog.Logger) error {
	var AwsS3Client aws_config.AwsS3BucketClient
	AwsS3Client.S3Client = la.S3Client

	err := AwsS3Client.CreateBucket(ctx, &options.BlobOpsParams{
		BucketName: params.BucketName,
		Region:     params.Region,
	})
	if err != nil {
		logger.Err(err).Msg("failed to Delete S3 Bucket:")
		return err
	}
	return nil
}

func (la *AWSInstaller) DeleteBucket(ctx context.Context, bucketName string, logger zerolog.Logger) error {
	var AwsS3Client aws_config.AwsS3BucketClient
	AwsS3Client.S3Client = la.S3Client
	err := AwsS3Client.DeleteBucket(ctx, &options.BlobOpsParams{
		BucketName: bucketName,
	})
	if err != nil {
		logger.Err(err).Msg("failed to Delete S3 Bucket:")
		return err
	}

	return nil
}
