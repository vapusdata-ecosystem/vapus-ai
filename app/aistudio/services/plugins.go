package services

import (
	"context"
	"fmt"
	"slices"
	"strings"

	guuid "github.com/google/uuid"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/aistudio/utils"
	aidmstore "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo/aistudio"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	models "github.com/vapusdata-ecosystem/vapusai/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	"github.com/vapusdata-ecosystem/vapusai/core/plugins"
	processes "github.com/vapusdata-ecosystem/vapusai/core/process"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

type PluginIntAgentOpts func(*PluginManagerAgent)

func WithPluginAgentManagerRequest(managerRequest *pb.PluginManagerRequest) PluginIntAgentOpts {
	return func(v *PluginManagerAgent) {
		v.managerRequest = managerRequest
	}
}

func WithPluginAgentGetterRequest(getterRequest *pb.PluginGetterRequest) PluginIntAgentOpts {
	return func(v *PluginManagerAgent) {
		v.getterRequest = getterRequest
	}
}

func WithPluginAgentManagerAction(action string) PluginIntAgentOpts {
	return func(v *PluginManagerAgent) {
		v.Action = action
	}
}

type PluginManagerAgent struct {
	*processes.VapusInterfaceBase
	managerRequest *pb.PluginManagerRequest
	getterRequest  *pb.PluginGetterRequest
	result         []*mpb.Plugin
	organization   *models.Organization
	dmStore        *aidmstore.AIStudioDMStore
}

func (s *AIStudioServices) NewPluginManagerAgent(ctx context.Context, opts ...PluginIntAgentOpts) (*PluginManagerAgent, error) {
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		s.Logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, encryption.ErrInvalidJWTClaims
	}

	organization, err := s.DMStore.GetOrganization(ctx, vapusPlatformClaim[encryption.ClaimOrganizationKey], vapusPlatformClaim)
	if err != nil {
		s.Logger.Error().Err(err).Ctx(ctx).Msg("error while getting organization from datastore")
		return nil, dmerrors.DMError(apperr.ErrOrganization404, err)
	}
	agent := &PluginManagerAgent{
		result:       make([]*mpb.Plugin, 0),
		dmStore:      s.DMStore,
		organization: organization,
		VapusInterfaceBase: &processes.VapusInterfaceBase{
			CtxClaim: vapusPlatformClaim,
			// Ctx:      ctx,
			InitAt: dmutils.GetEpochTime(),
		},
	}
	agent.SetAgentId()
	for _, opt := range opts {
		opt(agent)
	}
	agent.Logger = pkgs.GetSubDMLogger(types.AIPROMPTAGENT.String(), agent.AgentId)
	return agent, nil
}

func (v *PluginManagerAgent) GetAgentId() string {
	return v.AgentId
}

func (v *PluginManagerAgent) GetResult() []*mpb.Plugin {
	v.FinishAt = dmutils.GetEpochTime()
	v.FinalLog()
	return v.result
}

func (v *PluginManagerAgent) Act(ctx context.Context) error {
	switch v.GetAction() {
	case mpb.ResourceLcActions_ADD.String():
		return v.configurePlugin(ctx)
	case mpb.ResourceLcActions_UPDATE.String():
		return v.patchPlugin(ctx)
	case mpb.ResourceLcActions_ARCHIVE.String():
		return v.archivePlugins(ctx)
	case mpb.ResourceLcActions_LIST.String():
		return v.listPlugins(ctx)
	case mpb.ResourceLcActions_GET.String():
		return v.describePlugin(ctx)
	default:
		v.Logger.Error().Msg("invalid action")
		return apperr.ErrInvalidAction
	}
}

