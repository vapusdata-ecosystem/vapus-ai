package opensearch

// import (
// 	"context"
// 	"crypto/tls"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
// 	"github.com/aws/aws-sdk-go-v2/config"
// 	"github.com/opensearch-project/opensearch-go"
// 	"github.com/opensearch-project/opensearch-go/opensearchtransport"
// )

// type signingRoundTripper struct {
// 	next   http.RoundTripper
// 	signer *v4.Signer
// 	cfg    aws.Config
// }

// func (srt *signingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
// 	creds, err := srt.cfg.Credentials.Retrieve(context.Background()) // Retrieve credentials
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = srt.signer.SignHTTP(context.Background(), creds, req, "", "es", srt.cfg.Region, time.Now())
// 	if err != nil {
// 		return nil, err
// 	}
// 	return srt.next.RoundTrip(req)
// }

// func main() {
// 	cfg, err := config.LoadDefaultConfig(context.Background())
// 	if err != nil {
// 		log.Fatalf("failed to load AWS config: %v", err)
// 	}

// 	endpoint := "https://your-opensearch-ORGANIZATION.us-east-1.es.amazonaws.com"

// 	transport := &http.Transport{
// 		TLSClientConfig: &tls.Config{
// 			// For production, ensure InsecureSkipVerify is false.
// 			InsecureSkipVerify: false,
// 		},
// 	}
// 	signer := v4.NewSigner()

// 	signingTransport := &signingRoundTripper{
// 		next:   transport,
// 		signer: signer,
// 		cfg:    cfg,
// 	}
// 	httpClient := &http.Client{
// 		Transport: signingTransport,
// 	}

// 	var opensearchconfig opensearchtransport.Config

// 	opensearchCfg := opensearch.Config{
// 		Addresses: []string{endpoint},
// 		Transport: opensearchtransport.New(opensearchconfig).WithHTTPClient(httpClient),
// 	}
// 	client, err := opensearch.NewClient(opensearchCfg)
// 	if err != nil {
// 		log.Fatalf("Error creating OpenSearch client: %v", err)
// 	}

// 	res, err := client.Info()
// 	if err != nil {
// 		log.Fatalf("Error retrieving OpenSearch info: %v", err)
// 	}
// 	defer res.Body.Close()
// }
