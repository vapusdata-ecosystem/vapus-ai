package services

import (
	"context"
	"slices"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/aistudio/utils"
	aidmstore "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo/aistudio"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	apppkgss "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	models "github.com/vapusdata-ecosystem/vapusai/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	processes "github.com/vapusdata-ecosystem/vapusai/core/process"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

type AIGuardrailIntAgent struct {
	*processes.VapusInterfaceBase
	managerRequest *pb.GuardrailsManagerRequest
	getterRequest  *pb.GuardrailsGetterRequest
	result         []*mpb.AIGuardrails
	dmStore        *aidmstore.AIStudioDMStore
}

type AIGuardrailIntAgentOpts func(*AIGuardrailIntAgent)

func WithGuardrailManagerRequest(managerRequest *pb.GuardrailsManagerRequest) AIGuardrailIntAgentOpts {
	return func(v *AIGuardrailIntAgent) {
		v.managerRequest = managerRequest
	}
}

func WithGuardrailGetterRequest(getterRequest *pb.GuardrailsGetterRequest) AIGuardrailIntAgentOpts {
	return func(v *AIGuardrailIntAgent) {
		v.getterRequest = getterRequest
	}
}

func WithGuardrailManagerAction(action mpb.ResourceLcActions) AIGuardrailIntAgentOpts {
	return func(v *AIGuardrailIntAgent) {
		v.Action = action.String()
	}
}

func (s *AIStudioServices) NewAIGuardrailIntAgent(ctx context.Context, opts ...AIGuardrailIntAgentOpts) (*AIGuardrailIntAgent, error) {
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		s.Logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}
	guardrail := &AIGuardrailIntAgent{
		result:  make([]*mpb.AIGuardrails, 0),
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

func (v *AIGuardrailIntAgent) GetAgentId() string {
	return v.AgentId
}

func (v *AIGuardrailIntAgent) GetResult() []*mpb.AIGuardrails {
	v.FinishAt = dmutils.GetEpochTime()
	v.FinalLog()
	return v.result
}

func (v *AIGuardrailIntAgent) Act(ctx context.Context) error {
	switch v.GetAction() {
	case mpb.ResourceLcActions_CREATE.String():
		return v.configureAIGuardrails(ctx)
	case mpb.ResourceLcActions_UPDATE.String():
		return v.updateAIGuardrails(ctx)
	case mpb.ResourceLcActions_GET.String():
		return v.describeAIGuardrails(ctx)
	case mpb.ResourceLcActions_LIST.String():
		return v.listAIGuardrails(ctx)
	case mpb.ResourceLcActions_ARCHIVE.String():
		return v.archiveAIGuardrail(ctx)
	default:
		v.Logger.Error().Msg("invalid action for AIGuardrailIntAgent")
		return dmerrors.DMError(apperr.ErrInvalidAction, nil)
	}
}

func (v *AIGuardrailIntAgent) updateCachePool(guardrail *models.AIGuardrails) {
	pkgs.GuardrailPoolManager.UpdateGuardrailPool(guardrail)
	return
}

func (v *AIGuardrailIntAgent) GetResourceCreateResponse(guardrail *models.AIGuardrails) {
	pkgs.GuardrailPoolManager.UpdateGuardrailPool(guardrail)
	return
}

func (v *AIGuardrailIntAgent) configureAIGuardrails(ctx context.Context) error {
	aiGuardrails := (&models.AIGuardrails{}).ConvertFromPb(v.managerRequest.GetSpec())
	aiGuardrails.PreSaveCreate(v.CtxClaim)
	aiGuardrails.Status = mpb.CommonStatus_ACTIVE.String()
	aiGuardrails.Organization = v.CtxClaim[encryption.ClaimOrganizationKey]
	if aiGuardrails.Editors == nil {
		aiGuardrails.Editors = []string{v.CtxClaim[encryption.ClaimUserIdKey]}
	} else if !slices.Contains(aiGuardrails.Editors, v.CtxClaim[encryption.ClaimUserIdKey]) {
		aiGuardrails.Editors = append(aiGuardrails.Editors, v.CtxClaim[encryption.ClaimUserIdKey])
	}
	aiGuardrails.Schema = BuildGuardrailSchema(aiGuardrails)
	err := v.dmStore.ConfigureAIGuardrails(ctx, aiGuardrails, v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while configuring ai guardrail")
		return dmerrors.DMError(apperr.ErrAIGuardrailCreate400, err)
	}
	v.updateCachePool(aiGuardrails)
	v.SetCreateResponse(mpb.Resources_AIGUARDRAILS, aiGuardrails.VapusID)
	return nil
}

func (v *AIGuardrailIntAgent) updateAIGuardrails(ctx context.Context) error {
	existingAIGuardrail, err := v.dmStore.GetAIGuardrail(ctx, v.managerRequest.GetSpec().GetGuardrailId(), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai guardrail from datastore")
		return dmerrors.DMError(apperr.ErrAIGuardrail404, err)
	}
	if existingAIGuardrail.Status != mpb.CommonStatus_ACTIVE.String() || existingAIGuardrail.DeletedAt != 0 {
		v.Logger.Error().Msg("ai guardrail is not active")
		return dmerrors.DMError(apperr.ErrResourceNotActive, nil)
	}
	if !(existingAIGuardrail.CreatedBy == v.CtxClaim[encryption.ClaimUserIdKey] && existingAIGuardrail.Organization == v.CtxClaim[encryption.ClaimOrganizationKey]) {
		v.Logger.Error().Msg("error while validating user access to ai guardrail")
		return dmerrors.DMError(apperr.ErrAIGuardrail403, nil)
	}
	aiGuardrails := (&models.AIGuardrails{}).ConvertFromPb(v.managerRequest.GetSpec())
	existingAIGuardrail.PreSaveUpdate(v.CtxClaim[encryption.ClaimUserIdKey])

	existingAIGuardrail.Description = aiGuardrails.Description
	if aiGuardrails.Editors != nil {
		existingAIGuardrail.Editors = aiGuardrails.Editors
	}
	existingAIGuardrail.Contents = aiGuardrails.Contents
	existingAIGuardrail.Topics = aiGuardrails.Topics
	existingAIGuardrail.Words = aiGuardrails.Words
	existingAIGuardrail.SensitiveDataset = aiGuardrails.SensitiveDataset
	existingAIGuardrail.Editors = aiGuardrails.Editors
	existingAIGuardrail.Status = mpb.CommonStatus_ACTIVE.String()
	existingAIGuardrail.Organization = v.CtxClaim[encryption.ClaimOrganizationKey]
	existingAIGuardrail.Description = aiGuardrails.Description
	existingAIGuardrail.FailureMessage = aiGuardrails.FailureMessage
	existingAIGuardrail.Schema = BuildGuardrailSchema(existingAIGuardrail)
	existingAIGuardrail.ScanMode = aiGuardrails.ScanMode
	existingAIGuardrail.GuardModel = aiGuardrails.GuardModel
	err = v.dmStore.PutAIGuardrails(ctx, existingAIGuardrail, v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while configuring ai guardrail")
		return dmerrors.DMError(apperr.ErrAIGuardrailPatch400, err)
	}
	v.updateCachePool(existingAIGuardrail)
	v.result = []*mpb.AIGuardrails{existingAIGuardrail.ConvertToPb()}
	return nil
}

func (v *AIGuardrailIntAgent) describeAIGuardrails(ctx context.Context) error {
	v.Logger.Info().Msg("getting ai guardrail describe from datastore")
	result, err := v.dmStore.GetAIGuardrail(ctx, v.getterRequest.GetGuardrailId(), v.CtxClaim)
	if err != nil || result == nil {
		v.Logger.Error().Err(err).Msg("error while getting ai guardrail describe from datastore")
		return dmerrors.DMError(apperr.ErrAIGuardrail404, err)
	}
	v.result = []*mpb.AIGuardrails{result.ConvertToPb()}
	return nil
}

func (v *AIGuardrailIntAgent) listAIGuardrails(ctx context.Context) error {
	v.Logger.Info().Msg("getting ai guardrail list from datastore")
	result, err := v.dmStore.ListAIGuardrails(ctx,
		apppkgss.ListResourceWithGovernance(v.CtxClaim), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai guardrails list from datastore")
		return dmerrors.DMError(apperr.ErrAIGuardrail404, err)
	}
	v.result = utils.AIGDListObjToPb(result)
	return nil
}

func (v *AIGuardrailIntAgent) archiveAIGuardrail(ctx context.Context) error {
	if v.getterRequest.GetGuardrailId() == "" {
		v.Logger.Error().Msg("error: invalid ai GuardrailId")
		return dmerrors.DMError(apperr.ErrMissingResource, nil)
	}
	result, err := v.dmStore.GetAIGuardrail(ctx, v.getterRequest.GetGuardrailId(), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while archiveAIGuardrail from datastore")
		return dmerrors.DMError(err, nil)
	}
	result.Status = mpb.CommonStatus_DELETED.String()
	result.PreSaveDelete(v.CtxClaim)
	err = v.dmStore.PutAIGuardrails(ctx, result, v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while archiving ai guardrail")
		return dmerrors.DMError(err, nil)
	}
	v.result = []*mpb.AIGuardrails{}
	return nil
}
