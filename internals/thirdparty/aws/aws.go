package vaws

import (
	"context"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awscreds "github.com/aws/aws-sdk-go-v2/credentials"
	bedrockService "github.com/aws/aws-sdk-go-v2/service/bedrock"
	bedrockRuntime "github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
)

type AWSConfig struct {
	SecretAccessKey, AccessKeyId, Region string
	KMSKey                               string
}

type BedrockConfig struct {
	AwsConfig      aws.Config
	BedrockService *bedrockService.Client
	BedrockRuntime *bedrockRuntime.Client
}

// getAwsCLientConfig returns the AWS configuration for the given AWSConfig.
// It validates the required fields (Region, AccessKeyId, SecretAccessKey) and returns the AWS configuration.
func (a *AWSConfig) getAwsCLientConfig(ctx context.Context) (aws.Config, error) {
	if a.Region == "" || a.AccessKeyId == "" || a.SecretAccessKey == "" {
		return aws.Config{}, dmerrors.DMError(ErrInvalidAwsConfig, nil)
	}
	return awsConfig.LoadDefaultConfig(
		ctx,
		awsConfig.WithRegion(a.Region),
		awsConfig.WithCredentialsProvider(
			awscreds.NewStaticCredentialsProvider(a.AccessKeyId, a.SecretAccessKey, ""),
		),
	)
}

func GetAwsCLientConfig(ctx context.Context, a *AWSConfig) (aws.Config, error) {
	return a.getAwsCLientConfig(ctx)
}

func GetBedrockConfig(ctx context.Context, a *AWSConfig) (*BedrockConfig, error) {
	cfg, err := a.getAwsCLientConfig(ctx)
	if err != nil {
		return nil, err
	}
	genAiCl := bedrockRuntime.NewFromConfig(cfg)
	service := bedrockService.NewFromConfig(cfg)
	return &BedrockConfig{
		AwsConfig:      cfg,
		BedrockRuntime: genAiCl,
		BedrockService: service,
	}, nil
}
