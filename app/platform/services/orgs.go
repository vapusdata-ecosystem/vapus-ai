package services

import (
	"context"
	"fmt"
	"slices"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	processes "github.com/vapusdata-ecosystem/vapusdata/core/process"
	types "github.com/vapusdata-ecosystem/vapusdata/core/types"
	dmstores "github.com/vapusdata-ecosystem/vapusdata/platform/datastoreops"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/platform/pkgs"
	utils "github.com/vapusdata-ecosystem/vapusdata/platform/utils"
)

type DMIntAgentOpts func(*OrganizationAgent)

func WithDmAgentManagerRequest(managerRequest *pb.OrganizationManagerRequest) DMIntAgentOpts {
	return func(v *OrganizationAgent) {
		v.managerRequest = managerRequest
	}
}

func WithDmAgentUserManagerRequest(userAddRequest *pb.OrganizationAdduserRequest) DMIntAgentOpts {
	return func(v *OrganizationAgent) {
		v.userAddRequest = userAddRequest
	}
}

func WithDmAgentGetterRequest(getterRequest *pb.OrganizationGetterRequest) DMIntAgentOpts {
	return func(v *OrganizationAgent) {
		v.getterRequest = getterRequest
	}
}

func WithDmAgentManagerAction(action string) DMIntAgentOpts {
	return func(v *OrganizationAgent) {
		v.Action = action
	}
}

type OrganizationAgent struct {
	result         *pb.OrganizationResponse
	dbStore        *dmstores.DMStore
	DMServices     *DMServices
	managerRequest *pb.OrganizationManagerRequest
	userAddRequest *pb.OrganizationAdduserRequest
	getterRequest  *pb.OrganizationGetterRequest
	organization   *models.Organization
	errors         []error
	*processes.VapusInterfaceBase
}

func (x *OrganizationAgent) GetResult() *pb.OrganizationResponse {
	x.FinishAt = dmutils.GetEpochTime()
	return x.result
}

func (x *OrganizationAgent) LogAgent() {
	x.Logger.Info().Msgf("OrganizationAgent - %v action started at %v and finished at %v with status %v", x.AgentId, x.InitAt, x.FinishAt, x.Status)
}

func (x *DMServices) NewOrganizationAgent(ctx context.Context, opts ...DMIntAgentOpts) (*OrganizationAgent, error) {
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		x.logger.Error().Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}

	organization, err := x.DMStore.GetOrganization(ctx, vapusPlatformClaim[encryption.ClaimOrganizationKey], vapusPlatformClaim)
	if err != nil {
		return nil, dmerrors.DMError(apperr.ErrOrganization404, err)
	}

	agent := &OrganizationAgent{
		dbStore:      x.DMStore,
		DMServices:   x,
		organization: organization,
		result:       &pb.OrganizationResponse{Output: &pb.OrganizationResponse_OrganizationOutput{Users: make([]*pb.OrganizationResponse_OrganizationOutput_OrganizationUsers, 0)}},
		VapusInterfaceBase: &processes.VapusInterfaceBase{
			InitAt: dmutils.GetEpochTime(),
			// Ctx:       ctx,
			CtxClaim:  vapusPlatformClaim,
			AgentType: types.ORGANIZATIONAGENT.String(),
		},
	}
	for _, opt := range opts {
		opt(agent)
	}
	agent.SetAgentId()
	agent.Logger = pkgs.GetSubDMLogger(types.ORGANIZATIONAGENT.String(), agent.AgentId)
	return agent, nil
}

func (x *OrganizationAgent) Act(ctx context.Context, action string) error {

	if action != "" {
		x.Action = action
	}
	switch x.Action {
	case mpb.ResourceLcActions_ADD.String():
		x.organization = utils.DmNodeToObj(x.managerRequest)
		return x.configureOrganization(ctx)
	case mpb.ResourceLcActions_LIST.String():
		return x.listOrganizations(ctx)
	case mpb.ResourceLcActions_UPGRADE.String():
		return x.ugradeOrganizationArtifacts(ctx)
	case mpb.ResourceLcActions_UPDATE.String():
		x.organization = utils.DmNodeToObj(x.managerRequest)
		return x.patchOrganization(ctx)
	case mpb.ResourceLcActions_ADD_USERS.String():
		if x.userAddRequest != nil && x.userAddRequest.GetUsers() != nil {
			return x.addOrganizationUsers(ctx)
		}
		return dmerrors.DMError(apperr.ErrInvalidManageAgentActions, nil) //nolint:wrapcheck
	case mpb.ResourceLcActions_GET.String():
		return x.describeOrganization(ctx)
	default:
		return dmerrors.DMError(apperr.ErrInvalidManageAgentActions, nil) //nolint:wrapcheck
	}
}

