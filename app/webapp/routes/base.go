package routes

type NavMenuId string

type RouteBaseStruct struct {
	ItemName string
	Url      string
	ItemId   string
	Svg      string
	Children []RouteBaseStruct
}

var (
	HomeNav            NavMenuId = "vapushomeNavMenu"
	SettingsNav        NavMenuId = "settingsNavMenu"
	MyOrganizationNav  NavMenuId = "domainsNavMenu"
	DataMarketplaceNav NavMenuId = "dataMarketplaceNavMenu"
	ManageAINav        NavMenuId = "manageAINavMenu"
	ExploreNav         NavMenuId = "exploreNavMenu"
	VapusStudioNav     NavMenuId = "vapusStudioNavMenu"
	DevelopersNav      NavMenuId = "developersNavMenu"
	NabhikAINav        NavMenuId = "nabhikAINavMenu"
	InsightsNav        NavMenuId = "insightsNavMenu"
	DataQueryServerNav NavMenuId = "dataQueryServerNavMenu"
)

func (p NavMenuId) String() string {
	return string(p)
}

type SidebarId string

var (
	DataProductsPage                   SidebarId = "dataProducts"
	DataMarketplaceDataProductsPage    SidebarId = "dataMarketplaceDataProducts"
	DiscoverDataMarketplacePage        SidebarId = "discoverDataMarketplace"
	OrganizationSettingsPage           SidebarId = "domainSettings"
	DataMarketplaceSettingsPage        SidebarId = "dataMarketplaceSettingsPage"
	OrganizationDataSourcesPage        SidebarId = "domainDataSources"
	OrganizationsPage                  SidebarId = "domains"
	OrganizationDataWorkersPage        SidebarId = "domainDataWorkers"
	OrganizationVdcDeploymentsPage     SidebarId = "domainVdcDeployments"
	OrganizationUsersPage              SidebarId = "domainUsers"
	OrganizationDataProductspage       SidebarId = "domainDataProducts"
	OrganizationWorkersDeploymentsPage SidebarId = "domainWorkersDeployments"
	DataMarketplacePage                SidebarId = "dataMarketplace"
	DataMarketplaceOverviewPage        SidebarId = "dataMarketplaceOverview"
	OrganizationOverviewPage           SidebarId = "domainOverview"
	ManageAIOverviewPage               SidebarId = "aiStudioOverview"
	ManageAIModelNodesPage             SidebarId = "aiStudioModelNodes"
	AIPromptsPage                      SidebarId = "aiPrompts"
	AIModelDeploymentPage              SidebarId = "aiModelDeployments"
	AITrainerPage                      SidebarId = "aiTrainers"
	ManageAIAgentsPage                 SidebarId = "aiStudioAgents"
	ManageAIModelInterfacePage         SidebarId = "aiStudioModelInterface"
	DataCatalogPage                    SidebarId = "dataCatalog"
	MyDataProductsPage                 SidebarId = "myDataProducts"
	DataProductDiscoverPage            SidebarId = "dataProductDiscover"
	RequestNewDataProductPage          SidebarId = "requestNewDataProduct"
	SettingsPlatformOrganizationsPage  SidebarId = "platformOrganizations"
	SettingsProfilePage                SidebarId = "settingsProfile"
	SettingsPlatformPage               SidebarId = "settingsPlatform"
	SettingsIntegrationsPage           SidebarId = "settingsIntegrations"
	SettingsTokensPage                 SidebarId = "settingsTokens"
	SettingsUsersPage                  SidebarId = "settingsUsers"
	SettingsSecretStorePage            SidebarId = "settingsSecretStore"
	AccessRequestsPage                 SidebarId = "accessRequests"
	LogoutPage                         SidebarId = "logout"
	DevelopersResourcesPage            SidebarId = "settingsResources"
	DevelopersEnumsPage                SidebarId = "settingsEnums"
	CatalogdataProductsPage            SidebarId = "catalogDataProducts"
	ExpDatasourcesPage                 SidebarId = "exploreDatasources"
	SettingsPluginsPage                SidebarId = "settingsPlugins"
	AIStudioPage                       SidebarId = "aiStudio"
	AgentStudioPage                    SidebarId = "agentStudio"
	DataStudioPage                     SidebarId = "dataStudio"
	ManageAIGuardrailsPage             SidebarId = "aiGuardrails"
	OrganizationObservabilityPage      SidebarId = "domainObservability"
	LLMInsightsPage                    SidebarId = "llmInsights"
	NabhikTaskPage                     SidebarId = "nabhikTask"
)

func (p SidebarId) String() string {
	return string(p)
}

var (
	NavMenuRoutes                []RouteBaseStruct
	DataMarketplaceSidebarRoutes []RouteBaseStruct
	OrganizationSidebarRoutes    []RouteBaseStruct
	ManageAISidebarRoutes        []RouteBaseStruct
)

