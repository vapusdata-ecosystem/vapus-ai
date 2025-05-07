package pkg

const (
	APPNAME              = "vapusctl"
	DATAMARKETPLACE      = "datamarketplace"
	DOMAIN               = "domain"
	DATASOURCE           = "dataSource"
	DEFAULTRPCCONTEXTSEC = 2

	// Actions
	GET      = "get"
	MANAGE   = "manage"
	DELETE   = "delete"
	LIST     = "list"
	DESCRIBE = "describe"

	NAME      = "name"
	NAMESPACE = "namespace"
	PORT      = "port"
	URL       = "url"
	ADDRESS   = "address"
)

// constants for the cli nase actions
const (
	ContextsCmdName           = "contexts"
	GetAction                 = "get"
	ResourcesCmdName          = "resources"
	AuthAction                = "auth"
	GetOps                    = "get"
	DescribeOps               = "describe"
	ActOps                    = "act"
	SearchOpts                = "search"
	GetPrompt                 = "prompt"
	SvcInfoResource           = "svcinfo"
	ConfigureOps              = "configure"
	SpecsOps                  = "spec"
	DeployOps                 = "deploy"
	ClearOps                  = "clear"
	ExplainOps                = "explain"
	ContextsOps               = "context"
	ConnectOps                = "connect"
	AuthOps                   = "auth"
	OperatorOps               = "operator"
	InstallerOps              = "install"
	SetupCmd                  = "setup"
	InstallerSecretSpecGenOps = "gen-secret-spec"
	InstallerUpgradeOps       = "upgrade"
	BuildOps                  = "build"
	UpgradeOps                = "upgrade"
	ValidateOps               = "validate"
	RunOps                    = "run"
	PublishOps                = "publish"
	UnPublishOps              = "unpublish"
	PruneOps                  = "prune"
	RestartOps                = "restart"
	StopOps                   = "stop"
	CreateOps                 = "create"
	UpdateOps                 = "update"
	AddUserOps                = "adduser"
	ArchiveOps                = "archive"
	SyncOps                   = "sync"
	ResumeOps                 = "resume"
	PauseOps                  = "pause"
	CloneOps                  = "clone"
)

// commands and resource var/constants
const (
	MarketplaceResource           = "marketplace"
	DataSourceResource            = "datasources"
	OrganizationResource          = "domains"
	GenTemplate                   = "gen-template"
	AccountResource               = "account"
	MetaDataResource              = "metadata"
	DataProductResource           = "dataproducts"
	DataWorkerResource            = "dataworkers"
	DataCatalogResource           = "dataCatalogs"
	AIModelNodeResource           = "aimodelnodes"
	DataWorkerDeploymentResource  = "workerdeployments"
	DataProductDeploymentResource = "vdcdeployments"
	UserResource                  = "users"
	LoginResource                 = "login"
	ConfigResource                = "config"
	SearchResource                = "search"
)

// request params keys
const (
	DataproductKey        = "dataproduct"
	DatacatalogKey        = "datacatalog"
	DatasourceKey         = "datasource"
	DataworkerKey         = "dataworker"
	DatamarketplaceKey    = "datamarketplace"
	DeploymentKey         = "deployment"
	SearchqueryKey        = "searchquery"
	OrganizationKey       = "domain"
	WithPassword          = "with-password"
	WorkerDeplomentKey    = "workerDeployment"
	DatasourceMetadataKey = "datasource-metadata"
)

const (
	ARTIFACTSTORE = "artifactstore"
	CACHESTORE    = "cachestore"
	DBSTORE       = "dbstore"
	SECRETSTORE   = "secretstore"
	AICONFIG      = "aiconfig"
	AUTHNSECRET   = "authnsecret"
	JWTSECRET     = "jwtsecret"
)

const (
	HemlInstallerPackage = "oci://asia-south1-docker.pkg.dev/vapusdata-beta/vapusdata-ecosystem/vapusdata-platform-helmchart"
)
