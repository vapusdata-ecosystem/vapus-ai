package utils

import "time"

const (
	//default port for datamarketplace server
	DEFAULTPORT          = 8900
	DEFAULTSERVERKEYTLS  = "x509/server_key.pem"
	DEFAULTSERVERCERTTLS = "x509/server_cert.pem"
	DEFAULT_CONFIG_TYPE  = "toml"
	DOT                  = "."
	EMPTYSTR             = ""
	HASHICORPVAULT       = "hashicorpVault"
	AWS_SECRETMANAGER    = "awsSecretManager"
)

// context Contexts

type ctxKeys string

const (
	// ContextKey for the context
	PARSED_CLAIMS    ctxKeys = "parsedClaims"
	ACCESS_TOKEN     ctxKeys = "accessToken"
	USERID_CTX       ctxKeys = "userId"
	CATALOG_IDEN             = "organizationCatalog"
	MARKETPLACE_IDEN         = "marketplaceId"
)

const (
	// Vault Secret engines for different resources
	DATAMARKETPLACE_SE_VAULT      = "vapus-datamarketplace"
	DataMarketplaceNODES_SE_VAULT = "vapus-DataMarketplacenodes"
	DATANODES_SE_VAULT            = "vapus-datanodes"
	VPR_SE_VAULT                  = "vapus-vpr"
	VDR_SE_VAULT                  = "vapus-vdr"
)

// User constants

const (
	DEFAULT_USER_INVITE_EXPIRE_TIME = time.Hour * 24 * 30
	DEFAULT_PLATFORM_AT_VALIDITY    = time.Hour * 24
	DEFAULT_PLATFORM_RT_VALIDITY    = time.Hour * 24 * 30
)

// REDIS action and keys that are constant
const (
	DATAMARKETPLACE_NODE_MAP = "datamarketplace_node_map"
	DOMAINNODES_TABLE        = "organizationnodes-table"
	DATANODES_TABLE          = "datanodes-table"
	DATAMARKETPLACE_TABLE    = "datamarketplace-table"
	DATANODES_METADATA_TABLE = "datamarketplace-metadata"
)

// ES Indexes
const (
	DATANODE_METADATA_ES_INDEX = "vapusdata-datasource-metadata"
	DOMAINNODES_INDEX          = "vapusdata-organization"
	DATANODES_INDEX            = "vapusdata-datasources"
	DATAMARKETPLACE_INDEX      = "vapusdata-datamarketplace"
	ACCOUNT_INDEX              = "vapusdata-accounts"
)

// DB action and keys that are configurable
var (
	DataSourcePackages   = "schema::datasource::packages::%v"
	DataSourceDatabase   = "schema::datasource::database::%v"
	DATANODE_MEADATA_KEY = "datanode-scheme-%v"
)
