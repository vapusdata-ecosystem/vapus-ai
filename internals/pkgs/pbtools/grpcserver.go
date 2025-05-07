package pbtools

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	prometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
	grpc "google.golang.org/grpc"
	credentials "google.golang.org/grpc/credentials"
	grpc_insecure "google.golang.org/grpc/credentials/insecure"
)

// // grpcHandlerFunc routes HTTP and gRPC traffic to the respective servers.
// func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Check if the request is for gRPC based on the content type and protocol.
// 		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
// 			grpcServer.ServeHTTP(w, r)
// 		} else {
// 			otherHandler.ServeHTTP(w, r)
// 		}
// 	})
// }

// GRPCServer struct with configuration options to start and run the grpc server
type GRPCServer struct {
	ServerPort    int32
	ServerAddress string
	GrpcServ      *grpc.Server
	Netlis        net.Listener
	Logger        zerolog.Logger
	PlainTls      bool
	MTls          bool
	tlsCertFile   string
	tlsKeyFile    string
	tlsCaCertFile string
	Insecure      bool
	SetUpHttp     bool
	MuxInstance   *runtime.ServeMux
	HttpGwPort    int32
	Http1Address  string
}

type GrpcOptions func(*GRPCServer)

// NewGRPCServer creates a new grpc server with configuration options
func NewGRPCServer(opts ...GrpcOptions) *GRPCServer {
	grpcServer := &GRPCServer{}
	for _, o := range opts {
		o(grpcServer)
	}
	return grpcServer
}

func WithMtls(caFile, serverFile, keyFile string) GrpcOptions {
	return func(opt *GRPCServer) {
		opt.tlsCaCertFile = caFile
		opt.tlsCertFile = serverFile
		opt.tlsKeyFile = keyFile
		opt.MTls = true
	}
}

func WithPlainTls(serverFile, keyFile string) GrpcOptions {
	return func(opt *GRPCServer) {
		opt.tlsCertFile = serverFile
		opt.tlsKeyFile = keyFile
		opt.PlainTls = true
	}
}

