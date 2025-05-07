package qdrant

import (
	"github.com/qdrant/go-client/qdrant"
	"github.com/rs/zerolog"
)

type QdrantOpts struct {
	Host   string
	Port   int
	ApiKey string
}

type QdrantStore struct {
	Opts   *QdrantOpts
	Client *qdrant.Client
	logger zerolog.Logger
}

func NewConnectQdrant(opts *QdrantOpts, l zerolog.Logger) (*QdrantStore, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   opts.Host,
		Port:   opts.Port,
		APIKey: opts.ApiKey,
		UseTLS: true, // uses default config with minimum TLS version set to 1.3
		// TLSConfig: &tls.Config{...},
		// GrpcOptions: []grpc.DialOption{},
	})
	if err != nil {
		l.Err(err).Msgf("Failed to create Client")
		return nil, err
	}

	return &QdrantStore{
		Opts:   opts,
		Client: client,
		logger: l,
	}, nil
}

func (m *QdrantStore) Close() {
	m.Client.Close()
	// m.DB = nil
}