func (x *OrganizationAgent) addOrganizationUsers(ctx context.Context) error {
	if x.organization.VapusID != x.CtxClaim[encryption.ClaimOrganizationKey] {
		return dmerrors.DMError(apperr.ErrInvalidOrganizationRequested, nil)
	}
	for _, user := range x.userAddRequest.GetUsers() {
		if !slices.Contains(x.organization.Users, user.GetUserId()) {
			err := x.attachOrganization2User(ctx, user.GetUserId(), user.InviteIfNotFound, &models.UserOrganizationRole{
				OrganizationId: x.organization.VapusID,
				RoleArns:       user.GetRole(),
			})
			if err != nil {
				x.Logger.Err(err).Ctx(ctx).Msgf("error while mapping user %s to this organization %v", user, x.organization)
				return dmerrors.DMError(apperr.ErrUserOrganizationMapping, err) //nolint:wrapcheck
			}
		}
	}
	return nil
}

func (x *OrganizationAgent) listOrganizationUsers(ctx context.Context) error {
	var err error
	organizationObj, err := x.dbStore.GetOrganization(ctx, x.organization.VapusID, x.CtxClaim)
	if err != nil {
		return dmerrors.DMError(apperr.ErrInvalidOrganizationRequested, err) //nolint:wrapcheck
	}
	filter := fmt.Sprintf(`organization_roles @> '[{"organizationId": "%s"}]'`, organizationObj.VapusID)
	users, err := x.dbStore.ListUsers(ctx, filter, x.CtxClaim)
	if err != nil {
		return dmerrors.DMError(apperr.ErrInvalidOrganizationRequested, err) //nolint:wrapcheck
	}
	x.result.Output.Users = []*pb.OrganizationResponse_OrganizationOutput_OrganizationUsers{{
		Users:        utils.DmUListingToPb(users, x.organization.VapusID),
		Organization: organizationObj.VapusID,
	}}
	x.result.Output.Organizations = utils.DmArrToPb([]*models.Organization{organizationObj})
	return nil
}

func (x *OrganizationAgent) configureOrganization(ctx context.Context) error {
	var err error
	if x.organization == nil {
		return dmerrors.DMError(apperr.ErrInvalidAddOrganizationRequest, nil)
	}
	if x.organization.OrganizationType == mpb.OrganizationType_SERVICE_ORGANIZATION.String() {
		return dmerrors.DMError(apperr.ErrCannotCreateServiceOrganization, nil)
	}
	x.organization, err = organizationConfigureTool(ctx, x.organization, x.dbStore, x.Logger, x.CtxClaim)
	if err != nil {
		return dmerrors.DMError(err, nil)
	}
	err = x.attachOrganization2User(ctx, x.CtxClaim[encryption.ClaimUserIdKey], true, &models.UserOrganizationRole{
		OrganizationId: x.organization.VapusID,
		RoleArns:       []string{mpb.UserRoles_ORG_USER.String()},
	})
	if err != nil {
		x.Logger.Err(err).Ctx(ctx).Msgf("error while mapping user as organization owner to this organization %v", x.organization)
		return dmerrors.DMError(apperr.ErrUserOrganizationMapping, err) //nolint:wrapcheck
	}
	x.SetCreateResponse(mpb.Resources_ORGANIZATIONS, x.organization.VapusID)
	return nil
}

func (x *OrganizationAgent) attachOrganization2User(ctx context.Context, userId string, invite bool, obj *models.UserOrganizationRole) error {
	var user *models.Users
	var err error
	user, err = x.dbStore.GetUser(ctx, userId, x.CtxClaim)
	if err != nil {
		if !invite {
			return dmerrors.DMError(apperr.ErrInvalidUserRequested, err) //nolint:wrapcheck
		}
		x.Logger.Err(err).Ctx(ctx).Msgf("User with id %s not found, inviting.", userId)
		userAgent, err := x.DMServices.NewUserManagerAgent(ctx, &pb.UserManagerRequest{
			Action: pb.UserManagerActions_INVITE_USERS,
			Spec: &mpb.User{
				Email: userId,
			},
			Organization: obj.OrganizationId,
			RoleArn:      obj.RoleArns,
		}, nil)
		if err != nil {
			return dmerrors.DMError(apperr.ErrInvalidUserRequested, err) //nolint:wrapcheck
		}
		err = userAgent.Act(ctx, pb.UserManagerActions_INVITE_USERS.String())
		if err != nil {
			return dmerrors.DMError(apperr.ErrInvalidUserRequested, err) //nolint:wrapcheck
		}
		userAgent.LogAgent()
		result := userAgent.GetResult()
		if len(result.GetOutput().Users) > 0 {
			user, err = x.dbStore.GetUser(ctx, result.GetOutput().Users[0].UserId, x.CtxClaim)
			if err != nil {
				return dmerrors.DMError(apperr.ErrInvalidUserRequested, err) //nolint:wrapcheck
			}
		}
	}
	return x.dbStore.PutUser(ctx, user, x.CtxClaim)
}