var NavMenuList = []RouteBaseStruct{
	{
		ItemName: "Dashboard",
		ItemId:   HomeNav.String(),
		Url:      UIRoute + HomeGroup,
		Svg: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="h-5 w-5 m-1" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M3 9l9-6 9 6v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
  <path d="M9 22V12h6v10"></path>
</svg>
`,
	},
	{
		ItemName: "Insights",
		ItemId:   InsightsNav.String(),
		Url:      UIRoute + InsightsGroup + LLMInsights,
		Svg: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="h-6 w-6 m-1" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-6 w-6 m-1">
  <title>Observability Icon 1</title>
  <!-- Line graph element -->
  <polyline points="3 17 8 11 13 13 17 8 21 12" />
  <!-- Bar elements (using fill) -->
  <rect x="5" y="18" width="4" height="3" stroke="none" fill="currentColor"/>
  <rect x="10" y="16" width="4" height="5" stroke="none" fill="currentColor"/>
  <rect x="15" y="17" width="4" height="4" stroke="none" fill="currentColor"/>
  <!-- Optional Axes (thinner stroke) -->
  <!-- <line x1="3" y1="21" x2="21" y2="21" stroke-width="1"/> -->
  <!-- <line x1="3" y1="21" x2="3" y2="6" stroke-width="1"/> -->
</svg>
`,
		Children: InsightsNavSideList,
	},
	{
		ItemName: "Studios",
		ItemId:   VapusStudioNav.String(),
		Url:      UIRoute + StudioGroup + AIStudio,
		Svg: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20" class="h-5 w-5 m-1"><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-miterlimit="22.926" stroke-width="1.5" d="m6.514 9.06-3.988-.383 3.217-3.216a3.36 3.36 0 0 1 3.925-.595M10.95 13.55l.377 3.924 3.217-3.216a3.354 3.354 0 0 0 .52-4.06"></path><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-miterlimit="22.926" stroke-width="1.5" d="M15.005 10.166q-.36.416-.76.817c-1.342 1.341-2.838 2.347-4.315 2.978a1.11 1.11 0 0 1-1.24-.24l-2.416-2.414a1.11 1.11 0 0 1-.24-1.24c.632-1.477 1.638-2.972 2.98-4.314 2.815-2.814 6.309-4.151 8.882-3.65.454 2.33-.599 5.412-2.89 8.063M13.673 2.65l3.556 3.555"></path><path fill="currentColor" d="M13.56 6.44a1.5 1.5 0 1 1-2.12 2.12 1.5 1.5 0 0 1 2.12-2.12"></path><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-miterlimit="22.926" stroke-width="1.5" d="M6.015 13.987 2 18M4.42 12.392l-1.357 1.356M7.61 15.581l-1.356 1.356"></path></svg>

`,
		// Children: StudioNavList,
	},
	{
		ItemName: "AI Center",
		ItemId:   ManageAINav.String(),
		Url:      UIRoute + ManageAIGroup + ManageAIModelNodes,
		Svg: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-6 w-6 m-1">
  <title>AI Studio Icon 5</title>
  <!-- Code Brackets -->
  <polyline points="8 3 3 9 8 15"/>
  <polyline points="16 3 21 9 16 15"/>
  <!-- AI Node/Core -->
  <circle cx="12" cy="18" r="3" fill="currentColor"/>
  <line x1="12" y1="15" x2="12" y2="9"/> <!-- Connection line -->
</svg>
`,
		Children: ManageAINavSideList,
	},
}

var BottonMenuList = []RouteBaseStruct{
	{
		ItemName: "Settings",
		ItemId:   SettingsNav.String(),
		Url:      UIRoute + SettingsGroup + SettingsOrganizations,
		Svg: `<svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-settings h-5 w-5 m-1" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
  <path stroke="none" d="M0 0h24v24H0z" fill="none"/>
  <path d="M10.325 4.317a1.724 1.724 0 0 1 3.35 0l.333 1.35a7.04 7.04 0 0 1 1.817 .621l1.308 -.478a1.724 1.724 0 0 1 2.156 2.156l-.478 1.308a7.034 7.034 0 0 1 .621 1.817l1.35 .333a1.724 1.724 0 0 1 0 3.35l-1.35 .333a7.034 7.034 0 0 1 -.621 1.817l.478 1.308a1.724 1.724 0 0 1 -2.156 2.156l-1.308 -.478a7.04 7.04 0 0 1 -1.817 .621l-.333 1.35a1.724 1.724 0 0 1 -3.35 0l-.333 -1.35a7.04 7.04 0 0 1 -1.817 -.621l-1.308 .478a1.724 1.724 0 0 1 -2.156 -2.156l.478 -1.308a7.034 7.034 0 0 1 -.621 -1.817l-1.35 -.333a1.724 1.724 0 0 1 0 -3.35l1.35 -.333a7.034 7.034 0 0 1 .621 -1.817l-.478 -1.308a1.724 1.724 0 0 1 2.156 -2.156l1.308 .478a7.04 7.04 0 0 1 1.817 -.621z" />
  <circle cx="12" cy="12" r="3" />
</svg>
`,
		Children: SettingsSideList,
	},
	{
		ItemName: "Developers",
		ItemId:   DevelopersNav.String(),
		Url:      UIRoute + DevelopersGroup + DevelopersResources,
		Svg: `<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 m-1" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <polyline points="16 18 22 12 16 6"></polyline>
  <polyline points="8 6 2 12 8 18"></polyline>
</svg>
`,
		Children: DevelopersSideList,
	},
}

