package elasticsearch

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"

	elasticSearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/rs/zerolog"
)

type ElasticSearchOpts struct {
	URL            string
	Port           int
	ApiKey         string
	Username       string
	Password       string
	CloudHost      string
	MaxConnections int32
}

type ElasticSearchStore struct {
	Opts    *ElasticSearchOpts
	ES      *elasticSearch.Client
	TClient *elasticSearch.TypedClient
	RES     *esapi.Response
	logger  zerolog.Logger
}

func NewElasticSearchStore(opts *ElasticSearchOpts, l zerolog.Logger) (*ElasticSearchStore, error) {

	// cfg := getDsn(opts)
	// l.Debug().Msgf("Connecting to Sql Server with dsn: %s", cfg)

	// es, err := elasticsearch.NewClient(cfg)
	// if err != nil {
	// 	log.Fatalf("Error creating the client: %s", err)
	// }

	cfg := elasticSearch.Config{
		// CloudID: opts.URL,
		Addresses: []string{
			opts.URL,
		},
	}

	if opts.ApiKey != "" {
		cfg.APIKey = opts.ApiKey
	} else if opts.Username != "" && opts.Password != "" {
		cfg.Username = opts.Username
		cfg.Password = opts.Password
	}

	cfg.Transport = &http.Transport{
		MaxIdleConnsPerHost: 10,
		DialContext:         (&net.Dialer{Timeout: time.Second}).DialContext,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: false},
	}

	tes, err := elasticSearch.NewTypedClient(cfg)
	if err != nil {
		l.Err(err).Msg("Error creating ES client")
		return nil, err
	}

	es, err := elasticSearch.NewClient(cfg)
	if err != nil {
		l.Err(err).Msg("Error creating ES client")
		return nil, err
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	// defer res.Body.Close()

	l.Info().Msg("Connected to Elastic Search database successfully")

	return &ElasticSearchStore{
		Opts:    opts,
		RES:     res,
		TClient: tes,
		logger:  l,
		ES:      es,
	}, nil
}

func (m *ElasticSearchStore) Close() {
	m.RES.Body.Close()

}

// func getDsn(opts *ElasticSearchOpts) elasticSearch.Config {
// 	var cfg elasticSearch.Config
// if opts.ApiKey != "" && opts.CloudHost != "" {
// 	cfg = elasticsearch.Config{
// 		CloudID: opts.CloudHost, // Get this from Elastic Cloud Console
// 		APIKey:  opts.ApiKey,    // Generate this in Elastic Cloud Console
// 	}
// } else if opts.Username != "" && opts.Password != "" {
// 	cfg = elasticSearch.Config{
// 		Addresses: []string{
// 			opts.URL,
// 		},
// 		Username: opts.Username,
// 		Password: opts.Password,
// 	}
// } else if opts.ApiKey != "" && opts.URL != "" {
// 	cfg = elasticSearch.Config{
// 		Addresses: []string{
// 			opts.URL,
// 		},
// 		APIKey: opts.ApiKey,
// 	}
// }

// 	return cfg
// }