func (v *PluginManagerAgent) configurePlugin(ctx context.Context) error {
	fmt.Println("========== I am creating Plugins ==========")
	if v.managerRequest.GetSpec().GetPluginType() != mpb.IntegrationPluginTypes_GUARDRAILS {
		cQ := fmt.Sprintf("status = 'ACTIVE' AND deleted_at IS NULL AND plugin_type = '%s'", v.managerRequest.GetSpec().GetPluginType())
		switch v.managerRequest.GetSpec().GetScope() {
		case mpb.ResourceScope_ORGANIZATION_SCOPE.String():
			cQ += fmt.Sprintf(" AND scope = '%s' AND organization = '%s'", mpb.ResourceScope_USER_SCOPE.String(), v.CtxClaim[encryption.ClaimOrganizationKey])
		case mpb.ResourceScope_USER_SCOPE.String():
			cQ += fmt.Sprintf(" AND scope = '%s' AND created_by = '%s' AND organization = '%s'", mpb.ResourceScope_USER_SCOPE.String(), v.CtxClaim[encryption.ClaimUserIdKey], v.CtxClaim[encryption.ClaimOrganizationKey])
		case mpb.ResourceScope_PLATFORM_SCOPE.String():
			cQ += fmt.Sprintf(" AND scope = '%s'", mpb.ResourceScope_PLATFORM_SCOPE.String())
		default:
			return dmerrors.DMError(apperr.ErrInvalidPluginScopeParams, nil)
		}
		count, err := v.dmStore.CountPlugins(ctx, cQ, v.CtxClaim)
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while fetching the plugin creds")
			return err
		}
		if count > 0 {
			v.Logger.Error().Msg("plugin already exists")
			return dmerrors.DMError(apperr.ErrPluginServiceScopeExists, nil)
		}
		fmt.Println("Inside If block==============>>>>>>>>>>>")
	}
	fmt.Println("Outside of if block")
	validScope, ok := plugins.PluginTypeScopeMap[v.managerRequest.GetSpec().GetPluginType().String()]
	if !ok {
		v.Logger.Error().Msg("invalid plugin type")
		return dmerrors.DMError(apperr.ErrUnSupportedPluginType, nil)
	}
	if !slices.Contains(validScope, v.managerRequest.GetSpec().GetScope()) {
		v.Logger.Error().Msg("invalid plugin scope")
		return dmerrors.DMError(apperr.ErrPluginScope403, nil)
	}
	plugin := (&models.Plugin{}).ConvertFromPb(v.managerRequest.GetSpec())
	plugin.PreSaveCreate(v.CtxClaim)
	if plugin.NetworkParams.SecretName == types.EMPTYSTR {
		plugin.NetworkParams.SecretName = dmutils.GetSecretName("plugin", "", "creds-"+guuid.New().String())
		err := apppkgs.SaveCredentialsCreds(ctx, plugin.NetworkParams.SecretName, plugin.NetworkParams.Credentials, v.dmStore.VapusStore, v.Logger)
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while saving plugin creds")
			return err
		}
	}

	if plugin.Scope == mpb.ResourceScope_PLATFORM_SCOPE.String() {
		if v.organization.OrganizationType != mpb.OrganizationType_SERVICE_ORGANIZATION.String() {
			return dmerrors.DMError(apperr.ErrPluginPlatformScope403, nil)
		}
	}
	if plugin.Scope == mpb.ResourceScope_ORGANIZATION_SCOPE.String() {
		if !strings.Contains(v.CtxClaim[encryption.ClaimOrganizationRolesKey], mpb.OrgRoles_ORG_OWNER.String()) {
			return dmerrors.DMError(apperr.ErrPluginOrganizationScope403, nil)
		}
	}
	plugin.Organization = v.CtxClaim[encryption.ClaimOrganizationKey]
	plugin.NetworkParams.Credentials = nil
	plugin.NetworkParams.IsAlreadyInSecretBS = true
	err := v.dmStore.ConfigurePlugin(ctx, plugin, v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while configuring plugin")
		return dmerrors.DMError(apperr.ErrPluginCreate400, err)
	}
	v.SetCreateResponse(mpb.Resources_DATASOURCES, plugin.VapusID)
	return nil
}