var DatamanagerSideList = []RouteBaseStruct{
	{ItemName: "Data Sources", ItemId: OrganizationDataSourcesPage.String(), Url: UIRoute + MyOrganizationGroup + DataSources},
	{ItemName: "Observability", ItemId: OrganizationObservabilityPage.String(), Url: UIRoute + MyOrganizationGroup + OrganizationObservability},
}

var ManageAINavSideList = []RouteBaseStruct{
	{ItemName: "Models Registry", ItemId: ManageAIModelNodesPage.String(), Url: UIRoute + ManageAIGroup + ManageAIModelNodes},
	{ItemName: "Prompts", ItemId: AIPromptsPage.String(), Url: UIRoute + ManageAIGroup + ManageAIPrompts},
	{ItemName: "Agents", ItemId: ManageAIAgentsPage.String(), Url: UIRoute + ManageAIGroup + ManageAIAgents},
	{ItemName: "Guardrails", ItemId: ManageAIGuardrailsPage.String(), Url: UIRoute + ManageAIGroup + ManageAIGuardrails},
	{ItemName: "Nabhik Task", ItemId: NabhikTaskPage.String(), Url: UIRoute + ManageAIGroup + NabhikTasks},
}

var InsightsNavSideList = []RouteBaseStruct{
	{ItemName: "LLMs", ItemId: LLMInsightsPage.String(), Url: UIRoute + InsightsGroup + LLMInsights},
}

var StudioNavList = []RouteBaseStruct{
	{ItemName: "Data Fabric", ItemId: DataStudioPage.String(), Url: UIRoute + StudioGroup + DataStudio},
	{ItemName: "AI Studio", ItemId: AIStudioPage.String(), Url: UIRoute + StudioGroup + AIStudio},
	{ItemName: "Agent Studio", ItemId: AgentStudioPage.String(), Url: UIRoute + StudioGroup + AgentStudio},
}

var SettingsSideList = []RouteBaseStruct{
	// {ItemName: "Profile", ItemId: SettingsProfilePage.String(), Url: UIRoute + SettingsGroup},
	{ItemName: "Organization", ItemId: OrganizationSettingsPage.String(), Url: UIRoute + SettingsGroup + SettingsOrganizations},
	{ItemName: "Platform", ItemId: SettingsPlatformPage.String(), Url: UIRoute + SettingsGroup + SettingsPlatform},
	// {ItemName: "Integrations", ItemId: SettingsIntergationsPage.String(), Url: UIRoute + SettingsGroup + SettingsIntergation},
	{ItemName: "Users", ItemId: SettingsUsersPage.String(), Url: UIRoute + SettingsGroup + SettingsUsers},
	{ItemName: "Plugins", ItemId: SettingsPluginsPage.String(), Url: UIRoute + SettingsGroup + SettingsPlugins},
	{ItemName: "Platform Organizations", ItemId: SettingsPlatformOrganizationsPage.String(), Url: UIRoute + SettingsGroup + SettingsPlatformOrganizations},
	{ItemName: "Secret Store", ItemId: SettingsSecretStorePage.String(), Url: UIRoute + SettingsGroup + SecretStoreList},
}

var DevelopersSideList = []RouteBaseStruct{
	{ItemName: "Resources", ItemId: DevelopersResourcesPage.String(), Url: UIRoute + DevelopersGroup + DevelopersResources},
	{ItemName: "Enums", ItemId: DevelopersEnumsPage.String(), Url: UIRoute + DevelopersGroup + DevelopersEnums},
	// {ItemName: "Tokens", ItemId: SettingTokenPage.String(), Url: UIRoute + SettingsGroup + SettingToken},
}

var (
	UIHome = UIRoute + HomeGroup

	ManageOrganizationHome = UIRoute + MyOrganizationGroup

	NabhikHome = UIRoute + NabhikAI

	DataServerHome = UIRoute + UIRoute + DataQueryServer

	ManageAIHome = UIRoute + ManageAIGroup + ManageAIModelNodes

	AIStudioHome = UIRoute + StudioGroup + AIStudio

	AgentStudioHome = UIRoute + StudioGroup + AgentStudio

	SecretServiceHome = UIRoute + SettingsGroup + SecretStoreDetails
)
