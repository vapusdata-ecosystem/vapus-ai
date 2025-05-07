package types

var (
	MARKETPLACEIDTEMPLATE = "marketplace-%s"

	ORGANIZATIONIDTEMPLATE = "dmn-%s"

	ACCOUNTIDTEMPLATE = "acc-%s"

	DATA_SOURCEID_TEMPLATE = "ds-%s"

	DATA_CATALOG_ID = "dc-%s"

	DATA_PRODUCT_ID = "dp-%v"

	DATA_PRODUCT_DEPLOYMENT_ID = "vdcp-%s"

	DATA_WORKER_DEPLOYMENT_ID = "dwd-%s"

	DATA_WORKER_ID = "dw-%s"

	VAPUS_AIMODEL_NODE_ID = "aimn-%s"

	DATA_PRODUCT_VERSION = "v%d"

	DATA_PRODUCT_ACCESS_REQUEST_ID = "dpar-%s"

	PROMPT_ID = "pr-%s"

	RESOURCE_ARN = "arn:vapusdata::%s:%s::"
)

var BASE_TABLE_TYPE = "BASE TABLE"
var VIEW_TYPE = "VIEW"

var (
	// AddressStr is a string template to generate full address including host and port.
	AddressStr = "%s:%d"

	IdSeparator = "_"
)

var (
	ClientRetryLimit = 3

	ClientRetryStart = 0
)

var (
	DataProductDescriptionVectorDimensions int = 1536
	DataSourceMetadataVectorDimensions     int = 1536
)

const (
	SENDEMAIL = "sendEmail"
)

// REDIS action and keys that are constant
const (
	LIST              = "list"
	ADD               = "Add"
	EXISTS            = "exists"
	COUNT             = "count"
	DEL               = "del"
	MADD              = "add-mulitple"
	ACCOUNT_KEY       = "account"
	ACCOUNT_DM_CF_KEY = "account:datamarketplace"
	DM_IDENTIFIER     = "datamarketplace"
	MARKETPLACEID     = "marketplaceId"
	DMNID             = "dmnId"
	DNID              = "dnId"
)

var FileContentTypes = map[string]string{
	"csv":  "text/csv",
	"json": "application/json",
	"xml":  "application/xml",
	"yaml": "application/yaml",
	"yml":  "application/yaml",
	"txt":  "text/plain",
	"tsv":  "text/tab-separated-values",
	"tab":  "text/tab-separated-values",
	"xls":  "application/vnd.ms-excel",
	"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"doc":  "application/msword",
	"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"pdf":  "application/pdf",
	"ppt":  "application/vnd.ms-powerpoint",
	"pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"zip":  "application/zip",
	"tar":  "application/x-tar",
	"gz":   "application/gzip",
	"bz2":  "application/x-bzip2",
	"rar":  "application/x-rar-compressed",
	"gif":  "image/gif",
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"svg":  "image/svg+xml",
	"bmp":  "image/bmp",
	"ico":  "image/vnd.microsoft.icon",
	"webp": "image/webp",
	"mp4":  "video/mp4",
	"html": "text/html",
	"htm":  "text/html",
	"css":  "text/css",
	"js":   "application/javascript",
	"woff": "font/woff",
}

const (
	MarketplaceId      = "marketplaceId"
	OrganizationId     = "OrganizationId"
	DataProductId      = "dataproductId"
	CatalogId          = "catalogId"
	PromptId           = "promptId"
	WorkerDeploymentId = "workerDeploymentId"
	VdcDeploymentId    = "vdcDeploymentId"
	UserId             = "userId"
	DataWorkerId       = "dataworkerId"
	DataSourceId       = "dataSourceId"
	AIModelNodeId      = "aiModelNodeId"
	AgentId            = "agentId"
	PluginId           = "pluginId"
	GuardrailId        = "guardrailId"
	NabhikChatId       = "chatId"
	AIStudioChatId     = "aiChatId"
	VapusAgentId       = "vapusAgentId"
	SecretStoreName    = "secretstoreName"
	NabhikTaskId       = "nabhikTaskId"
)

var UrlResouceIdMap = map[string]bool{
	MarketplaceId:      true,
	OrganizationId:     true,
	DataProductId:      true,
	CatalogId:          true,
	PromptId:           true,
	WorkerDeploymentId: true,
	VdcDeploymentId:    true,
	UserId:             true,
	DataWorkerId:       true,
	DataSourceId:       true,
	AIModelNodeId:      true,
	AgentId:            true,
	PluginId:           true,
	GuardrailId:        true,
	NabhikChatId:       true,
	VapusAgentId:       true,
	SecretStoreName:    true,
	NabhikTaskId:       true,
}

const (
	NabhikChatFileKey = "nabhikChat"
	VapusAgentFileKey = "nabhikAgent"
)

const (
	DSTORES = "dmtores"
)

const (
	ACCESS_TOKEN       = "access_token"
	REFRESH_TOKEN      = "refresh_token"
	USER_PROFILE       = "user_profile"
	AIMODEL_HEADER_KEY = "x-aimodelnode"
)
