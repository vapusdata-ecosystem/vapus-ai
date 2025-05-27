package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/databricks/databricks-sql-go/logger"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"

	bedrock "github.com/vapusdata-ecosystem/vapusai/core/aistudio/guardrails/bedrock"
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
	managerRequest *pb.GuardrailsManagerRequest
	getterRequest  *pb.GuardrailsGetterRequest
	result         *pb.GuardrailsTypeResponse
	dmStore        *aidmstore.AIStudioDMStore
}

type GuardrailPluginsIntAgentOpts func(*GuardrailPluginsIntAgent)

func WithGuardrailPluginsManagerRequest(managerRequest *pb.GuardrailsManagerRequest) GuardrailPluginsIntAgentOpts {
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

	bedrockList := []*pb.BedrockGuardrailList{}
	for _, val := range listPlugins {
		if val.PluginService == types.Bedrock_Guardrail.String() {
			// fetching the secrets
			secrets := &models.GenericCredentialModel{}
			secretStr, err := v.dmStore.VapusStore.SecretStore.ReadSecret(ctx, val.NetworkParams.SecretName)
			if err != nil {
				v.Logger.Error().Err(err).Msg("error while fetching the secrets")
				return dmerrors.DMError(apperr.ErrVapusSecret404, err)
			}
			err = json.Unmarshal([]byte(dmutils.AnyToStr(secretStr)), secrets)
			if err != nil {
				logger.Err(err).Ctx(ctx).Msg("error while unmarshelling creds from secret store.")
				return dmerrors.DMError(apperr.ErrDataSourceCredsSecretGet, err)
			}

			// storing the secrets in creds
			bedrockClient, err := bedrock.NewBedrockGuardrails(ctx, &bedrock.BedrockOpts{
				SecretAccessKey: secrets.AwsCreds.SecretAccessKey,
				AccessKeyId:     secrets.AwsCreds.AccessKeyId,
				Region:          secrets.AwsCreds.Region,
			}, v.Logger)
			if err != nil {
				logger.Err(err).Msg("Error while creating AWS client")
				return err
			}

			// Fetching the list
			identifiers := ""
			list, err := bedrockClient.ListGuardrails(ctx, &identifiers, v.Logger)
			if err != nil {
				logger.Err(err).Msg("Unable to list of guardrails")
				return err
			}
			v.Logger.Debug().Msg("List of guardrails fetched successfully")
			// Adding the data to the response...
			for _, val := range list {
				bedrockList = append(bedrockList, &pb.BedrockGuardrailList{
					Id:   val.Id,
					Arn:  val.ARN,
					Name: val.Name,
				})
			}
		}
	}
	// setting up the response
	v.result.Output = append(v.result.Output, bedrockList...)
	return nil
}
