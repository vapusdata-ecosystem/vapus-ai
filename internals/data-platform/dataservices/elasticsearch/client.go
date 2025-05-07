package elasticsearch

// import (
// 	"crypto/tls"
// 	"net"
// 	"net/http"
// 	"time"

// 	elasticSearch "github.com/elastic/go-elasticsearch/v8"

// 	"github.com/rs/zerolog"
// )

// type ElasticSearchStore struct {
// 	Client  *elasticSearch.Client
// 	TClient *elasticSearch.TypedClient
// 	logger  zerolog.Logger
// }

// func NewElasticSearchStoreNEW(opts *ElasticSearch, logger zerolog.Logger) (*ElasticSearchStore, error) {
// 	cfg := elasticSearch.Config{
// 		Addresses: []string{
// 			opts.URL,
// 		},
// 	}

// 	if opts.ApiKey != "" {
// 		cfg.APIKey = opts.ApiKey
// 	} else if opts.Username != "" && opts.Password != "" {
// 		cfg.Username = opts.Username
// 		cfg.Password = opts.Password
// 	}

// 	cfg.Transport = &http.Transport{
// 		MaxIdleConnsPerHost: 10,
// 		DialContext:         (&net.Dialer{Timeout: time.Second}).DialContext,
// 		TLSClientConfig:     &tls.Config{InsecureSkipVerify: false},
// 	}
// 	tes, err := elasticSearch.NewTypedClient(cfg)
// 	if err != nil {
// 		logger.Err(err).Msg("Error creating ES client")
// 		return nil, err
// 	}
// 	es, err := elasticSearch.NewClient(cfg)
// 	if err != nil {
// 		logger.Err(err).Msg("Error creating ES client")
// 		return nil, err
// 	}
// 	return &ElasticSearchStore{
// 		Client:  es,
// 		TClient: tes,
// 		logger:  logger,
// 	}, nil
// }
