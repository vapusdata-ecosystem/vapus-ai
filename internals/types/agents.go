package types

type AgentStatus string

func (a AgentStatus) String() string {
	return string(a)
}

const (
	Active      AgentStatus = "ACTIVE"
	Initialized AgentStatus = "INITIALIZED"
	Waiting     AgentStatus = "WAITING"
	Failed      AgentStatus = "FAILED"
	Completed   AgentStatus = "COMPLETED"
	Retrying    AgentStatus = "RETRYING"
	Working     AgentStatus = "WORKING"
)

const (
	DataVisualizationRecommendation = "dataVisualizationRecommendation"
	DataAnalyticsRecommendation     = "dataAnalyticsRecommendation"
)

type NabhikMessageType string

func (a NabhikMessageType) String() string {
	return string(a)
}

const (
	TaskMessage NabhikMessageType = "TaskMessage"
	ChatMessage NabhikMessageType = "ChatMessage"
)

type AssetType string

func (x AssetType) String() string {
	return string(x)
}

const (
	ChartAssetType      AssetType = "ChartAssetType"
	DatasetAssetType    AssetType = "DatasetAssetType"
	AnalyzedDataset     AssetType = "AnalyzedDataset"
	SearchResultDataset AssetType = "SearchResultDataset"
	GenericDataset      AssetType = "GenericDataset"
	Reasoningset        AssetType = "Reasoningset"
	Knowledgeset        AssetType = "Knowledgeset"
	FilteredDataset     AssetType = "FilteredDataset"
	NextStepAssets      AssetType = "NextStepAssets"
	SummarySet          AssetType = "SummarySet"
)

type LeaderAgent string

func (a LeaderAgent) String() string {
	return string(a)
}

const (
	DataVisualizer      LeaderAgent = "DataVisualizer"
	Instructor          LeaderAgent = "Instructor"
	Communicator        LeaderAgent = "Communicator"
	FileOperator        LeaderAgent = "FileOperator"
	DataAnalyst         LeaderAgent = "DataAnalyst"
	DataQuerier         LeaderAgent = "DataQuerier"
	MetadataRetriever   LeaderAgent = "MetadataRetriever"
	DataFinder          LeaderAgent = "DataFinder"
	ContentWriter       LeaderAgent = "ContentWriter"
	DocuManager         LeaderAgent = "DocuManager"
	DataEngineer        LeaderAgent = "DataEngineer"
	DataAuditor         LeaderAgent = "DataAuditor"
	CodeManager         LeaderAgent = "CodeManager"
	MeetingScheduler    LeaderAgent = "MeetingScheduler"
	ProjectManager      LeaderAgent = "ProjectManager"
	TaskScheduler       LeaderAgent = "TaskScheduler"
	ChatInstructor      LeaderAgent = "ChatInstructor"
	TaskManager         LeaderAgent = "TaskManager"
	NabhikDone          LeaderAgent = "Done"
	HSNGSTRateValidator LeaderAgent = "HSNGSTRateValidator"
)

var AgentCapabilities = map[LeaderAgent]string{
	DataVisualizer:      "Creates and recommends data visualizations including charts like pie, bar, line, heatmap, scatter plots",
	Instructor:          "Provides instructions, guidance and suggestions to users",
	Communicator:        "Handles communication features such as sending emails or messages",
	FileOperator:        "Manages Files, file,datasets, assets, file storage, and cloud storage operations like uploading, downloading,listing, deleting and sharing",
	DataAnalyst:         "Performs data analysis including aggregation, filtering, grouping, and summarization",
	DataQuerier:         "Executes data query operations for selecting, fetching, or listing data from data sources",
	MetadataRetriever:   "Lists and searches metadata, especially when data source information is required",
	DataFinder:          "Searches and retrieves data products based on user requirements",
	ContentWriter:       "Creates written content like blogs, posts, articles, pdf, ppts, design doc, presentation, and documentation",
	DocuManager:         "Manages documentation and knowledge bases",
	DataEngineer:        "Handles ETL processes and data transformation pipelines",
	DataAuditor:         "Performs compliance checks and data quality audits",
	CodeManager:         "Manages code repositories and software development tasks",
	MeetingScheduler:    "Schedules and manages calendar events and meetings",
	ProjectManager:      "Handles project management tasks and team collaboration",
	TaskScheduler:       "Sets up and manages scheduled recurring tasks",
	TaskManager:         "Manage, Build, Run, Monitor, and Schedule tasks",
	ChatInstructor:      "Provides instructions, guidance, and user suggestions",
	HSNGSTRateValidator: "Validates HSN codes and GST rates, including classification and compliance checks",
}