func (x *OrganizationAgent) ugradeOrganizationArtifacts(ctx context.Context) error {
	x.Logger.Debug().Ctx(ctx).Msgf("ugradeOrganizationArtifacts for organization %v", x.organization)
	var err error
	organization, err := x.dbStore.GetOrganization(ctx, x.CtxClaim[encryption.ClaimOrganizationKey], x.CtxClaim)
	if err != nil {
		return dmerrors.DMError(apperr.ErrInvalidOrganizationRequested, err) //nolint:wrapcheck
	}
	organization.PreSaveUpdate(x.CtxClaim[encryption.ClaimUserIdKey])
	err = x.dbStore.PutOrganization(ctx, organization, x.CtxClaim)
	if err != nil {
		x.Logger.Err(err).Ctx(ctx).Msgf("error while configuring organization %v", organization)
		return dmerrors.DMError(apperr.ErrCreateOrganization, err) //nolint:wrapcheck
	}

	x.result.Output.Organizations = utils.DmArrToPb([]*models.Organization{organization})
	return nil
}

func (x *OrganizationAgent) patchOrganization(ctx context.Context) error {
	var err error
	newObj := utils.DmNodeToObj(x.managerRequest)
	if x.organization.OrganizationType == mpb.OrganizationType_SERVICE_ORGANIZATION.String() {
		return dmerrors.DMError(apperr.ErrCannotCreateServiceOrganization, nil)
	}
	x.organization.PreSaveUpdate(x.CtxClaim[encryption.ClaimUserIdKey])
	for _, nwUser := range x.userAddRequest.GetUsers() {
		if !slices.Contains(x.organization.Users, nwUser.GetUserId()) {
			err := x.attachOrganization2User(ctx, nwUser.GetUserId(), nwUser.InviteIfNotFound, &models.UserOrganizationRole{
				OrganizationId: x.organization.VapusID,
				RoleArns:       nwUser.GetRole(),
			})
			if err != nil {
				x.Logger.Err(err).Ctx(ctx).Msgf("error while mapping user %s to this organization %v", nwUser, x.organization)
			}
		}
	}
	if newObj != nil && newObj.DataProductInfraPlatform != nil {
		err = setOrganizationDPK8sInfra(ctx, x.organization, newObj.DataProductInfraPlatform, x.dbStore)
		if err != nil {
			return dmerrors.DMError(err, nil)
		}
	}
	if newObj != nil && newObj.ArtifactStorage != nil {
		resp, err := setOrganizationArtifactBEStore(ctx, newObj, x.dbStore)
		if err != nil {
			x.Logger.Err(err).Ctx(ctx).Msgf("error while setting organization artifact store %v", x.organization)
			return dmerrors.DMError(apperr.ErrSettingOrganizationArtifactStore, err) //nolint:wrapcheck
		}
		if resp != nil {
			x.organization.ArtifactStorage = resp
		}
	}
	if newObj != nil {
		x.organization.DataProductInfraPlatform = newObj.DataProductInfraPlatform
	}
	x.organization.DisplayName = newObj.DisplayName
	err = x.dbStore.PutOrganization(ctx, x.organization, x.CtxClaim)
	if err != nil {
		x.Logger.Err(err).Ctx(ctx).Msgf("error while configuring organization %v", x.organization)
		return dmerrors.DMError(apperr.ErrCreateOrganization, err) //nolint:wrapcheck
	}

	x.result.Output.Organizations = utils.DmArrToPb([]*models.Organization{x.organization})
	return nil
}

func (x *OrganizationAgent) listOrganizations(ctx context.Context) error {
	var filter string = ""
	dmIds := utils.GetFilterParams(x.getterRequest.GetSearchParam(), types.ORGANIZATIONSK.String())
	if len(dmIds) > 0 {
		filter = fmt.Sprintf("vapus_id IN (%s)", dmIds)
	} else {
		filter = ""
	}
	organizations, err := x.dbStore.ListOrganizations(ctx, filter, x.CtxClaim)
	if err != nil {
		return dmerrors.DMError(apperr.ErrInvalidOrganizationRequested, err) //nolint:wrapcheck
	}
	x.result.Output.Organizations = utils.DmListToPb(organizations)
	return nil
}

func (x *OrganizationAgent) describeOrganization(ctx context.Context) error {
	organization, err := x.dbStore.GetOrganization(ctx, x.CtxClaim[encryption.ClaimOrganizationKey], x.CtxClaim)
	if err != nil {
		return dmerrors.DMError(apperr.ErrInvalidOrganizationRequested, err) //nolint:wrapcheck
	}
	x.result.Output.Organizations = utils.DmArrToPb([]*models.Organization{organization})
	return nil
}
