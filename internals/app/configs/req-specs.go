package appconfigs

import (
	"maps"
	"strings"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusai/core/pkgs/pbtools"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func MarshalSpecs(spec protoreflect.ProtoMessage) string {
	// req := &mpb.MainRequestSpec{
	// 	ApiVersion: "v1alpha1",
	// 	Kind:       spec.ProtoReflect().Descriptor().Name(),
	// }
	byess, err := pbtools.ProtoYamlMarshal(spec)
	if err != nil {
		return ""
	}
	specStr := string(byess)
	return strings.TrimSpace(specStr)
}

var SpecMap = map[mpb.Resources]string{
	mpb.Resources_ORGANIZATIONS: MarshalSpecs(OrganizationManagerRequest),
	mpb.Resources_AIMODELS:      MarshalSpecs(AinodeConfiguratorRequest),
	mpb.Resources_ACCOUNT:       MarshalSpecs(AccountManagerRequest),
	mpb.Resources_AIPROMPTS:     MarshalSpecs(AiPromptManagerRequest),
	mpb.Resources_PLUGINS:       MarshalSpecs(PluginManagerRequest),
	mpb.Resources_AIGUARDRAILS:  MarshalSpecs(AiGuardrailManagerRequest),
	mpb.Resources_AIAGENTS:      MarshalSpecs(AIAgentsManagerRequest),
}

var ResourceActionsMap = map[mpb.Resources][]string{
	mpb.Resources_ACCOUNT: {
		pb.AccountAgentActions_UPDATE_PROFILE.String(),
		pb.AccountAgentActions_CONFIGURE_AISTUDIO_MODEL.String(),
	},

	mpb.Resources_PLUGINS: {
		pb.PluginAgentAction_CONFIGURE_PLUGIN.String(),
		pb.PluginAgentAction_PATCH_PLUGIN.String(),
		pb.PluginAgentAction_TEST_PLUGIN.String(),
	},
	mpb.Resources_AIAGENTS: {
		pb.PluginAgentAction_CONFIGURE_PLUGIN.String(),
		pb.PluginAgentAction_PATCH_PLUGIN.String(),
		pb.PluginAgentAction_TEST_PLUGIN.String(),
	},
	mpb.Resources_AIGUARDRAILS: {
		pb.PluginAgentAction_CONFIGURE_PLUGIN.String(),
		pb.PluginAgentAction_PATCH_PLUGIN.String(),
		pb.PluginAgentAction_TEST_PLUGIN.String(),
	},
	mpb.Resources_AIPROMPTS: {
		pb.PluginAgentAction_CONFIGURE_PLUGIN.String(),
		pb.PluginAgentAction_PATCH_PLUGIN.String(),
		pb.PluginAgentAction_TEST_PLUGIN.String(),
	},
	mpb.Resources_AIMODELS: {
		pb.PluginAgentAction_CONFIGURE_PLUGIN.String(),
		pb.PluginAgentAction_PATCH_PLUGIN.String(),
		pb.PluginAgentAction_TEST_PLUGIN.String(),
	},

	mpb.Resources_ORGANIZATIONS: {
		pb.PluginAgentAction_CONFIGURE_PLUGIN.String(),
		pb.PluginAgentAction_PATCH_PLUGIN.String(),
		pb.PluginAgentAction_TEST_PLUGIN.String(),
	},
}

var AccountManagerRequest *pb.AccountManagerRequest = &pb.AccountManagerRequest{
	Actions: pb.AccountAgentActions_UPDATE_PROFILE,
	Spec: &mpb.Account{
		AiAttributes: &mpb.AIAttributes{},
	},
}

var AinodeConfiguratorRequest *pb.AIModelNodeManagerRequest = &pb.AIModelNodeManagerRequest{
	Spec: &mpb.AIModelNode{
		Attributes: &mpb.AIModelNodeAttributes{
			NetworkParams: &mpb.AIModelNodeNetworkParams{
				Credentials: DataSourceCreds,
			},
		},
	},
}

var AIAgentsManagerRequest *pb.AgentManagerRequest = &pb.AgentManagerRequest{
	Spec: &mpb.VapusAgent{
		Attributes: &mpb.AgentAttributes{
			Schedule: &mpb.VapusSchedule{
				CronTab: &mpb.CronTab{
					FrequencyTab: []*mpb.FrequencyTab{{}},
				},
			},
		},
		Specs: []*mpb.AgentSpec{
			{
				InputToolCalls: []*mpb.FunctionCall{
					{},
				},
			},
		},
	},
}

