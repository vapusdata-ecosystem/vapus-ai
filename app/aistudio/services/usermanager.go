package services

import (
	"context"
	"fmt"
	"slices"
	"strings"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	dmstores "github.com/vapusdata-ecosystem/vapusai/aistudio/datastoreops"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	utils "github.com/vapusdata-ecosystem/vapusai/aistudio/utils"
	aidmstore "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo/aistudio"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	models "github.com/vapusdata-ecosystem/vapusai/core/models"
	"github.com/vapusdata-ecosystem/vapusai/core/options"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	processes "github.com/vapusdata-ecosystem/vapusai/core/process"
	types "github.com/vapusdata-ecosystem/vapusai/core/types"
)

type UserManagerAgent struct {
	managerRequest *pb.UserManagerRequest
	getterRequest  *pb.UserGetterRequest
	dmStore        *aidmstore.AIStudioDMStore
	user           *models.Users
	result         *pb.UserResponse
	organization   *models.Organization
	*processes.VapusInterfaceBase
}

func (x *UserManagerAgent) GetResult() *pb.UserResponse {
	x.FinishAt = dmutils.GetEpochTime()
	return x.result
}

func (x *UserManagerAgent) LogAgent() {
	x.Logger.Info().Msgf("UserManagerAgent - %v action started at %v and finished at %v with status %v", x.AgentId, x.InitAt, x.FinishAt, x.Status)
}

func (x *AIStudioServices) NewUserManagerAgent(ctx context.Context, managerRequest *pb.UserManagerRequest, getterRequest *pb.UserGetterRequest) (*UserManagerAgent, error) {
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		x.Logger.Error().Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}

	organization, err := x.DMStore.GetOrganization(ctx, vapusPlatformClaim[encryption.ClaimOrganizationKey], vapusPlatformClaim)
	if err != nil {
		return nil, dmerrors.DMError(apperr.ErrOrganization404, err)
	}
	agent := &UserManagerAgent{
		result:         &pb.UserResponse{Output: &pb.UserResponse_VapusUser{}},
		managerRequest: managerRequest,
		getterRequest:  getterRequest,
		dmStore:        x.DMStore,
		organization:   organization,
		VapusInterfaceBase: &processes.VapusInterfaceBase{
			InitAt: dmutils.GetEpochTime(),
			// Ctx:       ctx,
			CtxClaim:  vapusPlatformClaim,
			AgentType: types.USERMANAGERAGENT.String(),
		},
	}
	agent.SetAgentId()
	if managerRequest != nil {
		agent.Action = managerRequest.GetAction().String()
	} else if getterRequest != nil {
		if getterRequest.GetUserId() == "" {
			agent.Action = getterRequest.GetAction().String()
		} else {
			agent.Action = pb.UserGetterActions_GET_USER.String()
		}
	} else {
		agent.Action = ""
	}
	agent.Logger = pkgs.GetSubDMLogger(types.DATAPRODUCTAGENT.String(), agent.AgentId)
	return agent, nil
}

func (x *UserManagerAgent) Act(ctx context.Context, action string) error {
	if action != "" {
		x.Action = action
	}
	switch x.Action {
	case pb.UserManagerActions_INVITE_USERS.String():
		return x.InviteUser(ctx)
	case pb.UserManagerActions_PATCH_USER.String():
		return x.PatchUser(ctx)
	case pb.UserGetterActions_GET_USER.String():
		return x.GetUser(ctx, x.getterRequest.GetUserId())
	case pb.UserGetterActions_LIST_USERS.String():
		return x.ListUsers(ctx)
	case pb.UserGetterActions_LIST_PLATFORM_USERS.String():
		return x.ListPlatformUsers(ctx)
	default:
		return dmerrors.DMError(apperr.ErrInvalidUserManagerAction, nil)
	}
}

