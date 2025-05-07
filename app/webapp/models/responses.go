package models

import (
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	vpb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusai/webapp/routes"
)

type ResourceManagerParams struct {
	API              string
	SupportedActions []string
	YamlSpec         string
	ActionMap        map[string]string
	ChatAPI          string
}

type GlobalContexts struct {
	NavMenuMap          []routes.RouteBaseStruct
	BottomMenuMap       []routes.RouteBaseStruct
	SidebarMap          []routes.RouteBaseStruct
	CurrentNav          string
	CurrentSideBar      string
	UserInfo            *mpb.User
	Account             *mpb.Account
	CurrentOrganization *mpb.Organization
	LoginUrl            string
	AccessTokenKey      string
	Manager             bool
	OrganizationMap     map[string]string
	CurrentUrl          string
	IsOrganizationOwner bool
	IsPlatformOwner     bool
}

type DataMarketplaceSvcResponse struct {
	CurrentDataMarketplace string
	// DataMarketplaceProducts []*mpb.DataProduct
	BackListingLink     string
	HideBackListingLink bool
	DetailLink          string
	DataServerUrl       string
	ActionParams        *ResourceManagerParams
	ActionRules         []*ActionRule
}
type HomePageResponse struct {
	Dashboard           *mpb.OrganizationDashboard
	BackListingLink     string
	HideBackListingLink bool
}
type NabhikServerResponse struct {
	BackListingLink         string
	HideBackListingLink     bool
	SelectedDataProductId   string
	QueryServerUrl          string
	QueryStreamServerUrl    string
	ContainerQueryStreamUri string
	ContainerQueryUri       string
	ActionUrisJson          string
	AIAgents                []*mpb.VapusAgent
	InvalidFabrciChat       *InvalidFabrciChat
	DataNabhikUiUrl         string
}

type InvalidFabrciChat struct {
	Error    string
	StartNew bool
}

type ExploreResponse struct {
	Organizations       []*mpb.Organization
	OrganizationUser    []*mpb.User
	YourOrganizations   []string
	ActionParams        *ResourceManagerParams
	ActionRules         []*ActionRule
	Users               []*mpb.User
	BackListingLink     string
	HideBackListingLink bool
}

type DataContainerResponse struct {
	ConfigUri string
	Actions   []string
}

type SettingsResponse struct {
	ActionParams        *ResourceManagerParams
	ActionRules         []*ActionRule
	CurrentOrganization *mpb.Organization
	Users               []*mpb.User
	User                *mpb.User
	BackListingLink     string
	HideBackListingLink bool
	SpecMap             map[mpb.Resources]string
	ResourceActionsMap  map[mpb.Resources][]string
	Enums               map[string]map[string]int32
	Plugins             []*mpb.Plugin
	Plugin              *mpb.Plugin
	PluginTypeMap       []*vpb.PluginTypeMap
	ResourceId          string
	YamlSpec            string
	CreateActionParams  *ActionRule
}

type SearchResponse struct {
	ActionParams        *ResourceManagerParams
	ActionRules         []*ActionRule
	CurrentOrganization *mpb.Organization
}

type OrganizationSvcResponse struct {
	CurrentOrganization *mpb.Organization
	DataSources         []*mpb.DataSource
	Users               []*mpb.User
	DataSource          *mpb.DataSource
	User                *mpb.User
	BackListingLink     string
	HideBackListingLink bool
	ActionRules         []*ActionRule
	ResourceId          string
	YamlSpec            string
	DataContainerOps    *DataContainerResponse
	SupportedActions    []string
	ActionParams        *ResourceManagerParams
	CreateActionParams  *ActionRule
}

type AIStudioResponse struct {
	AIModelNodes        []*mpb.AIModelNode
	AIModelNode         *mpb.AIModelNode
	ActionRules         []*ActionRule
	ResourceId          string
	BackListingLink     string
	HideBackListingLink bool
	AIPrompts           []*mpb.AIPrompt
	AIPrompt            *mpb.AIPrompt
	AIAgents            []*mpb.VapusAgent
	AIAgent             *mpb.VapusAgent
	AIGuardrail         *mpb.AIGuardrails
	AIGuardrails        []*mpb.AIGuardrails
	YamlSpec            string
	ActionParams        *ResourceManagerParams
	AIStudioChats       []*pb.AIStudioChat
	AIStudioChat        *pb.AIStudioChat
	CreateActionParams  *ActionRule
	DataSources         []*mpb.DataSource
	AIModelNodeInsights []*mpb.ModelNodeObservability
}

type AgentStudioResponse struct {
	AIModelNodes        []*mpb.AIModelNode
	AIModelNode         *mpb.AIModelNode
	ActionParams        *ResourceManagerParams
	ActionRules         []*ActionRule
	BackListingLink     string
	HideBackListingLink bool
	AIPrompts           []*mpb.AIPrompt
	AIPrompt            *mpb.AIPrompt
	AIAgents            []*mpb.VapusAgent
	AIAgent             *mpb.VapusAgent
	YamlSpec            string
}

type ActionRule struct {
	API              string
	Method           string
	YamlSpec         string
	ActionMap        map[string]string
	StreamAPI        string
	Action           string
	ResourceId       string
	Title            string
	IsRedirect       bool `default:"false"`
	SupportedActions []string
	Weblink          string
	JSONSpec         string
	AiStudioURL      string
}

type SecretServiceResponse struct {
	SecretStores        []*mpb.SecretStore
	SecretStore         *mpb.SecretStore
	BackListingLink     string
	HideBackListingLink bool
	ActionParams        *ResourceManagerParams
	ActionRules         []*ActionRule
	CreateActionParams  *ActionRule
	ResourceId          string
	YamlSpec            string
}

type NabhikTaskResponse struct {
	BackListingLink     string
	HideBackListingLink bool
	ActionParams        *ResourceManagerParams
	ActionRules         []*ActionRule
	CreateActionParams  *ActionRule
	ResourceId          string
	YamlSpec            string
}