type TeamMember string

func (a TeamMember) String() string {
	return string(a)
}

const (
	TaskLead                  TeamMember = "TaskLead"
	ToolDiscovery             TeamMember = "ToolDiscovery"
	DataVisualizationPlanner  TeamMember = "DataVisualizationPlanner"
	ChartReader               TeamMember = "ChartReader"
	ChartBuilder              TeamMember = "ChartBuilder"
	DashboardBuilder          TeamMember = "DashboardBuilder"
	DataAggregator            TeamMember = "DataAggregator"
	DataAnalysisPlanner       TeamMember = "DataAnalysisPlanner"
	DataSummarizer            TeamMember = "DataSummarizer"
	EmailOperator             TeamMember = "EmailOperator"
	GoogleDriverOperator      TeamMember = "GoogleDriverOperator"
	BoxUploader               TeamMember = "BoxUploader"
	DropBoxUploader           TeamMember = "DropBoxUploader"
	OneDriveUploader          TeamMember = "OneDriveUploader"
	BucketOperator            TeamMember = "BucketOperator"
	FilterDataset             TeamMember = "FilterDataset"
	SlackOperator             TeamMember = "SlackOperator"
	DatasetJoiner             TeamMember = "DatasetJoiner"
	BlogWriter                TeamMember = "BlogWriter"
	DocChatter                TeamMember = "DocChatter"
	ComplianceChecker         TeamMember = "ComplianceChecker"
	DocuWriter                TeamMember = "DocuWriter"
	ETLRunner                 TeamMember = "ETLRunner"
	TDMmanager                TeamMember = "TDMmanager"
	VectorGenerator           TeamMember = "VectorGenerator"
	GoogleCalendarOperator    TeamMember = "GoogleCalendarOperator"
	ZoomOperator              TeamMember = "ZoomOperator"
	MSTeamOperator            TeamMember = "MSTeamOperator"
	GithubOperator            TeamMember = "GithubOperator"
	GitlabOperator            TeamMember = "GitlabOperator"
	BitbucketOperator         TeamMember = "BitbucketOperator"
	JiraOperator              TeamMember = "JiraOperator"
	AsanaOperator             TeamMember = "AsanaOperator"
	ConfluenceOperator        TeamMember = "ConfluenceOperator"
	BlobHelper                TeamMember = "BlobHelper"
	DatabaseOperator          TeamMember = "DatabaseOperator"
	DataStreamOperator        TeamMember = "DataStreamOperator"
	APIOperator               TeamMember = "APIOperator"
	WebOperator               TeamMember = "WebOperator"
	FSOperator                TeamMember = "FSOperator"
	SuggestionManager         TeamMember = "SuggestionManager"
	DataProductNabhikOperator TeamMember = "DataProductNabhikOperator"
	DataProductRetriever      TeamMember = "DataProductRetriever"
	CronConvertor             TeamMember = "CronConvertor"
	DataAnalystEngineer       TeamMember = "DataAnalystEngineer"
	PythonProgrammer          TeamMember = "PythonProgrammer"
	HSNGSTRateOperator        TeamMember = "HSNGSTRateValidator"
	FieldMapDetector          TeamMember = "FieldMapDetector"
	SearchEngineOperator      TeamMember = "SearchEngineOperator"
	GSTCalculator             TeamMember = "GSTCalculator"
)

var ValidCrewMap = map[LeaderAgent][]TeamMember{
	DataVisualizer: {ChartBuilder, DashboardBuilder, DataVisualizationPlanner, ChartReader},
	Communicator:   {EmailOperator, SlackOperator, MSTeamOperator},
	DataAnalyst:    {DataAnalysisPlanner, FilterDataset, DatasetJoiner, DataSummarizer, DataAggregator, DataAnalystEngineer},
	FileOperator:   {GoogleDriverOperator, BoxUploader, DropBoxUploader, OneDriveUploader, BlobHelper},
	DataQuerier:    {DatabaseOperator, DataStreamOperator, APIOperator, WebOperator, FSOperator},
	Instructor:     {SuggestionManager},
}
