package types

const (
	WINDOWS = "windows"
	MACOS   = "macos"
	LINUX   = "linux"

	EMPTYSTR  = ""
	LOCALHOST = "localhost"

	DEFAULTSERVERKEYTLS  = "x509/server_key.pem"
	DEFAULTSERVERCERTTLS = "x509/server_cert.pem"

	DEFAULT_CONFIG_TYPE = "json"
	JSON                = "json"
	TOML                = "toml"
	DOT                 = "."

	SERVICE_CONFIG_READ_ERROR = "failed to read the service configuration file"

	HASHICORPVAULT = "hashicorpVault"

	CUSTOM_MESSAGE = "customMessage"

	ERROR_CTX       = "error"
	WARNING_CTX     = "warning"
	INFO_CTX        = "info"
	QueryKey        = "query"
	DataProductsKey = "dataproducts"
)