var AiGuardrailManagerRequest *pb.GuardrailsManagerRequest = &pb.GuardrailsManagerRequest{
	Spec: &mpb.AIGuardrails{
		Contents:         &mpb.ContentGuardrailLevel{},
		Topics:           []*mpb.TopicGuardrails{{}},
		Words:            []*mpb.WordGuardRails{{}},
		SensitiveDataset: []*mpb.SensitiveDataGuardrails{{}},
		ResourceBase:     &mpb.VapusBase{},
		GuardModel:       &mpb.GuardModels{},
	},
}

var AiPromptManagerRequest *pb.PromptManagerRequest = &pb.PromptManagerRequest{
	Spec: &mpb.AIPrompt{
		Spec: &mpb.PromptSpec{
			Sample: &mpb.Sample{},
			Tools: []*mpb.ToolPrompts{{
				Schema: &mpb.FunctionCall{},
			}},
			ResponseFormat: &mpb.StructuredResponseFormat{
				JsonSchema: &mpb.ResponseJsonSchema{},
			},
		},
	},
}

var PluginManagerRequest *pb.PluginManagerRequest = &pb.PluginManagerRequest{
	Action: pb.PluginAgentAction_CONFIGURE_PLUGIN,
	Spec: &mpb.Plugin{
		NetworkParams: &mpb.PluginNetworkParams{
			Credentials: DataSourceCreds,
		},
		DynamicParams: []*mpb.Mapper{{}},
	},
}

var DatasourceManagerRequest *pb.DataSourceManagerRequest = &pb.DataSourceManagerRequest{
	Spec: &mpb.DataSource{
		NetParams: &mpb.DataSourceNetParams{
			DsCreds: []*mpb.DataSourceCreds{
				{
					Credentials: DataSourceCreds,
				},
			},
		},
		Attributes:    &mpb.DataSourceAttributes{},
		Tags:          []*mpb.Mapper{{}},
		SharingParams: &mpb.DataSourceSharingParams{},
	},
}

var DataSourceCreds *mpb.GenericCredentialObj = &mpb.GenericCredentialObj{
	AwsCreds:   &mpb.AWSCreds{},
	GcpCreds:   &mpb.GCPCreds{},
	AzureCreds: &mpb.AzureCreds{},
}

var OrganizationManagerRequest *pb.OrganizationManagerRequest = &pb.OrganizationManagerRequest{
	Spec: &mpb.Organization{
		SecretPasscode: &mpb.CredentialSalt{},
		Attributes: &mpb.OrganizationAttributes{
			AuthnJwtParams: &mpb.JWTParams{},
		},
		BackendSecretStorage: &mpb.BackendStorages{},
		ArtifactStorage:      &mpb.BackendStorages{},
		DataProductInfraPlatform: []*mpb.K8SInfraParams{
			{
				Credentials: DataSourceCreds,
			},
		},
	},
}

