package services

import (
	"context"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/aistudio/pkgs"
	"github.com/vapusdata-ecosystem/vapusdata/aistudio/utils"
	aidmstore "github.com/vapusdata-ecosystem/vapusdata/core/app/datarepo/aistudio"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	apppkgss "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	models "github.com/vapusdata-ecosystem/vapusdata/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	processes "github.com/vapusdata-ecosystem/vapusdata/core/process"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
)

type AIPromptIntAgent struct {
	*processes.VapusInterfaceBase
	managerRequest *pb.PromptManagerRequest
	getterRequest  *pb.PromptGetterRequest
	result         []*mpb.AIPrompt
	dmStore        *aidmstore.AIStudioDMStore
}

type AIPromptIntAgentOpts func(*AIPromptIntAgent)

func WithPromptManagerRequest(managerRequest *pb.PromptManagerRequest) AIPromptIntAgentOpts {
	return func(v *AIPromptIntAgent) {
		v.managerRequest = managerRequest
	}
}

func WithPromptGetterRequest(getterRequest *pb.PromptGetterRequest) AIPromptIntAgentOpts {
	return func(v *AIPromptIntAgent) {
		v.getterRequest = getterRequest
	}
}

func WithPromptManagerAction(action mpb.ResourceLcActions) AIPromptIntAgentOpts {
	return func(v *AIPromptIntAgent) {
		v.Action = action.String()
	}
}

func (s *AIStudioServices) NewAIPromptIntAgent(ctx context.Context, opts ...AIPromptIntAgentOpts) (*AIPromptIntAgent, error) {
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		s.Logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}
	promptAgent := &AIPromptIntAgent{
		result:  make([]*mpb.AIPrompt, 0),
		dmStore: s.DMStore,
		VapusInterfaceBase: &processes.VapusInterfaceBase{
			CtxClaim: vapusPlatformClaim,
			// Ctx:      ctx,
			InitAt: dmutils.GetEpochTime(),
		},
	}
	promptAgent.SetAgentId()
	for _, opt := range opts {
		opt(promptAgent)
	}
	promptAgent.Logger = pkgs.GetSubDMLogger(types.AIPROMPTAGENT.String(), promptAgent.AgentId)
	return promptAgent, nil
}

func (v *AIPromptIntAgent) GetAgentId() string {
	return v.AgentId
}

func (v *AIPromptIntAgent) GetResult() []*mpb.AIPrompt {
	v.FinishAt = dmutils.GetEpochTime()
	v.FinalLog()
	return v.result
}

func (v *AIPromptIntAgent) Act(ctx context.Context) error {
	switch v.GetAction() {
	case mpb.ResourceLcActions_CREATE.String():
		return v.configureAIPrompts(ctx)
	case mpb.ResourceLcActions_UPDATE.String():
		return v.updateAIPrompts(ctx)
	case mpb.ResourceLcActions_GET.String():
		return v.describeAIPrompts(ctx)
	case mpb.ResourceLcActions_LIST.String():
		return v.listAIPrompts(ctx)
	case mpb.ResourceLcActions_ARCHIVE.String():
		return v.archiveAIPrompt(ctx)
	default:
		v.Logger.Error().Msg("invalid action for AIModelIntAgent")
		return dmerrors.DMError(apperr.ErrInvalidAction, nil)
	}
}

func (v *AIPromptIntAgent) configureAIPrompts(ctx context.Context) error {
	prompt := utils.AIPROPb2Obj(v.managerRequest.GetSpec())
	if prompt == nil {
		v.Logger.Error().Msg("error while unmarshalling ai model prompt from managerRequest")
		return dmerrors.DMError(apperr.ErrInvalidAIPromptRequestSpec, nil)
	}
	var err error
	prompt.SetPromptId()
	prompt.PreSaveCreate(v.CtxClaim)
	prompt.Status = mpb.CommonStatus_ACTIVE.String()
	prompt.Organization = v.CtxClaim[encryption.ClaimOrganizationKey]
	if prompt.Spec == nil {
		v.Logger.Error().Msg("error: invalid ai prompt spec, missing spec")
		return dmerrors.DMError(apperr.ErrInvalidAIPromptRequestSpec, nil)
	}
	if len(prompt.Spec.Tools) > 0 {
		err = BuildPromptSchema(prompt, v.Logger)
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while building prompt schema")
			return dmerrors.DMError(apperr.ErrInvalidAIPromptRequestSpec, err)
		}
	}
	BuildAIPromptTemplate(prompt)
	err = BuildPromptResponseFormat(prompt, v.Logger)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while building prompt response format")
		return dmerrors.DMError(apperr.ErrInvalidAIPromptRequestSpec, err)
	}
	err = BuildPromptVariables(prompt.Spec, v.Logger)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while building prompt variables")
		return dmerrors.DMError(apperr.ErrInvalidAIPromptRequestSpec, err)
	}
	err = v.dmStore.ConfigureAIPrompts(ctx, prompt, v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while saving ai prompt to datastore")
		return dmerrors.DMError(apperr.ErrSavingAIPrompt, err)
	}
	v.SetCreateResponse(mpb.Resources_AIPROMPTS, prompt.VapusID)
	return nil
}

