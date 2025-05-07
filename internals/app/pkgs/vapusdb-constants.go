package apppkgs

import (
	"github.com/vapusdata-ecosystem/vapusai/core/models"
)

const (
	AccountsTable                        = "accounts"
	UsersTable                           = "users"
	DataSourcesTable                     = "data_sources"
	AIModelsNodesTable                   = "ai_models"
	OrganizationsTable                   = "organizations"
	JwtLogsTable                         = "jwt_logs"
	UserTeamsTable                       = "teams"
	DataproductQueryLogTable             = "dataproduct_query_logs"
	AIModelPromptTable                   = "ai_prompts"
	RefreshTokenLogsTable                = "refresh_token_logs"
	VapusAIAgentsTable                   = "ai_agents"
	UpDownVotesTable                     = "up_down_votes"
	StarReviewsTable                     = "star_reviews"
	VapusResourceArnTable                = "vapus_resource_arns"
	PluginsTable                         = "plugins"
	AIStudioLogsTable                    = "ai_studio_logs"
	AIAgentThreadsTable                  = "ai_agent_threads"
	VapusGuardrailsTable                 = "ai_guardrails"
	AIToolCallLogTable                   = "ai_tool_call_logs"
	AIStudioUsagesTable                  = "ai_studio_usages"
	NabrunnerLogTable                    = "nabrunner_logs"
	SecretStoreTable                     = "secret_stores"
	DatamarketplaceFederatedCatalogTable = "datamarketplace_federated_catalogs"
	AIGuardrailsLogsTable                = "ai_guardrails_logs"
	AIStudioChatsTable                   = "ai_studio_chats"
	VapusAgentsTable                     = "vapus_agents"
	VapusAgentLogTable                   = "vapus_agents_logs"
	AIModelPriceListTable                = "ai_model_price_list"
	FileStoreLogTable                    = "file_store_logs"
)

var DBTablesMap = map[string]any{
	AccountsTable:         &models.Account{},
	UsersTable:            &models.Users{},
	DataSourcesTable:      &models.DataSource{},
	AIModelsNodesTable:    &models.AIModelNode{},
	OrganizationsTable:    &models.Organization{},
	JwtLogsTable:          &models.JwtLog{},
	UserTeamsTable:        &models.Team{},
	AIModelPromptTable:    &models.AIPrompt{},
	RefreshTokenLogsTable: &models.RefreshTokenLog{},
	UpDownVotesTable:      &models.UpDownVote{},
	StarReviewsTable:      &models.StarReview{},
	VapusResourceArnTable: &models.VapusResourceArn{},
	PluginsTable:          &models.Plugin{},
	AIStudioLogsTable:     &models.AIStudioLog{},
	VapusGuardrailsTable:  &models.AIGuardrails{},
	AIToolCallLogTable:    &models.AIToolCallLog{},
	AIStudioUsagesTable:   &models.AIStudioUsages{},
	NabrunnerLogTable:     &models.NabrunnerLog{},
	SecretStoreTable:      &models.SecretStore{},
	AIGuardrailsLogsTable: &models.AIGuardrailsLog{},
	AIStudioChatsTable:    &models.AIStudioChat{},
	VapusAgentsTable:      &models.VapusAgents{},
	VapusAgentLogTable:    &models.VapusAgentLog{},
	AIModelPriceListTable: &models.AIModelPriceList{},
	FileStoreLogTable:     &models.FileStoreLog{},
}