func (x *UserManagerAgent) InviteUser(ctx context.Context) error {
	var user *models.Users
	var err error
	if x.managerRequest == nil {
		x.Logger.Error().Msg("error while getting user invite request")
		return dmerrors.DMError(apperr.ErrUserInviteCreateFailed, nil)
	}
	userObj := utils.DmUPbToObj(x.managerRequest.GetSpec())
	if userObj == nil {
		x.Logger.Error().Msg("error while converting user object from pb")
		return dmerrors.DMError(apperr.ErrUserInviteCreateFailed, nil)
	}
	for _, organizationRole := range userObj.Roles {
		organizationRole.InvitedOn = dmutils.GetEpochTime()
	}
	exists := x.dmStore.UserInviteExists(ctx, userObj.Email, x.CtxClaim)
	if exists {
		user, err := x.dmStore.GetUser(ctx, userObj.Email, x.CtxClaim)
		if err != nil {
			x.Logger.Error().Msgf("invite already exists for email - %v, but not registered user", userObj.Email)
			return dmerrors.DMError(apperr.ErrUserInviteExists, nil)
		}
		for _, organizationRole := range userObj.Roles {
			for _, existOrganizationRole := range user.Roles {
				if organizationRole.OrganizationId == existOrganizationRole.OrganizationId {
					existOrganizationRole.RoleArns = append(existOrganizationRole.RoleArns, organizationRole.RoleArns...)
				} else {

					user.Roles = append(user.Roles, organizationRole)
				}
			}
		}
		user.PreSaveUpdate(x.CtxClaim[encryption.ClaimUserIdKey])
		err = x.dmStore.PutUser(ctx, user, x.CtxClaim)
		if err != nil {
			x.Logger.Error().Msgf("error while updating user - %v", err)
			return dmerrors.DMError(apperr.ErrUserInviteCreateFailed, err)
		}
	} else {
		userObj.InvitedType = mpb.UserInviteType_INVITE_ACCESS.String()
		userObj.SetUserId()
		userObj.Status = mpb.CommonStatus_INVITED.String()
		userObj.PreSaveInvite(x.CtxClaim, utils.DEFAULT_USER_INVITE_EXPIRE_TIME) // fetch user rom MD
		user, err = x.dmStore.CreateUser(ctx, nil, userObj, x.CtxClaim)
		if err != nil {
			x.Logger.Error().Msgf("error while creating user while inviting - %v", err)
			return dmerrors.DMError(apperr.ErrUserInviteCreateFailed, err)
		}
		emailer := pkgs.PluginServiceManager.PlatformPlugins.Emailer
		body := types.UserInviteEMailTemplate
		body = strings.ReplaceAll(body, "{Account}", dmstores.DMStoreManager.Account.Name)
		body = strings.ReplaceAll(body, "{Name}", userObj.FirstName+" "+userObj.LastName)
		body = strings.ReplaceAll(body, "{Link}", pkgs.NetworkConfigManager.ExternalURL+"/login")
		err = emailer.SendRawEmail(ctx, &options.SendEmailRequest{
			To:               []string{userObj.Email},
			Subject:          "Welcome to VapusData Platform",
			HtmlTemplateBody: body,
		}, x.AgentId)
		if err != nil {
			x.Logger.Error().Msgf("error while sending email on user invite - %v", err)
			return dmerrors.DMError(apperr.ErrUserInviteCreateFailed, err)
		}
	}
	x.result.Output.Users = utils.DmUArToPb([]*models.Users{user}, x.CtxClaim[encryption.ClaimOrganizationKey])
	return nil
}

func (x *UserManagerAgent) GetUser(ctx context.Context, userId string) error {
	if userId == "" {
		userId = x.CtxClaim[encryption.ClaimUserIdKey]
	}
	user, err := x.dmStore.GetUser(ctx, userId, x.CtxClaim)
	if err != nil {
		x.Logger.Error().Ctx(ctx).Msgf("error while getting user - %v", err)
		return dmerrors.DMError(apperr.ErrUser404, err)
	}
	organizationMap := make(map[string]string)
	organizationIds := ""
	for _, organizationRole := range user.Roles {
		organizationIds = fmt.Sprintf("%s'%s',", organizationIds, organizationRole.OrganizationId)
	}
	organizationIds = strings.TrimRight(organizationIds, ",")
	organizations, err := x.dmStore.ListOrganizations(ctx,
		"vapus_id in ("+organizationIds+")", x.CtxClaim)
	if err != nil {
		x.Logger.Error().Ctx(ctx).Msgf("error while getting organization - %v", err)
		return dmerrors.DMError(apperr.ErrUser404, err)
	}
	for _, organization := range organizations {
		organizationMap[organization.VapusID] = organization.Name
	}
	x.result.OrganizationMap = organizationMap
	if user.IsValidUserByOrganization(x.CtxClaim[encryption.ClaimOrganizationKey]) {
		x.result.Output.Users = utils.DmUArToPb([]*models.Users{user}, x.CtxClaim[encryption.ClaimOrganizationKey])
		return nil
	} else if user.Organization == x.CtxClaim[encryption.ClaimOrganizationKey] {
		dmRole := user.GetOrganizationRole(x.CtxClaim[encryption.ClaimOrganizationKey])
		if len(dmRole) > 0 && slices.Contains(dmRole[0].RoleArns, mpb.UserRoles_ORG_USER.String()) {
			x.result.Output.Users = utils.DmUArToPb([]*models.Users{user}, x.CtxClaim[encryption.ClaimOrganizationKey])
			return nil
		} else {
			x.Logger.Error().Msg("error while getting organization role from user object")
			return dmerrors.DMError(apperr.ErrUser404, nil)
		}
	} else {
		x.Logger.Error().Ctx(ctx).Msgf("error while getting user - %v", err)
		return dmerrors.DMError(apperr.ErrUser404, err)
	}
}

