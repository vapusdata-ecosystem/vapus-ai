package routes

var (
	LoginRedirect    = "/login/redirect"
	Login            string
	LoginCallBack    string
	Logout           string
	OrganizationAuth = "/auth/domain/:domain"
	UIRoute          = "/ui"
	ManagePrefix     = "/manage"
	Register         = "/register"
)

const (
	ManageAIGroup                   = "/ai/manage"
	ManageAIModelNodes              = "/model-nodes"
	CreateAIModelNodes              = "/model-nodes/create"
	UpdateAIModelNodes              = "/model-nodes/:aiModelNodeId/update"
	ManageAIAgents                  = "/agents"
	CreateAIAgents                  = "/agents/create"
	UpdateAIAgents                  = "/agents/:agentId/update"
	ManageAIAgentsResource          = "/agents/:agentId"
	ManageAIPrompts                 = "/prompts"
	CreateAIPrompts                 = "/prompts/create"
	UpdateAIPrompts                 = "/prompts/:promptId/update"
	ManageAIPromptResource          = "/prompts/:promptId"
	ManageAIGuardrails              = "/guardrails"
	CreateAIGuardrails              = "/guardrails/create"
	UpdateAIGuardrails              = "/guardrails/:guardrailId/update"
	ManageAIGuardrailResource       = "/guardrails/:guardrailId"
	ManageAIManageModelNodeResource = "/model-nodes/:aiModelNodeId"
	// Nabhik Task
	NabhikTasks      = "/nabhiktasks"
	NabhikTaskDetail = "/nabhiktasks/:nabhikTaskId"
	// Nabhik Task Log
	// NabhikTaskLogList				= "/nabhik-"
)

const (
	StudioGroup = ""
	AIStudio    = "/ai-studio"
	AgentStudio = "/agents-studio"
	DataStudio  = "/fabric"
)

const (
	NabhikAI = "/nabhik"
)

const (
	NabhikTaskGroup = "/nabhik-task"
)
const (
	DataQueryServer = "/query-server"
)

const (
	InsightsGroup = "/insights"
	LLMInsights   = "/llms"
)

const (
	SettingsGroup         = "/settings"
	SettingsOrganizations = "/organisation"
	SettingToken          = "/tokens"
	SettingsPlatform      = "/platform"
	SettingsIntergation   = "/integrations"
	SettingsUsers         = "/users"
	SettingsUserResource  = "/users/:userId"

	SettingsPlugins               = "/plugins"
	SettingsPluginResource        = "/plugins/:pluginId"
	SettingsPlatformOrganizations = "/platform-organisations"
	SettingsPluginsCreate         = "/plugins/create"
	SettingsPluginsUpdate         = "/plugins/:pluginId/update"

	SecretStoreList    = "/secretstores"
	SecretStoreDetails = "/secretstores/:secretstoreName"
	CreateSecretStore  = "/secretstores/create"
	UpdateSecretStore  = "/secretstores/:updateName/update"
)

const (
	DevelopersGroup     = "/developers"
	DevelopersResources = "/resources"
	DevelopersEnums     = "/enums"
)

const (
	HomeGroup = "/"
)

const (
	MyOrganizationGroup       = "/data/manage"
	DataSources               = "/data-sources"
	DataSourcesResource       = "/data-sources/:dataSourceId"
	OrganizationObservability = "/observability"
	CreateDataSource          = "/data-sources/create"
)

const (
	ExploreGroup            = "/explore"
	ExpOrganizations        = "/domains"
	ExpDatasources          = "/datasources"
	ExpDatasourceResource   = "/datasources/:datasourceId"
	ExpOrganizationResource = "/domains/:domainId"
	ExpOrganizationUsers    = "/domains/:domainId/users"
)

const (
	SearchGroup = "/search"
)
