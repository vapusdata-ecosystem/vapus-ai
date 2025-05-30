package services

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/databricks/databricks-sql-go/logger"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	bedrock "github.com/vapusdata-ecosystem/vapusai/core/aistudio/guardrails/bedrock"
	mistral "github.com/vapusdata-ecosystem/vapusai/core/aistudio/guardrails/mistral"
	pangea "github.com/vapusdata-ecosystem/vapusai/core/aistudio/guardrails/pangea"
	aidmstore "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo/aistudio"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	processes "github.com/vapusdata-ecosystem/vapusai/core/process"
	types "github.com/vapusdata-ecosystem/vapusai/core/types"
)

type GuardrailPluginsIntAgent struct {
	*processes.VapusInterfaceBase
	managerRequest *pb.GuardrailsTypeGetterRequest
	getterRequest  *pb.GuardrailsGetterRequest
	result         *pb.GuardrailsTypeResponse
	dmStore        *aidmstore.AIStudioDMStore
}

type GuardrailPluginsIntAgentOpts func(*GuardrailPluginsIntAgent)

func WithGuardrailPluginsManagerRequest(managerRequest *pb.GuardrailsTypeGetterRequest) GuardrailPluginsIntAgentOpts {
	return func(v *GuardrailPluginsIntAgent) {
		v.managerRequest = managerRequest
	}
}

func WithGuardrailPluginsGetterRequest(getterRequest *pb.GuardrailsGetterRequest) GuardrailPluginsIntAgentOpts {
	return func(v *GuardrailPluginsIntAgent) {
		v.getterRequest = getterRequest
	}
}

func WithGuardrailPluginsManagerAction(action mpb.ResourceLcActions) GuardrailPluginsIntAgentOpts {
	return func(v *GuardrailPluginsIntAgent) {
		v.Action = action.String()
	}
}

func (s *AIStudioServices) NewGuardrailPluginsIntAgent(ctx context.Context, opts ...GuardrailPluginsIntAgentOpts) (*GuardrailPluginsIntAgent, error) {
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		s.Logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}
	guardrail := &GuardrailPluginsIntAgent{
		result:  &pb.GuardrailsTypeResponse{},
		dmStore: s.DMStore,
		VapusInterfaceBase: &processes.VapusInterfaceBase{
			CtxClaim: vapusPlatformClaim,
			InitAt:   dmutils.GetEpochTime(),
		},
	}
	for _, opt := range opts {
		opt(guardrail)
	}
	guardrail.SetAgentId()
	guardrail.Logger = pkgs.GetSubDMLogger(types.AIPROMPTAGENT.String(), guardrail.AgentId)
	return guardrail, nil
}

func (v *GuardrailPluginsIntAgent) GetAgentId() string {
	return v.AgentId
}

func (v *GuardrailPluginsIntAgent) GetResult() *pb.GuardrailsTypeResponse {
	v.FinishAt = dmutils.GetEpochTime()
	v.FinalLog()
	return v.result
}

func (v *GuardrailPluginsIntAgent) Act(ctx context.Context) error {
	switch v.GetAction() {
	case mpb.ResourceLcActions_LIST.String():
		return v.listGuardrailPlugins(ctx)
	default:
		v.Logger.Error().Msg("invalid action for AIGuardrailIntAgent")
		return dmerrors.DMError(apperr.ErrInvalidAction, nil)
	}
}