func (x *UserManagerAgent) ListUsers(ctx context.Context) error {
	users, err := x.dmStore.GetOrganizationUsers(ctx, x.CtxClaim[encryption.ClaimOrganizationKey], x.CtxClaim)
	if err != nil {
		x.Logger.Error().Ctx(ctx).Msgf("error while getting user - %v", err)
		return dmerrors.DMError(apperr.ErrUser404, err)
	}
	x.result.Output.Users = utils.DmUListingToPb(users, x.CtxClaim[encryption.ClaimOrganizationKey])
	return nil
}

func (x *UserManagerAgent) ListPlatformUsers(ctx context.Context) error {
	users, err := x.dmStore.ListUsers(ctx, "", x.CtxClaim)
	if err != nil {
		x.Logger.Error().Ctx(ctx).Msgf("error while getting user - %v", err)
		return dmerrors.DMError(apperr.ErrUser404, err)
	}
	x.result.Output.Users = utils.DmUListingToPb(users, x.CtxClaim[encryption.ClaimOrganizationKey])
	return nil
}

func (x *UserManagerAgent) PatchUser(ctx context.Context) error {
	userObj := utils.DmUPbToObj(x.managerRequest.GetSpec())
	if userObj == nil || userObj.UserId == "" {
		x.Logger.Error().Msg("error while converting user object from pb")
		return dmerrors.DMError(apperr.ErrUser404, nil)
	}
	exUser, err := x.dmStore.GetUser(ctx, userObj.UserId, x.CtxClaim)
	if err != nil {
		x.Logger.Error().Ctx(ctx).Msgf("error while getting user - %v", err)
		return dmerrors.DMError(apperr.ErrUser404, err)
	}
	if x.CtxClaim[encryption.ClaimUserIdKey] != exUser.UserId {
		if exUser.GetOrganizationRole(x.CtxClaim[encryption.ClaimOrganizationKey]) == nil {
			x.Logger.Error().Msg("error while getting organization role from user object")
			return dmerrors.DMError(apperr.ErrUser404, nil)
		}
		if !strings.Contains(x.CtxClaim[encryption.ClaimRoleKey], mpb.UserRoles_ORG_OWNER.String()) {
			x.Logger.Error().Msg("error while getting organization role from user object")
			return dmerrors.DMError(apperr.ErrUser404, nil)
		}
	}
	updatedRoles := []string{}
	for _, organizationRole := range userObj.Roles {
		if organizationRole.OrganizationId == x.CtxClaim[encryption.ClaimOrganizationKey] {
			updatedRoles = append(updatedRoles, organizationRole.RoleArns...)
		}
	}
	for _, organizationRole := range exUser.Roles {
		if organizationRole.OrganizationId == x.CtxClaim[encryption.ClaimOrganizationKey] {
			organizationRole.RoleArns = updatedRoles
		}
	}
	exUser.PreSaveUpdate(x.CtxClaim[encryption.ClaimUserIdKey])

	if userObj.FirstName != "" {
		exUser.FirstName = userObj.FirstName
	}
	if userObj.LastName != "" {
		exUser.LastName = userObj.LastName
	}
	if userObj.DisplayName != "" {
		exUser.DisplayName = userObj.DisplayName
	}
	if userObj.Profile != nil {
		if exUser.Profile.Addresses != nil {
			exUser.Profile.Addresses = userObj.Profile.Addresses
		}
		if userObj.Profile.Avatar != "" {
			exUser.Profile.Avatar = userObj.Profile.Avatar
		}
		if userObj.Profile.DateOfBirth != "" {
			exUser.Profile.DateOfBirth = userObj.Profile.DateOfBirth
		}
		if userObj.Profile.Description != "" {
			exUser.Profile.Description = userObj.Profile.Description
		}
		if userObj.Profile.Gender != "" {
			exUser.Profile.Gender = userObj.Profile.Gender
		}
	}
	err = x.dmStore.PutUser(ctx, exUser, x.CtxClaim)
	x.result.Output.Users = utils.DmUArToPb([]*models.Users{exUser}, x.CtxClaim[encryption.ClaimOrganizationKey])
	return nil
}
