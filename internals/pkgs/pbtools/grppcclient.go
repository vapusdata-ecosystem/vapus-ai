package pbtools

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/rs/zerolog"

	grpc "google.golang.org/grpc"
	credentials "google.golang.org/grpc/credentials"
	grpc_insecure "google.golang.org/grpc/credentials/insecure"
)

// GrpcClient struct with configuration options to start and run the grpc server
type GrpcClient struct {
	serverAddress   string
	Connection      *grpc.ClientConn
	Logger          zerolog.Logger
	plainTls        bool
	mTls            bool
	tlsCertFile     string
	tlsKeyFile      string
	tlsCaCertFile   string
	Insecure        bool
	TlsORGANIZATION string
}

type GrpcClientOptions func(*GrpcClient)

// NewGrpcClient creates a new grpc client with configuration options
func NewGrpcClient(log zerolog.Logger, opts ...GrpcClientOptions) *GrpcClient {
	grpcClient := &GrpcClient{}
	for _, o := range opts {
		o(grpcClient)
	}
	grpcClient.initClient()
	return grpcClient
}

func ClientWithServiceAddress(addr string) GrpcClientOptions {
	return func(opt *GrpcClient) {
		opt.serverAddress = addr
	}
}

func ClientWithMtls(caFile, clientCert, keyFile string) GrpcClientOptions {
	return func(opt *GrpcClient) {
		opt.tlsCaCertFile = caFile
		opt.tlsCertFile = clientCert
		opt.tlsKeyFile = keyFile
		opt.mTls = true
	}
}

func ClientWithPlainTls(serverFile, keyFile string) GrpcClientOptions {
	return func(opt *GrpcClient) {
		opt.tlsCertFile = serverFile
		opt.tlsKeyFile = keyFile
		opt.plainTls = true
	}
}

func ClientWithInsecure(insecure bool) GrpcClientOptions {
	return func(opt *GrpcClient) {
		opt.Insecure = insecure
	}
}

// NewGrpcServer creates a new grpc server with confirguration and other options
// This function will return the grpc server and net listner
func (gc *GrpcClient) initClient() {
	var err error
	// Set the server address
	// If tls is enabled, generate the credentials
	creds := gc.initClientCreds()
	if creds == nil {
		gc.Logger.Err(err).Msg("Failed to generate client credentials")
	}

	conn, err := grpc.NewClient(gc.serverAddress, grpc.WithTransportCredentials(creds))
	if err != nil {
		gc.Logger.Err(err).Msg("Failed to connect to server")
	}
	gc.Connection = conn
}

func (gc *GrpcClient) initClientCreds() credentials.TransportCredentials {
	// If Plain tls is enabled, generate the credentials
	if gc.plainTls {
		return gc.initPlainTls()
	} else if gc.mTls {
		return gc.initMtls()
	} else if gc.Insecure {
		return gc.initInsecure()
	}
	return nil
}

func (gc *GrpcClient) initPlainTls() credentials.TransportCredentials {
	creds, err := credentials.NewClientTLSFromFile(gc.tlsCertFile, gc.TlsORGANIZATION)
	if err != nil {
		gc.Logger.Err(err).Msg("Failed to generate plaintls credentials for grpc client")
	}
	return creds
}

func (gc *GrpcClient) initMtls() credentials.TransportCredentials {
	gc.Logger.Info().Msg("Generating MTLS credentials")
	caPem, err := os.ReadFile(gc.tlsCaCertFile)
	if err != nil {
		gc.Logger.Err(err).Msg("Failed to read mtls ca certificate for grpc client")
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPem) {
		gc.Logger.Err(err).Msg("Failed to AppendCertsFromPEM mtls credentials for grpc client")
	}

	// read server cert & key
	serverCert, err := tls.LoadX509KeyPair(gc.tlsCertFile, gc.tlsKeyFile)
	if err != nil {
		gc.Logger.Err(err).Msg("Failed to LoadX509KeyPair mtls credentials for grpc client")
	}

	// configuration of the certificate what we want to
	conf := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		RootCAs:      certPool,
		ServerName:   gc.TlsORGANIZATION,
	}

	//create tls certificate
	tlsCredentials := credentials.NewTLS(conf)
	return tlsCredentials
}

func (gc *GrpcClient) initInsecure() credentials.TransportCredentials {
	return grpc_insecure.NewCredentials()
}

func (gc *GrpcClient) Close() {
	gc.Connection.Close()
}
