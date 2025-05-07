package pinecone

import (
	"github.com/pinecone-io/go-pinecone/v3/pinecone"
	"github.com/rs/zerolog"
)

type PineconeOpts struct {
	ApiKey string
}

type PineconeStore struct {
	Opts   *PineconeOpts
	Client *pinecone.Client
	logger zerolog.Logger
}

func NewConnectPinecone(opts *PineconeOpts, l zerolog.Logger) (*PineconeStore, error) {
	pc, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey: opts.ApiKey,
	})
	if err != nil {
		l.Err(err).Msgf("Failed to create Pinecone client")
		return nil, err
	}
	// ctx := context.Background()
	// indexName := "vapus-index"
	// metric := pinecone.Dotproduct
	// deletionProtection := pinecone.DeletionProtectionDisabled

	// idx, err := pc.CreatePodIndex(ctx, &pinecone.CreatePodIndexRequest{
	// 	Name:               indexName,
	// 	Metric:             &metric,
	// 	Dimension:          1024,
	// 	Environment:        "us-east1-gcp",
	// 	PodType:            "p1.x1",
	// 	DeletionProtection: &deletionProtection,
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to create pod-based index: %v /n", err)
	// } else {
	// 	fmt.Printf("Successfully created pod-based index: %v", idx.Name)
	// }
	return &PineconeStore{
		Opts:   opts,
		Client: pc,
		logger: l,
	}, nil
}

func (m *PineconeStore) Close() {
	// m.Client.Close()
	// m.DB = nil
}
