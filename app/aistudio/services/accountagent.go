package services

import (
	"context"
	"slices"
	"strings"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	appdrepo "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo"
	aidmstore "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo/aistudio"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	processes "github.com/vapusdata-ecosystem/vapusai/core/process"
	types "github.com/vapusdata-ecosystem/vapusai/core/types"
)

type AccountAgent struct {
	request      *pb.AccountManagerRequest
	response     *models.Account
	organization *models.Organization
	*processes.VapusInterfaceBase
	dmStore *aidmstore.AIStudioDMStore
}

func (dms *AIStudioServices) GetAccount(ctx context.Context) (*models.Account, error) {
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		dms.Logger.Error().Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}
	_, err := dms.DMStore.GetOrganization(ctx, vapusPlatformClaim[encryption.ClaimOrganizationKey], vapusPlatformClaim)
	if err != nil {
		return nil, dmerrors.DMError(apperr.ErrOrganization404, err)
	}
	resp, err := dms.DMStore.GetAccount(ctx, vapusPlatformClaim)
	if err != nil {
		return nil, dmerrors.DMError(apperr.ErrGetAccount, err)
	}
	return resp, nil
}

func (dms *AIStudioServices) NewAccountAgent(ctx context.Context, request *pb.AccountManagerRequest) (*AccountAgent, error) {
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		dms.Logger.Error().Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}
	organization, err := dms.DMStore.GetOrganization(ctx, vapusPlatformClaim[encryption.ClaimOrganizationKey], vapusPlatformClaim)
	if err != nil {
		return nil, dmerrors.DMError(apperr.ErrOrganization404, err)
	}
	if organization.OrganizationType != mpb.OrganizationType_SERVICE_ORGANIZATION.String() {
		return nil, dmerrors.DMError(apperr.ErrNotServiceOrganization, nil)
	}
	agent := &AccountAgent{
		request:      request,
		organization: organization,
		VapusInterfaceBase: &processes.VapusInterfaceBase{
			InitAt: dmutils.GetEpochTime(),
			// Ctx:       ctx,
			CtxClaim:  vapusPlatformClaim,
			Action:    request.GetActions().String(),
			AgentType: types.ACCOUNTAGENT.String(),
		},
	}
	agent.SetAgentId()
	agent.Logger = pkgs.GetSubDMLogger(types.ACCOUNTAGENT.String(), agent.AgentId)
	return agent, nil
}

func (x *AccountAgent) Act(ctx context.Context, action string) error {
	if action != "" {
		x.Action = action
	}

	switch x.Action {
	case pb.AccountAgentActions_CONFIGURE_AISTUDIO_MODEL.String():
		response, err := x.configureAIAttributes(ctx)
		if err != nil {
			return err
		}
		x.response = response
		return nil
	case pb.AccountAgentActions_UPDATE_PROFILE.String():
		response, err := x.updateAccount(ctx)
		if err != nil {
			return err
		}
		x.response = response
		return nil
	default:
		return dmerrors.DMError(apperr.ErrInvalidAccountAgentActions, nil)
	}
}

func (x *AccountAgent) GetResponse() *models.Account {
	return x.response
}

func (x *AccountAgent) configureAIAttributes(ctx context.Context) (*models.Account, error) {
	userRoles := strings.Split(x.CtxClaim[encryption.ClaimRoleKey], ",")
	if !slices.Contains(userRoles, mpb.UserRoles_SERVICE_OWNER.String()) || !slices.Contains(userRoles, mpb.UserRoles_SERVICE_OPERATOR.String()) {
		return nil, dmerrors.DMError(apperr.ErrAccountOps403, nil)

	}
	var err error
	reqObj := (&models.Account{}).ConvertFromPb(x.request.GetSpec())
	account, err := x.dmStore.GetAccount(ctx, x.CtxClaim)
	if err != nil {
		return nil, dmerrors.DMError(apperr.ErrGetAccount, err)
	}
	account.AIAttributes = reqObj.GetAiAttributes()
	account.PreSaveUpdate(x.CtxClaim[encryption.ClaimUserIdKey])
	err = x.dmStore.PutAccount(ctx, account, x.CtxClaim)
	if err != nil {
		return nil, dmerrors.DMError(apperr.ErrConfigureAIStudioModel, err)
	}
	res := appdrepo.BootAccountCache(x.dmStore.VapusStore, x.Logger)
	if res == nil {
		x.Logger.Fatal().Msg("error while booting account cache")
	}
	for _, acc := range res {
		account = acc
		break
	}
	return account, nil
}

func (x *AccountAgent) updateAccount(ctx context.Context) (*models.Account, error) {
	userRoles := strings.Split(x.CtxClaim[encryption.ClaimRoleKey], ",")
	if !slices.Contains(userRoles, mpb.UserRoles_SERVICE_OWNER.String()) || x.organization.OrganizationType != mpb.OrganizationType_SERVICE_ORGANIZATION.String() {
		return nil, dmerrors.DMError(apperr.ErrAccountOps403, nil)
	}
	var err error
	reqObj := (&models.Account{}).ConvertFromPb(x.request.GetSpec())
	account, err := x.dmStore.GetAccount(ctx, x.CtxClaim)
	if err != nil {
		return nil, dmerrors.DMError(apperr.ErrGetAccount, err)
	}
	account.Profile = reqObj.Profile
	account.AIAttributes = reqObj.GetAiAttributes()

	account.PreSaveUpdate(x.CtxClaim[encryption.ClaimUserIdKey])
	account.ArtifactStorage = reqObj.GetArtifactStorage()

	err = x.dmStore.PutAccount(ctx, account, x.CtxClaim)
	if err != nil {
		return nil, dmerrors.DMError(apperr.ErrConfigureAIStudioModel, err)
	}
	res := appdrepo.BootAccountCache(x.dmStore.VapusStore, x.Logger)
	if res == nil {
		x.Logger.Fatal().Msg("error while booting account cache")
	}
	for _, acc := range res {
		account = acc
		break
	}
	return account, nil
}

func (x *AccountAgent) LogAgent() {
	x.Logger.Info().Msgf("AccountAgent - %v action started at %v and finished at %v with status %v", x.AgentId, x.InitAt, x.FinishAt, x.Status)
}
