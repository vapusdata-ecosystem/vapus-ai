package services

import (
	"context"
	"encoding/json"
	"fmt"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	processes "github.com/vapusdata-ecosystem/vapusdata/core/process"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/platform/pkgs"
	utils "github.com/vapusdata-ecosystem/vapusdata/platform/utils"
)

type SecretsManagerAgentOpts func(*SecretsManagerAgent)

func WithSecretManagerRequest(req *pb.SecretManagerRequest) SecretsManagerAgentOpts {
	return func(a *SecretsManagerAgent) {
		a.managerRequest = req
	}
}

func WithSecretGetterRequest(req *pb.SecretGetterRequest) SecretsManagerAgentOpts {
	return func(a *SecretsManagerAgent) {
		a.getterRequest = req
	}
}

func WithSecreManagerAction(action string) SecretsManagerAgentOpts {
	return func(a *SecretsManagerAgent) {
		a.Action = action
	}
}

type SecretsManagerAgent struct {
	*processes.VapusInterfaceBase
	managerRequest *pb.SecretManagerRequest
	getterRequest  *pb.SecretGetterRequest
	*DMServices
	organization *models.Organization
	result       *pb.VapusSecretsResponse
	secretStore  *models.SecretStore
}

func (x *DMServices) NewSecretsManagerAgent(ctx context.Context, opts ...SecretsManagerAgentOpts) (*SecretsManagerAgent, error) {
	var err error
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		x.logger.Error().Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}
	organization, err := x.DMStore.GetOrganization(ctx, vapusPlatformClaim[encryption.ClaimOrganizationKey], vapusPlatformClaim)
	if err != nil {
		return nil, dmerrors.DMError(apperr.ErrOrganization404, err)
	}

	agent := &SecretsManagerAgent{
		organization: organization,
		DMServices:   x,
		result:       &pb.VapusSecretsResponse{},
		VapusInterfaceBase: &processes.VapusInterfaceBase{
			InitAt: dmutils.GetEpochTime(),
			// Ctx:       ctx,
			CtxClaim:  vapusPlatformClaim,
			AgentType: types.SECRETMANAGERAGENT.String(),
		},
	}
	for _, opt := range opts {
		opt(agent)
	}
	agent.SetAgentId()
	agent.Logger = pkgs.GetSubDMLogger(types.SECRETMANAGERAGENT.String(), agent.AgentId)
	return agent, nil
}

func (x SecretsManagerAgent) GetResult() *pb.VapusSecretsResponse {
	return x.result
}

func (x SecretsManagerAgent) Act(ctx context.Context) error {
	switch x.Action {
	case mpb.ResourceLcActions_CREATE.String():
		return x.createSecret(ctx)
	case mpb.ResourceLcActions_LIST.String():
		return x.listSecret(ctx)
	case mpb.ResourceLcActions_GET.String():
		return x.getSecret(ctx)
	case mpb.ResourceLcActions_UPDATE.String():
		x.secretStore = utils.SMPb2Obj(x.managerRequest.GetSpec())
		return x.updateSecret(ctx)
	case mpb.ResourceLcActions_ARCHIVE.String():
		return x.archiveSecret(ctx)
	default:
		x.Logger.Error().Msgf("invalid action - %v", x.Action)
		return dmerrors.DMError(apperr.ErrInvalidSecretManagerAction, nil) //nolint:wrapcheck
	}
}

