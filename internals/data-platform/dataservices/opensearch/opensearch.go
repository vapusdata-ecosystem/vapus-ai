package opensearch

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	"github.com/rs/zerolog"
	// "github.com/aws/aws-sdk-go-v2/config"
	// "github.com/opensearch-project/opensearch-go/v2"
	// "github.com/aws/aws-sdk-go-v2/credentials"
	// requestsigner "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	// aws_config "github.com/vapusdata-ecosystem/vapusai/core/thirdparty/aws"
)

type OpenSearchOpts struct {
	Endpoint  string // URL
	Region    string // DataSourceCreds => GenericCredentialModel => AWSCreds => Region
	AccessKey string // DataSourceCreds => GenericCredentialModel => AWSCreds => AccessKeyId
	SecretKey string // DataSourceCreds => GenericCredentialModel => AWSCreds => SecretAccessKey
	KMSKey    string // DataSourceCreds => GenericCredentialModel => AWSCreds => RoleArn
	Username  string
	Password  string
}

type OpenSearchStore struct {
	Opts   *OpenSearchOpts
	logger zerolog.Logger
	Client *opensearchapi.Client
	// AwsConfig aws.Config
	// Request   *http.Request
	// Signer    *requestsigner.Signer
	// AOS      *http.Client
	// Response *http.Response
	// AwsCred  aws.Credentials
}

func NewOpenSearchStore(opts *OpenSearchOpts, l zerolog.Logger) (*OpenSearchStore, error) {
	client, err := opensearchapi.NewClient(
		opensearchapi.Config{
			Client: opensearch.Config{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // For testing only. Use certificate for validation.
				},
				Addresses: []string{opts.Endpoint},
				Username:  opts.Username, // For testing only. Don't store credentials in code.
				Password:  opts.Password,
			},
		},
	)
	fmt.Println(client)
	if err != nil {
		l.Err(err).Msg("Error creating AWS client")
		return nil, err
	}

	l.Info().Msg("Connected to Open Search database successfully")

	return &OpenSearchStore{
		Opts:   opts,
		Client: client,
		logger: l,
	}, nil

}
func (m *OpenSearchStore) Close() {
	// m.Response.Body.Close()
	// m.Client.Close()
}

// func NewOpenSearchStore(opts *OpenSearchOpts, l zerolog.Logger) (*OpenSearchStore, error) {

// 	ctx := context.Background()

// 	awsCfg, err := aws_config.GetAwsCLientConfig(ctx, &aws_config.AWSConfig{
// 		KMSKey:          opts.KMSKey,
// 		Region:          opts.Region,
// 		AccessKeyId:     opts.AccessKey,
// 		SecretAccessKey: opts.SecretKey,
// 	})
// 	if err != nil {
// 		l.Err(err).Msg("Error creating AWS client")
// 		return nil, err
// 	}
// 	// credentials := credentials.NewStaticCredentialsProvider()
// 	signer := requestsigner.NewSigner()
// 	req, _ := http.NewRequest("GET", opts.Endpoint, nil)
// 	req.Header.Add("Content-Type", "application/json")

// 	creds, err := awsCfg.Credentials.Retrieve(context.Background())
// 	if err != nil {
// 		l.Err(err).Msg("Error getting AWS credential"	)
// 		return nil, err
// 	}

// 	// Sign the request
// 	err = signer.SignHTTP(ctx, creds, req, "", "es", awsCfg.Region, time.Now())
// 	if err != nil {
// 		l.Err(err).Msg("failed to sign request:")
// 		return nil, err
// 	}

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Fatalf("request failed: %v", err)
// 		l.Err(err).Msg("Request Failed: ")
// 		return nil, err
// 	}

// 	// Insert the data
// 	// defer resp.Body.Close()

// 	// Read response
// 	body, _ := io.ReadAll(resp.Body)
// 	fmt.Println(string(body))
// 	l.Info().Msg("Connected to Open Search database successfully")

// 	return &OpenSearchStore{
// 		Opts:      opts,
// 		AwsConfig: awsCfg,
// 		Request:   req,
// 		Signer:    signer,
// 		logger:    l,
// 		AOS:       client,
// 		Response:  resp,
// 		AwsCred:   creds,
// 	}, nil
// }

// func getCredentialProvider(accessKey, secretAccessKey, token string) aws.CredentialsProviderFunc {
// 	return func(ctx context.Context) (aws.Credentials, error) {
// 		c := &aws.Credentials{
// 			AccessKeyID:     accessKey,
// 			SecretAccessKey: secretAccessKey,
// 			SessionToken:    token,
// 		}
// 		return *c, nil
// 	}
// }
