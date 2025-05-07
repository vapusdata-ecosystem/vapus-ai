package services

import (
	"context"

	pb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	models "github.com/vapusdata-ecosystem/vapusdata/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	processes "github.com/vapusdata-ecosystem/vapusdata/core/process"
	types "github.com/vapusdata-ecosystem/vapusdata/core/types"
	dmstores "github.com/vapusdata-ecosystem/vapusdata/platform/datastoreops"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/platform/pkgs"
)

type AuthzManagerAgent struct {
	managerRequest *pb.AuthzManagerRequest
	getterRequest  *pb.AuthzGetterRequest
	dbStore        *dmstores.DMStore
	*DMServices
	user         *models.Users
	result       *pb.AuthzResponse
	organization *models.Organization
	*processes.VapusInterfaceBase
}

func (x *AuthzManagerAgent) GetResult() *pb.AuthzResponse {
	x.FinishAt = dmutils.GetEpochTime()
	return x.result
}

func (x *AuthzManagerAgent) LogAgent() {
	x.Logger.Info().Msgf("AuthzManagerAgent - %v action started at %v and finished at %v with status %v", x.AgentId, x.InitAt, x.FinishAt, x.Status)
}

func (x *DMServices) NewAuthzManagerAgent(ctx context.Context, managerRequest *pb.AuthzManagerRequest, getterRequest *pb.AuthzGetterRequest) (*AuthzManagerAgent, error) {
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		x.logger.Error().Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}

	organization, err := x.DMStore.GetOrganization(ctx, vapusPlatformClaim[encryption.ClaimOrganizationKey], vapusPlatformClaim)
	if err != nil {
		return nil, dmerrors.DMError(apperr.ErrOrganization404, err)
	}
	agent := &AuthzManagerAgent{
		result:         &pb.AuthzResponse{Output: &pb.AuthzResponse_AuthzRoles{}},
		managerRequest: managerRequest,
		getterRequest:  getterRequest,
		dbStore:        x.DMStore,
		DMServices:     x,
		organization:   organization,
		VapusInterfaceBase: &processes.VapusInterfaceBase{
			InitAt:    dmutils.GetEpochTime(),
			CtxClaim:  vapusPlatformClaim,
			AgentType: types.USERMANAGERAGENT.String(),
			// Ctx:       ctx,
		},
	}
	agent.SetAgentId()
	if managerRequest != nil {
		agent.Action = managerRequest.GetAction().String()
	} else if getterRequest != nil {
		if getterRequest.GetRoleArn() == "" {
			agent.Action = "pb.AuthzManagerRequest_GET_AUTHZ_ROLE.String()"
		} else {
			agent.Action = "pb.AuthzManagerRequest_LIST_AUTHZ_ROLES.String()"
		}
	} else {
		agent.Action = ""
	}
	agent.Logger = pkgs.GetSubDMLogger(types.DATAPRODUCTAGENT.String(), agent.AgentId)
	return agent, nil
}

func (x *AuthzManagerAgent) Act(action string) error {
	// switch x.Action {
	// case pb.AuthzManagerRequest_CREATE_AUTHZ_ROLE.String():
	// 	return x.CreateAuthzRole()
	// case pb.AuthzManagerRequest_GET_AUTHZ_ROLE.String():
	// 	return x.GetAuthzRole()
	// case pb.AuthzManagerRequest_UPDATE_AUTHZ_ROLE.String():
	// 	return x.UpdateAuthzRole()
	// case pb.AuthzManagerRequest_DELETE_AUTHZ_ROLE.String():
	// 	return x.DeleteAuthzRole()
	// case pb.AuthzManagerRequest_LIST_AUTHZ_ROLES.String():
	// 	return x.ListAuthzRoles()
	// default:
	// 	return dmerrors.DMError(apperr.ErrInvalidAction, nil)
	// }
	return nil
}
