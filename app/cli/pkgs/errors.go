package pkg

import "errors"

var (
	ErrNoArgs      = errors.New("no arguments provided for the command, please provide the required arguments")
	ErrInvalidArgs = errors.New("invalid arguments provided for the command, please provide the required arguments")

	ErrLoadingGrpcCert                   = errors.New("error loading grpc cert for grpc client, please check the cert path provided in cli config")
	ErrGrpcConn                          = errors.New("error establishing grpc connection with the datamarketplace server, please check the server address and authtoken provided in cli config")
	ErrNoDataMarketplaceFound            = errors.New("no datamarketplace found for the provided marketplaceid")
	ErrNoOrganizationNodeFound           = errors.New("no domain node found for the provided domain node id")
	ErrNoDataNodeFound                   = errors.New("no domain node found for the provided domain node id")
	ErrInvalidInterfaceGoal              = errors.New("invalid interface action provided, please provide a valid action")
	ErrNoCurrentContext                  = errors.New("no current context found, please set a current context")
	ErrInvalidAction                     = errors.New("invalid action provided, please provide a valid action")
	ErrMissingDataProductLogin           = errors.New("no data product provided for login")
	ErrInvalidDataWorkerType             = errors.New("invalid data worker type provided, please provide a valid data worker type")
	ErrMetaData404                       = errors.New("metadata not found for the provided data source")
	ErrVapusDataPlatformNotConnected     = errors.New("vapus data platform not connected, please connect to the platform using the connect command")
	ErrInvalidResource                   = errors.New("invalid resource provided, please provide a valid resource")
	ErrNoFileInput                       = errors.New("no file input provided, please provide a file input")
	ErrInvalidParamsForDescribe          = errors.New("invalid params provided for describe action, please provide a valid data source id")
	ErrInvalidRequestSpec                = errors.New("invalid request spec provided, please provide a valid request spec")
	ErrEmptyInput                        = errors.New("empty input provided, please provide a valid input")
	ErrSecretStoreParam404               = errors.New("secret store param file not found, please provide a valid file path")
	ErrFile404                           = errors.New("file not found, please provide a valid file path")
	ErrInvalidBoolInput                  = errors.New("invalid boolean input provided, please provide a valid boolean input (true/false)")
	ErrInvalidEncryptionAlgorithm        = errors.New("invalid encryption algorithm provided, please provide a valid encryption algorithm (ECDSA/RSA)")
	ErrInvalidEncryptionAlgorithmBitSize = errors.New("invalid encryption algorithm bit size provided, please provide a valid encryption algorithm bit size (256/384/521)")
)
