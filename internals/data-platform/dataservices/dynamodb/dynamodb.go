package dynamodb

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBOpts struct {
	Region    string // DataSourceCreds => GenericCredentialModel => AWSCreds => Region
	AccessKey string // DataSourceCreds => GenericCredentialModel => AWSCreds => AccessKeyId
	SecretKey string // DataSourceCreds => GenericCredentialModel => AWSCreds => SecretAccessKey
}

type DynamoDBStore struct {
	Opts      *DynamoDBOpts
	logger    zerolog.Logger
	Client    *dynamodb.Client
	AwsConfig aws.Config
}

func NewDynamoDBStore(opts *DynamoDBOpts, l zerolog.Logger) (*DynamoDBStore, error) {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(opts.Region),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     opts.AccessKey,
				SecretAccessKey: opts.SecretKey,
			},
		}),
	)
	if err != nil {
		l.Err(err).Msg("Error creating AWS client: 'Invalid Access Key or Secret AccessKey' ")
		return nil, err
	}
	dynamodbClient := dynamodb.NewFromConfig(awsCfg, func(opt *dynamodb.Options) {
		opt.Region = opts.Region
	})
	return &DynamoDBStore{
		Opts:      opts,
		AwsConfig: awsCfg,
		logger:    l,
		Client:    dynamodbClient,
	}, nil

}
func (m *DynamoDBStore) Close() {
	// m.Response.Body.Close()
	// m.Client.Close()
}