func WithInsecure(insecure bool) GrpcOptions {
	return func(opt *GRPCServer) {
		opt.Insecure = insecure
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getServerAddress returns the server address
func (gs *GRPCServer) setServerAddress() {
	// if server address is not provided, use localhost
	if gs.ServerAddress == types.EMPTYSTR {
		gs.ServerAddress = fmt.Sprintf(":%d", gs.ServerPort)
	}
}

// getServerAddress returns the server address
func (gs *GRPCServer) setHttpServerAddress() {
	// if server address is not provided, use localhost
	gs.Http1Address = fmt.Sprintf(":%d", gs.HttpGwPort)
}

// NewGrpcServer creates a new grpc server with confirguration and other options
// This function will return the grpc server and net listner
func (gs *GRPCServer) ConfigureGrpcServer(opts []grpc.ServerOption, debugFlag bool) {
	var err error
	// Set the server address
	gs.setServerAddress()

	opts = append(opts, grpc.MaxRecvMsgSize(1024*1024*100),
		grpc.MaxSendMsgSize(1024*1024*100))

	// If tls is enabled, generate the credentials
	opts = gs.initServerCreds(opts)

	grpcServer := grpc.NewServer(opts...)

	//register the grpc server with prometheus
	grpc_prometheus.EnableHandlingTimeHistogram(grpc_prometheus.WithHistogramBuckets(prometheus.ExponentialBucketsRange(0.5, 70, 9)))
	grpc_prometheus.Register(grpcServer)
	gs.Netlis, err = net.Listen("tcp", gs.ServerAddress)
	if err != nil {
		gs.Logger.Fatal().Err(err).Msgf("Failed to listen on %s", gs.ServerAddress)
	}
	gs.GrpcServ = grpcServer
	if gs.SetUpHttp {
		gs.setHttpServerAddress()
		if gs.MuxInstance == nil {
			gs.MuxInstance = runtime.NewServeMux()
		}
	}
}

func (gs *GRPCServer) Run() {
	var err error
	httpServer := &http.Server{
		Addr:    gs.Http1Address,
		Handler: corsMiddleware(gs.MuxInstance),
	}
	if gs.SetUpHttp {
		go func() {
			gs.Logger.Info().Msgf("Starting HTTP server on %s", gs.Http1Address)
			if gs.Insecure {
				err = httpServer.ListenAndServe()
				if err != nil {
					gs.Logger.Fatal().Err(err).Msg("Failed to start HTTP server")
				}
			} else {
				err = httpServer.ListenAndServeTLS(gs.tlsCertFile, gs.tlsKeyFile)
				if err != nil {
					gs.Logger.Fatal().Err(err).Msg("Failed to start HTTP server")
				}
			}
		}()
	}
	// Handle graceful shutdown in a separate goroutine
	go func() {
		// Handle graceful shutdown
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		gs.Logger.Info().Msg("Shutting down server")

		// Stop the gRPC server gracefully
		gs.GrpcServ.GracefulStop()
		if gs.SetUpHttp {
			// Shut down the HTTP server
			if err := httpServer.Shutdown(context.Background()); err != nil {
				gs.Logger.Err(err).Msg("Failed to shut down HTTP server")
			} else {
				gs.Logger.Info().Msg("HTTP server stopped")
			}

		}

		log.Println("Server stopped")
	}()

	// Keep the main goroutine alive
	err = gs.GrpcServ.Serve(gs.Netlis)
	if err != nil {
		gs.Logger.Panic().Err(err)
	}
	select {}
}

func (gs *GRPCServer) initServerCreds(opts []grpc.ServerOption) []grpc.ServerOption {

	// If Plain tls is enabled, generate the credentials
	if gs.PlainTls {
		return gs.initPlainTls(opts)
	} else if gs.MTls {
		return gs.initMtls(opts)
	} else if gs.Insecure {
		return gs.initInsecure(opts)
	}
	return opts
}

// initPlainTls initializes the Plain TLS credentials for the GRPCServer.
// It generates the credentials using the provided TLS certificate and key files.
// If the certificate or key file is not provided, it falls back to the default files.
// It returns the updated list of server options with the Plain TLS credentials.
func (gs *GRPCServer) initPlainTls(opts []grpc.ServerOption) []grpc.ServerOption {
	gs.Logger.Info().Msg("Generating Plain TLS credentials")
	if gs.tlsCertFile == "" {
		gs.Logger.Info().Msg("Using default server cert file")
		gs.tlsCertFile = types.DEFAULTSERVERCERTTLS
	}
	if gs.tlsKeyFile == "" {
		gs.Logger.Info().Msg("Using default server key file")
		gs.tlsKeyFile = types.DEFAULTSERVERKEYTLS
	}
	creds, err := credentials.NewServerTLSFromFile(gs.tlsCertFile, gs.tlsKeyFile)
	if err != nil {
		gs.Logger.Fatal().Err(err).Msg("Failed to generate credentials")
	}
	return append(opts, grpc.Creds(creds))
}

// initMtls initializes the MTLS (Mutual TLS) credentials for the GRPCServer.
// It generates the MTLS credentials using the provided CA certificate, server certificate, and key.
// The function returns a list of GRPC server options with the MTLS credentials appended.
func (gs *GRPCServer) initMtls(opts []grpc.ServerOption) []grpc.ServerOption {
	gs.Logger.Info().Msg("Generating MTLS credentials")
	caPem, err := os.ReadFile(gs.tlsCaCertFile)
	if err != nil {
		gs.Logger.Fatal().Err(err)
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPem) {
		gs.Logger.Fatal().Err(err)
	}

	// read server cert & key
	serverCert, err := tls.LoadX509KeyPair(gs.tlsCertFile, gs.tlsKeyFile)
	if err != nil {
		gs.Logger.Fatal().Err(err)
	}

	// configuration of the certificate what we want to
	conf := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	//create tls certificate
	tlsCredentials := credentials.NewTLS(conf)
	return append(opts, grpc.Creds(tlsCredentials))
}

func (gs *GRPCServer) initInsecure(opts []grpc.ServerOption) []grpc.ServerOption {
	return append(opts, grpc.Creds(grpc_insecure.NewCredentials()))
}

func (gs *GRPCServer) Close() {
	gs.GrpcServ.GracefulStop()
}
