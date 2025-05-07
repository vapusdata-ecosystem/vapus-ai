package types

type SearchKey string

const (
	DataProuductSK    SearchKey = "dataProduct"
	DataWorkerSK      SearchKey = "dataWorker"
	DataSourceSK      SearchKey = "dataSource"
	ORGANIZATIONSK    SearchKey = "ORGANIZATION"
	DataCatalogSK     SearchKey = "dataCatalog"
	UserSK            SearchKey = "user"
	VdcDeploymentSK   SearchKey = "vdcDeploymentId"
	AIModelNodeSK     SearchKey = "aiModelNode"
	WorkerDeplomentSK SearchKey = "workerDeployment"
	DataStore         SearchKey = "datastore"
)

func (x SearchKey) String() string {
	return string(x)
}

type AgentType string

func (at AgentType) String() string {
	return string(at)
}

const (
	NABHIK_AGENT                 AgentType = "nabhik"
	VAPUSDATACONTAINERAGENT      AgentType = "vapusDataContainerAgent"
	DATAWORKERAGENT              AgentType = "dataWorkerAgent"
	VDCDEPLOYMENTAGENT           AgentType = "vdcDeploymentAgent"
	DATASOURCEAGENT              AgentType = "dataSourceAgent"
	DATAPRODUCTAGENT             AgentType = "dataProductAgent"
	DATAMARKETPLACEAGENT         AgentType = "dataMarketplaceAgent"
	AISTUDIONODE                 AgentType = "aiStudioNode"
	ACCOUNTAGENT                 AgentType = "accountAgent"
	DATASOURCEMETADATAAGENT      AgentType = "dataSourceMetadataAgent"
	VAPUSDATAPLATFORMSEARCHAGENT AgentType = "vapusDataPlatformSearchAgent"
	ORGANIZATIONAGENT            AgentType = "ORGANIZATIONAgent"
	USERMANAGERAGENT             AgentType = "userManagerAgent"
	AUTHZMANAGERAGENT            AgentType = "authzManagerAgent"
	AIPROMPTAGENT                AgentType = "aiPromptAgent"
	VAPUSAIAGENT                 AgentType = "vapusAIAgent"
	NABHIKAGENT                  AgentType = "nabhikAgent"
	SECRETMANAGERAGENT           AgentType = "secretManagerAgent"
)
