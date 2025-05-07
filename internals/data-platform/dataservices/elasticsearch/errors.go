package elasticsearch

import "errors"

var (
	// Error constants for ElasticSearch operations
	ErrESConnection = errors.New("error while connecting to ElasticSearch")
)
