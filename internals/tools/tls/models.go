package tls

import (
	"crypto/x509"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
)

type TLSOperatorOpts struct {
	Logger zerolog.Logger
	Algo   mpb.EncryptionAlgo
}

type TLSCreateParams struct {
	BitSize  int
	Template x509.Certificate
}

type TLSCreateResponse struct {
	IsBase64Encoded bool
	CertPem         string
	KeyPem          string
}