func (x *SecretsManagerAgent) createSecret(ctx context.Context) error {
	if x.managerRequest == nil || x.managerRequest.GetSpec() == nil {
		x.Logger.Error().Msg("error while creating secret")
		return dmerrors.DMError(apperr.ErrInvalidSecretSpec, nil)
	}
	obj := &models.SecretStore{}
	switch x.managerRequest.GetSpec().SecretType {
	case mpb.VapusSecretType_VAPUS_CREDENTIAL:
		secObj := &models.GenericCredentialModel{}
		err := json.Unmarshal(x.managerRequest.GetSpec().GetData(), secObj)
		if err != nil {
			x.Logger.Error().Err(err).Msg("error while unmarshalling secret data to generic credential model")
			return dmerrors.DMError(apperr.ErrCreateSecret, err)
		}
		err = apppkgs.SaveCredentialsCreds(ctx, x.managerRequest.GetSpec().GetName(), secObj, x.DMStore.VapusStore, x.Logger)
		if err != nil {
			x.Logger.Error().Err(err).Msg("error while saving the data in the vault")
			return dmerrors.DMError(apperr.ErrCreateSecret, err)
		}
	case mpb.VapusSecretType_CUSTOM_SECRET:
		err := json.Unmarshal(x.managerRequest.GetSpec().GetData(), &map[string]interface{}{})
		if err != nil {
			x.Logger.Error().Err(err).Msg("error while unmarshalling secret data to custom secret map, storing string as it is")
		}
		err = x.DMStore.SecretStore.WriteSecret(ctx, string(x.managerRequest.GetSpec().GetData()), x.managerRequest.GetSpec().GetName())
		if err != nil {
			x.Logger.Error().Err(err).Msg("error while writing secret")
			return dmerrors.DMError(apperr.ErrCreateSecret, err)
		}
	default:
		x.Logger.Error().Msg("error while creating secret")
		return dmerrors.DMError(apperr.ErrInvalidSecretType, nil)
	}
	obj.Name = x.managerRequest.GetSpec().GetName()
	obj.SecretType = x.managerRequest.GetSpec().GetSecretType().String()
	obj.Provider = x.DMStore.SecretStore.GetCreds().DataSourceService
	obj.Description = x.managerRequest.GetSpec().GetDescription()
	obj.ExpireAt = x.managerRequest.GetSpec().GetExpireAt()
	obj.PreSaveCreate(x.CtxClaim)
	err := x.DMStore.CreateVapusSecret(ctx, obj, x.CtxClaim)
	if err != nil {
		x.Logger.Error().Err(err).Msg("error while creating secret")
		return dmerrors.DMError(apperr.ErrCreateSecret, err)
	}
	x.SetCreateResponse(mpb.Resources_SECRETS, obj.Name)
	return nil
}

func (x *SecretsManagerAgent) getSecret(ctx context.Context) error {
	result, err := x.DMStore.GetVapusSecret(ctx, x.getterRequest.GetName(), x.CtxClaim)
	if err != nil || result == nil {
		x.Logger.Error().Err(err).Msg("error while getting secret")
		return dmerrors.DMError(apperr.ErrGetSecret, err)
	} else {
		if result.CreatedBy != x.CtxClaim[encryption.ClaimUserIdKey] {
			x.Logger.Error().Msg("error while getting secret")
			return dmerrors.DMError(apperr.ErrSecret403, nil)
		}
		x.result.Output = []*mpb.SecretStore{
			result.ConvertToPb(),
		}
	}
	return nil
}
func (x *SecretsManagerAgent) listSecret(ctx context.Context) error {
	filter := fmt.Sprintf("organization = '%s' AND created_by = '%s'", x.CtxClaim[encryption.ClaimOrganizationKey], x.CtxClaim[encryption.ClaimUserIdKey])
	result, err := x.DMStore.ListVapusSecrets(ctx, filter, x.CtxClaim)
	if err != nil || result == nil {
		x.Logger.Error().Err(err).Msg("error while listing secret")
		return dmerrors.DMError(apperr.ErrListSecrets, err)
	} else {
		for _, r := range result {
			x.result.Output = append(x.result.Output, r.ConvertToPb())
		}
	}
	return nil
}
func (x *SecretsManagerAgent) updateSecret(ctx context.Context) error {
	return nil
}
func (x *SecretsManagerAgent) archiveSecret(ctx context.Context) error {
	result, err := x.DMStore.GetVapusSecret(ctx, x.getterRequest.GetName(), x.CtxClaim)
	if err != nil || result == nil {
		x.Logger.Error().Err(err).Msg("error while getting secret")
		return dmerrors.DMError(apperr.ErrGetSecret, err)
	} else {
		if result.CreatedBy != x.CtxClaim[encryption.ClaimUserIdKey] {
			x.Logger.Error().Msg("error while getting secret")
			return dmerrors.DMError(apperr.ErrSecret403, nil)
		}
		err := x.DMStore.SecretStore.DeleteSecret(ctx, result.Name)
		if err != nil {
			x.Logger.Error().Err(err).Msg("error while deleting secret")
			return dmerrors.DMError(apperr.ErrArchiveSecret, err)
		}
		result.PreSaveDelete(x.CtxClaim)
		err = x.DMStore.DeleteVapusSecret(ctx, result.Name, x.CtxClaim)
		if err != nil {
			x.Logger.Error().Err(err).Msg("error while archiving secret")
			return dmerrors.DMError(apperr.ErrArchiveSecret, err)
		}
	}
	return nil
}