func (v *GuardrailPluginsIntAgent) listGuardrailPlugins(ctx context.Context) error {
	v.Logger.Info().Msg("fetching Bedrock Guardrail list")

	// list of plugins fetched
	listPlugins, err := v.dmStore.ListPlugins(ctx,
		fmt.Sprintf(`status = 'ACTIVE' AND deleted_at IS NULL AND 
		(
			(scope='PLATFORM_SCOPE') OR (scope='DOMAIN_SCOPE' AND organization='%s') OR (scope='USER_SCOPE' AND organization='%s' AND created_by='%s')
		)
		 ORDER BY created_at DESC`, v.CtxClaim[encryption.ClaimOrganizationKey], v.CtxClaim[encryption.ClaimOrganizationKey], v.CtxClaim[encryption.ClaimUserIdKey]),
		v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting plugin from datastore")
		return dmerrors.DMError(apperr.ErrPlugin404, err)
	}
	bedrockList := []*pb.BedrockGuardrailModels{}
	pangeaList := []*pb.ThirdPartyGuardrailModels{}
	pangeaCnt := 0
	mistralList := []*pb.ThirdPartyGuardrailModels{}
	mistralCnt := 0
	secrets := &models.GenericCredentialModel{} // so, that we can get the URL for pangea
	for _, val := range listPlugins {
		// fetching the secrets
		switch val.PluginService {
		case types.Bedrock_Guardrail.String():
			secrets, err = v.getPluginSecrets(ctx, val.NetworkParams.SecretName)
			if err != nil {
				v.Logger.Error().Err(err).Msg("error while fetching the secrets")
				return dmerrors.DMError(apperr.ErrVapusSecret404, err)
			}
			bedrockResp, err := v.bedrock(ctx, secrets)
			if err != nil {
				v.Logger.Error().Err(err).Msg("error while fetching the bedrock guardrail list")
				return dmerrors.DMError(apperr.ErrVapusSecret404, err) // update
			}
			bedrockList = append(bedrockList, bedrockResp...) // user might have two-three bedrock accounts

		case types.Pangea_Guardrail.String():
			secrets, err = v.getPluginSecrets(ctx, val.NetworkParams.SecretName)
			if err != nil {
				v.Logger.Error().Err(err).Msg("error while fetching the secrets")
				return dmerrors.DMError(apperr.ErrVapusSecret404, err)
			}
			clientExist, err := v.pangea(ctx, secrets, val.NetworkParams.URL)
			if err != nil {
				v.Logger.Error().Err(err).Msg("error while creating the pangea client")
				return dmerrors.DMError(apperr.ErrVapusSecret404, err) // update
			}
			if clientExist && pangeaCnt == 0 {
				for _, val := range types.PanegaGuardrailList {
					val.String()
					pangeaList = append(pangeaList, &pb.ThirdPartyGuardrailModels{
						Name: val.String(),
						Id:   val.String(),
					})
				}
				pangeaCnt++
			}

		case types.Mistral_Guardrail.String():
			secrets, err = v.getPluginSecrets(ctx, val.NetworkParams.SecretName)
			if err != nil {
				v.Logger.Error().Err(err).Msg("error while fetching the secrets")
				return dmerrors.DMError(apperr.ErrVapusSecret404, err)
			}
			clientExist, err := v.mistral(ctx, secrets)
			if err != nil {
				v.Logger.Error().Err(err).Msg("error while creating the mistral client")
				return dmerrors.DMError(apperr.ErrVapusSecret404, err) // update
			}
			if clientExist && mistralCnt == 0 {
				for _, val := range types.MistralGuardrailList {
					val.String()
					mistralList = append(mistralList, &pb.ThirdPartyGuardrailModels{
						Name: val.String(),
						Id:   val.String(),
					})
				}
				mistralCnt++
			}
		}
		// Now Sectrets have been checked....

	}
	// setting up the response
	v.result.BedrockOutput = append(v.result.BedrockOutput, bedrockList...)
	v.result.PangeaOutput = append(v.result.PangeaOutput, pangeaList...)
	v.result.MistralOutput = append(v.result.MistralOutput, mistralList...)
	return nil
}

func (v *GuardrailPluginsIntAgent) getPluginSecrets(ctx context.Context, secretName string) (*models.GenericCredentialModel, error) {
	secrets := &models.GenericCredentialModel{}
	secretStr, err := v.dmStore.VapusStore.SecretStore.ReadSecret(ctx, secretName)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while fetching the secrets")
		return nil, dmerrors.DMError(apperr.ErrVapusSecret404, err)
	}
	fmt.Println("Secrets: ", reflect.ValueOf(secretStr))
	fmt.Println("Secrets: ", secretStr)
	err = json.Unmarshal([]byte(dmutils.AnyToStr(secretStr)), secrets)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while unmarshelling creds from secret store.")
		return nil, dmerrors.DMError(apperr.ErrDataSourceCredsSecretGet, err)
	}
	fmt.Println("secrets after unmarsahling", secrets)
	fmt.Println(secrets.ApiTokenType)
	fmt.Println(secrets.AwsCreds.SecretAccessKey)
	return secrets, nil
}

func (v *GuardrailPluginsIntAgent) bedrock(ctx context.Context, secrets *models.GenericCredentialModel) ([]*pb.BedrockGuardrailModels, error) {
	bedrockList := []*pb.BedrockGuardrailModels{}
	// storing the secrets in creds
	bedrockClient, err := bedrock.NewBedrockGuardrails(ctx, &bedrock.BedrockOpts{
		SecretAccessKey: secrets.AwsCreds.SecretAccessKey,
		AccessKeyId:     secrets.AwsCreds.AccessKeyId,
		Region:          secrets.AwsCreds.Region,
	}, v.Logger)
	if err != nil {
		logger.Err(err).Msg("Error while creating AWS client")
		return nil, err
	}

	// Fetching the list
	identifiers := ""
	list, err := bedrockClient.ListGuardrails(ctx, &identifiers, v.Logger)
	if err != nil {
		logger.Err(err).Msg("Unable to list of guardrails")
		return nil, err
	}
	v.Logger.Debug().Msg("List of guardrails fetched successfully")

	// Adding the data to the response...
	for _, val := range list {
		bedrockList = append(bedrockList, &pb.BedrockGuardrailModels{
			Id:   val.Id,
			Arn:  val.ARN,
			Name: val.Name,
		})
	}
	return bedrockList, nil
}

func (v *GuardrailPluginsIntAgent) pangea(ctx context.Context, secrets *models.GenericCredentialModel, domainURL string) (bool, error) {
	_, err := pangea.NewPangeaGuardrail(ctx, &pangea.PanegaOpts{
		Token:  secrets.ApiToken,
		Domain: domainURL, // "aws.us.pangea.cloud"
	}, v.Logger)
	if err != nil {
		logger.Err(err).Msg("Error while creating Pangea client")
		return false, err
	}
	return true, nil
}

func (v *GuardrailPluginsIntAgent) mistral(ctx context.Context, secrets *models.GenericCredentialModel) (bool, error) {
	_, err := mistral.NewMistralGuardrail(ctx, secrets.ApiToken, v.Logger)
	if err != nil {
		logger.Err(err).Msg("Error while creating Mistral client")
		return false, err
	}
	fmt.Println("Mistral Client created sucessfully")
	return true, nil
}