var EnumSpecs = map[string]map[string]int32{
	"AuthnMethod":                   mpb.AuthnMethod_value,
	"EncryptionAlgo":                mpb.EncryptionAlgo_value,
	"BEStoreAccessScope":            mpb.BEStoreAccessScope_value,
	"TlsType":                       mpb.TlsType_value,
	"RequestObjects":                mpb.Resources_value,
	"TableNature":                   mpb.TableNature_value,
	"ApiTokenType":                  mpb.ApiTokenType_value,
	"Frequency":                     mpb.Frequency_value,
	"BackendStorageTypes":           mpb.BackendStorageTypes_value,
	"BackendStorageOnboarding":      mpb.BackendStorageOnboarding_value,
	"StorageService":                mpb.DataSourceServices_value,
	"SvcProvider":                   mpb.ServiceProvider_value,
	"DataSourceType":                mpb.DataSourceType_value,
	"ArtifactTypes":                 mpb.ArtifactTypes_value,
	"DataSourceAccessScope":         mpb.DataSourceAccessScope_value,
	"LLMQueryType":                  mpb.LLMQueryType_value,
	"InfraService":                  mpb.InfraService_value,
	"K8SServiceType":                mpb.K8SServiceType_value,
	"ComplianceFieldType":           mpb.ComplianceFieldType_value,
	"ComplianceTypes":               mpb.ComplianceTypes_value,
	"DataWorkerRunType":             mpb.DataWorkerRunType_value,
	"DataProductType":               mpb.DataProductType_value,
	"VDCOrchestratorScope":          mpb.VDCOrchestratorScope_value,
	"DataSensitivityClassification": mpb.DataSensitivityClassification_value,
	"ClassifiedTransformerActions":  mpb.ClassifiedTransformerActions_value,
	"ResourceScope":                 mpb.ResourceScope_value,
	"VersionBumpType":               mpb.VersionBumpType_value,
	"OrderBys":                      mpb.OrderBys_value,
	"ContentFormats":                mpb.ContentFormats_value,
	"AIModelNodeHosting":            mpb.AIModelNodeHosting_value,
	"EmbeddingType":                 mpb.EmbeddingType_value,
	"AIAgentTypes":                  mpb.VapusAiAgentTypes_value,
	"IntegrationPluginTypes":        mpb.IntegrationPluginTypes_value,
	"AgentStepEnum":                 mpb.AgentStepEnum_value,
	"EmailSettings":                 mpb.EmailSettings_value,
	"GuardRailLevels":               mpb.GuardRailLevels_value,
	"AgentStepValueType":            mpb.AgentStepValueType_value,
	"AIToolCallType":                mpb.AIToolCallType_value,
	"AIGuardrailScanMode":           mpb.AIGuardrailScanMode_value,
	"AIResponseFormat":              mpb.AIResponseFormat_value,
	"VapusSecretType":               mpb.VapusSecretType_value,
}

func GetValidEnums() map[string]map[string]int32 {
	var validEnums = make(map[string]map[string]int32)
	for m, v := range EnumSpecs {
		validEnums[m] = v
		maps.DeleteFunc(v, func(k string, v int32) bool {
			return strings.Contains(k, "INVALID")
		})
	}
	return validEnums
}

var PluginTypes = map[mpb.IntegrationPluginTypes]map[types.PluginServices]string{
	mpb.IntegrationPluginTypes_EMAIL: {
		types.GMAIL:    "https://static.vecteezy.com/system/resources/previews/013/948/544/large_2x/gmail-logo-on-transparent-white-background-free-vector.jpg",
		types.SENDGRID: "https://vendure.io/_next/image?url=https%3A%2F%2Fhub.vendure.io%2Fassets%2Fpreview%2F51%2Fsendgrid-logo__preview.png%3Fh%3D800%26w%3D800%26q%3D90%26mode%3Dcrop%26fpy%3Dundefined%26fpx%3Dundefined%26format%3Dwebp&w=640&q=75",
		types.AWS_SES:  "https://symbols.getvecta.com/stencil_17/1_amazon-ses.f09d18050e.svg",
	},
	mpb.IntegrationPluginTypes_MESSAGING_CHANNELS: {
		types.SLACK: "https://static.vecteezy.com/system/resources/previews/048/759/332/large_2x/slack-transparent-icon-free-png.png",
	},
	mpb.IntegrationPluginTypes_FILESTORES: {
		types.GOOGLE_DRIVE: "https://www.logo.wine/a/logo/Google_Drive/Google_Drive-Logo.wine.svg",
	},
	mpb.IntegrationPluginTypes_SMS: {
		types.TWILIO:    "https://cdn.brandfetch.io/idT7wVo_zL/theme/dark/symbol.svg?c=1bxid64Mup7aczewSAYMX&t=1668515584517",
		types.AZURE_SMS: "https://swimburger.net/media/ppnn3pcl/azure.png",
	},
	mpb.IntegrationPluginTypes_SEARCHAPI: {
		types.SERPSEARCH: "https://s3-us-west-2.amazonaws.com/anchor-generated-image-bank/production/podcast_uploaded_nologo400/30457559/30457559-1662746398635-28be940c9e9a.jpg",
	},
	mpb.IntegrationPluginTypes_GUARDRAILS: {
		types.Bedrock_Guardrail: "https://a0.awsstatic.com/libra-css/images/logos/aws_logo_smile_1200x630.png",
		types.Pangea_Guardrail:  "https://awsmp-logos.s3.amazonaws.com/seller-nqyldztntiony/7f339af3025e7696a9c33a632da937f4.png",
		types.Mistral_Guardrail: "https://avatars.githubusercontent.com/u/99472018",
	},
}