func (v *AIPromptIntAgent) updateAIPrompts(ctx context.Context) error {
	pbObj := v.managerRequest.GetSpec()
	if pbObj == nil {
		v.Logger.Error().Msg("error while unmarshalling ai model prompt from managerRequest")
		return dmerrors.DMError(apperr.ErrInvalidAIPromptRequestSpec, nil)
	}
	var result []*models.AIPrompt
	exprompt, err := v.dmStore.GetAIPrompt(ctx, pbObj.PromptId, v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msgf("error while getting ai model prompt with id - %s from datastore for updateing", pbObj.PromptId)
		return dmerrors.DMError(apperr.ErrAIPrompt404, err)
	}
	if exprompt.Status != mpb.CommonStatus_ACTIVE.String() || exprompt.DeletedAt != 0 {
		v.Logger.Error().Msg("ai prompt is not active")
		return dmerrors.DMError(apperr.ErrResourceNotActive, nil)
	}
	if !(exprompt.CreatedBy == v.CtxClaim[encryption.ClaimUserIdKey] && exprompt.Organization == v.CtxClaim[encryption.ClaimOrganizationKey]) {
		v.Logger.Error().Msg("error while validating user access")
		return dmerrors.DMError(apperr.ErrPrompt403, nil)
	}
	// if !exprompt.Editable {
	// 	v.Logger.Error().Msgf("ai model prompt with id - %s is not editable", pbObj.PromptId)
	// 	return dmerrors.DMError(apperr.ErrAIPromptNotEditable, nil)
	// }
	newPrompt := (&models.AIPrompt{}).ConvertFromPb(pbObj)
	newPrompt.VapusBase = exprompt.VapusBase
	newPrompt.PreSaveUpdate(v.CtxClaim[encryption.ClaimUserIdKey])
	newPrompt.Organization = exprompt.Organization
	newPrompt.Status = mpb.CommonStatus_ACTIVE.String()
	if len(newPrompt.Spec.Tools) > 0 {
		err = BuildPromptSchema(newPrompt, v.Logger)
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while updating prompt schema")
			return dmerrors.DMError(apperr.ErrInvalidAIPromptRequestSpec, err)
		}
	}
	BuildAIPromptTemplate(newPrompt)
	err = BuildPromptResponseFormat(newPrompt, v.Logger)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while updating prompt response format")
		return dmerrors.DMError(apperr.ErrInvalidAIPromptRequestSpec, err)
	}
	err = BuildPromptVariables(newPrompt.Spec, v.Logger)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while updating prompt variables")
		return dmerrors.DMError(apperr.ErrInvalidAIPromptRequestSpec, err)
	}
	err = v.dmStore.PutAIPrompts(ctx, newPrompt, v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while updating ai prompt to datastore")
		return dmerrors.DMError(apperr.ErrSavingAIPrompt, err)
	}
	result = append(result, newPrompt)
	v.result = utils.AIPRArrObj2Pb(result)
	return nil
}

func (v *AIPromptIntAgent) describeAIPrompts(ctx context.Context) error {
	v.Logger.Info().Msg("getting ai model prompt describe from datastore")
	result, err := v.dmStore.GetAIPrompt(ctx, v.getterRequest.GetPromptId(), v.CtxClaim)
	if err != nil || result == nil {
		v.Logger.Error().Err(err).Msg("error while getting ai model prompt describe from datastore")
		return dmerrors.DMError(apperr.ErrAIPrompt404, err)
	}
	v.result = []*mpb.AIPrompt{result.ConvertToPb()}
	return nil
}

func (v *AIPromptIntAgent) listAIPrompts(ctx context.Context) error {
	v.Logger.Info().Msg("getting ai model prompt list from datastore")
	result, err := v.dmStore.ListAIPrompts(ctx,
		apppkgss.ListResourceWithGovernance(v.CtxClaim), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai model prompt list from datastore")
		return dmerrors.DMError(apperr.ErrAIPrompt404, err)
	}
	v.result = utils.AIPRListObj2Pb(result)
	return nil
}

func (v *AIPromptIntAgent) archiveAIPrompt(ctx context.Context) error {
	if v.getterRequest.GetPromptId() == "" {
		v.Logger.Error().Msg("error: invalid ai prompt id")
		return dmerrors.DMError(apperr.ErrAIPrompt404, nil)
	}
	result, err := v.dmStore.GetAIPrompt(ctx, v.getterRequest.GetPromptId(), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai model prompt describe from datastore")
		return dmerrors.DMError(err, nil)
	}
	result.Status = mpb.CommonStatus_DELETED.String()
	result.PreSaveDelete(v.CtxClaim)
	err = v.dmStore.PutAIPrompts(ctx, result, v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while archiving ai model node")
		return dmerrors.DMError(err, nil)
	}
	v.result = []*mpb.AIPrompt{}
	return nil
}
