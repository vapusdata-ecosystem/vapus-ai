package tls

import (
	"fmt"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
)

type TLSOperatorOptions func(*TLSOperator)

type TLSOperatorService interface {
	// GenerateTLS certs and keys
	GenerateTlsPvtKey(params *TLSCreateParams) (*TLSCreateResponse, error)
}

type TLSOperator struct {
	opts *TLSOperatorOpts
}

func WithTLSOperatorParams(opts *TLSOperatorOpts) TLSOperatorOptions {
	return func(t *TLSOperator) {
		t.opts = opts
	}
}

func NewTLSOperator(opts ...TLSOperatorOptions) (TLSOperatorService, error) {
	t := &TLSOperator{}
	for _, opt := range opts {
		opt(t)
	}
	switch t.opts.Algo {
	case mpb.EncryptionAlgo_ECDSA:
		t.opts.Logger.Info().Msg("Using ECDSA algorithm")
		return NewECDSATlsOperator(t.opts), nil
	default:
		t.opts.Logger.Error().Msgf("Invalid algorithm: %s", t.opts.Algo)
		return nil, fmt.Errorf("invalid algorithm: %s", t.opts.Algo)
	}
}