func (v *PluginManagerAgent) patchPlugin(ctx context.Context) error {
	existingObj, err := v.dmStore.GetPlugin(ctx, v.managerRequest.GetSpec().PluginId, v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting plugin from datastore")
		return dmerrors.DMError(apperr.ErrPlugin404, err)
	}
	if !existingObj.Editable {
		v.Logger.Error().Msg("plugin not editable")
		return dmerrors.DMError(apperr.ErrPluginNotEditable, nil)
	} else if existingObj.Organization != v.CtxClaim[encryption.ClaimOrganizationKey] && existingObj.CreatedBy != v.CtxClaim[encryption.ClaimUserIdKey] {
		v.Logger.Error().Msg("plugin not in organization")
		return dmerrors.DMError(apperr.ErrPlugin403, nil)
	}
	plugin := (&models.Plugin{}).ConvertFromPb(v.managerRequest.GetSpec())
	existingObj.Name = plugin.Name
	existingObj.PreSaveUpdate(v.CtxClaim[encryption.ClaimUserIdKey])
	if plugin.NetworkParams.Credentials != nil {
		existingObj.NetworkParams = plugin.NetworkParams
		existingObj.NetworkParams.SecretName = dmutils.GetSecretName("plugin", "", dmutils.GetStrEpochTime())
		err = apppkgs.SaveCredentialsCreds(ctx, existingObj.NetworkParams.SecretName, existingObj.NetworkParams.Credentials, v.dmStore.VapusStore, v.Logger)
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while saving plugin creds")
			return err
		}
		existingObj.NetworkParams.Credentials = nil
		existingObj.NetworkParams.IsAlreadyInSecretBS = true
	} else if plugin.NetworkParams.SecretName != types.EMPTYSTR {
		existingObj.NetworkParams.SecretName = plugin.NetworkParams.SecretName
	}
	existingObj.Status = mpb.CommonStatus_ACTIVE.String()
	// DynamicParams can also be updated
	existingObj.DynamicParams = plugin.DynamicParams
	err = v.dmStore.PutPlugin(ctx, existingObj, v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while configuring plugin")
		return dmerrors.DMError(apperr.ErrPluginPatch400, err)
	}
	v.result = []*mpb.Plugin{existingObj.ConvertToPb()}
	return nil
}

func (v *PluginManagerAgent) describePlugin(ctx context.Context) error {
	result, err := v.dmStore.GetPlugin(ctx, v.getterRequest.GetPluginId(), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting plugin from datastore")
		return dmerrors.DMError(apperr.ErrPlugin404, err)
	}
	if !result.Editable {
		v.Logger.Error().Msg("plugin not editable")
		return dmerrors.DMError(apperr.ErrPluginNotEditable, nil)
	} else if result.Organization != v.CtxClaim[encryption.ClaimOrganizationKey] && result.CreatedBy != v.CtxClaim[encryption.ClaimUserIdKey] {
		v.Logger.Error().Msg("plugin not in organization")
		return dmerrors.DMError(apperr.ErrPlugin403, nil)
	}
	v.result = []*mpb.Plugin{result.ConvertToPb()}
	return nil
}

func (v *PluginManagerAgent) listPlugins(ctx context.Context) error {
	v.Logger.Info().Msg("getting ai agent list from datastore")
	result, err := v.dmStore.ListPlugins(ctx,
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
	v.result = utils.DPPLObj2Pb(result)
	return nil
}

func (v *PluginManagerAgent) archivePlugins(ctx context.Context) error {
	v.Logger.Info().Msg("getting ai agent list from datastore")
	result, err := v.dmStore.GetPlugin(ctx, v.getterRequest.GetPluginId(), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting plugin from datastore")
		return dmerrors.DMError(apperr.ErrPlugin404, err)
	}
	if result.Organization != v.CtxClaim[encryption.ClaimOrganizationKey] && result.CreatedBy != v.CtxClaim[encryption.ClaimUserIdKey] {
		v.Logger.Error().Msg("plugin not onwed by current user")
		return dmerrors.DMError(apperr.ErrPlugin403, nil)
	}
	result.Status = mpb.CommonStatus_DELETED.String()
	result.DeletedAt = dmutils.GetEpochTime()
	result.DeletedBy = v.CtxClaim[encryption.ClaimUserIdKey]
	err = v.dmStore.PutPlugin(ctx, result, v.CtxClaim)
	v.result = nil
	return nil
}
