package tls

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"

	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

var ellipticCurveMap = map[int]elliptic.Curve{
	256: elliptic.P256(),
	384: elliptic.P384(),
	521: elliptic.P521(),
}

type ECDSATlsOperator struct {
	opts *TLSOperatorOpts
}

func NewECDSATlsOperator(opts *TLSOperatorOpts) *ECDSATlsOperator {
	return &ECDSATlsOperator{
		opts: opts,
	}
}

func (e *ECDSATlsOperator) buildTlsResponse(dir *os.Root) (*TLSCreateResponse, error) {
	cert, err := dir.Open("cert.pem")
	if err != nil {
		e.opts.Logger.Err(err).Msg("Failed to open cert.pem")
		return nil, err
	}
	defer cert.Close()
	certBytes, err := io.ReadAll(cert)
	if err != nil {
		e.opts.Logger.Err(err).Msg("Failed to read cert.pem")
		return nil, err
	}
	response := &TLSCreateResponse{}
	encodedBytes := base64.StdEncoding.EncodeToString(certBytes)
	response.CertPem = string(encodedBytes)
	response.IsBase64Encoded = true
	key, err := dir.Open("key.pem")
	if err != nil {
		e.opts.Logger.Err(err).Msg("Failed to open key.pem")
		return nil, err
	}
	defer key.Close()
	keyBytes, err := io.ReadAll(key)
	if err != nil {
		e.opts.Logger.Err(err).Msg("Failed to read key.pem")
		return nil, err
	}
	encodedBytes = base64.StdEncoding.EncodeToString(keyBytes)
	response.KeyPem = string(encodedBytes)
	return response, nil
}

func (e *ECDSATlsOperator) GenerateTlsPvtKey(params *TLSCreateParams) (*TLSCreateResponse, error) {
	fname := dmutils.GetUUID()
	os.Mkdir(fname, 0755)
	defer os.RemoveAll(fname)
	rootDir, err := os.OpenRoot(fname)
	if err != nil {
		e.opts.Logger.Err(err).Msg("Failed to open root directory")
		return nil, err
	}
	defer rootDir.Close()
	eCurve, ok := ellipticCurveMap[params.BitSize]
	if !ok {
		e.opts.Logger.Error().Msgf("Invalid bit size: %d", params.BitSize)
		return nil, fmt.Errorf("Invalid bit size: %d", params.BitSize)
	}

	privateKey, err := ecdsa.GenerateKey(eCurve, rand.Reader)
	if err != nil {
		e.opts.Logger.Err(err).Msg("Failed to generate ECDSA private key")
		return nil, err
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &params.Template, &params.Template, &privateKey.PublicKey, privateKey)
	if err != nil {
		e.opts.Logger.Err(err).Msg("Failed to generate ECDSA private key")
		return nil, err
	}

	certOut, err := rootDir.Create("cert.pem")
	if err != nil {
		e.opts.Logger.Err(err).Msg("Failed to create cert.pem")
		return nil, err
	}
	defer certOut.Close()
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		defer certOut.Close()
		e.opts.Logger.Err(err).Msg("Failed to encode certificate to PEM")
		return nil, err
	}

	keyOut, err := rootDir.OpenFile("key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600) // 0600 permissions: only owner can read/write
	if err != nil {
		e.opts.Logger.Err(err).Msg("Failed to generate ECDSA private key")
		return nil, err
	}
	defer keyOut.Close()
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		e.opts.Logger.Err(err).Msg("Failed marshalling ECDSA private key")
		return nil, err
	}

	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyBytes}); err != nil {
		e.opts.Logger.Err(err).Msg("Failed to write ECDSA private key to file")
		return nil, err
	}
	return e.buildTlsResponse(rootDir)
}
