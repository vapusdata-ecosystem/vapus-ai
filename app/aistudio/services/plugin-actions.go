package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	dmstores "github.com/vapusdata-ecosystem/vapusai/aistudio/datastoreops"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	appdrepo "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo"
	aidmstore "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo/aistudio"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	models "github.com/vapusdata-ecosystem/vapusai/core/models"
	searchengine "github.com/vapusdata-ecosystem/vapusai/core/operator/search"
	"github.com/vapusdata-ecosystem/vapusai/core/options"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	processes "github.com/vapusdata-ecosystem/vapusai/core/process"
)

type PluginActionsAgent struct {
	*processes.VapusInterfaceBase
	request *pb.PluginActionRequest
	result  []*mpb.Plugin
	dmStore *aidmstore.AIStudioDMStore
}

func (s *AIStudioServices) NewPluginActionsAgent(ctx context.Context, request *pb.PluginActionRequest) (*PluginActionsAgent, error) {
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		s.Logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, encryption.ErrInvalidJWTClaims
	}
	agent := &PluginActionsAgent{
		request: request,
		result:  make([]*mpb.Plugin, 0),
		VapusInterfaceBase: &processes.VapusInterfaceBase{
			CtxClaim: vapusPlatformClaim,
			// Ctx:      ctx,
			InitAt: dmutils.GetEpochTime(),
		},
	}
	agent.SetAgentId()
	agent.Logger = pkgs.GetSubDMLogger(request.PluginType, agent.AgentId)
	return agent, nil
}

func (v *PluginActionsAgent) GetAgentId() string {
	return v.AgentId
}

func (v *PluginActionsAgent) GetResult() []*mpb.Plugin {
	v.FinishAt = dmutils.GetEpochTime()
	v.FinalLog()
	return v.result
}

func (v *PluginActionsAgent) Act(ctx context.Context) error {
	switch v.request.PluginType {
	case mpb.IntegrationPluginTypes_EMAIL.String():
		return v.sendEmail(ctx)
	case mpb.IntegrationPluginTypes_FILESTORES.String():
		return v.fileStoreAction(ctx)
	default:
		v.Logger.Error().Ctx(ctx).Msg("invalid plugin type for action")
		return apperr.ErrInvalidPluginTypeForAction
	}
}

func (v *PluginActionsAgent) sendEmail(ctx context.Context) error {
	if pkgs.PluginServiceManager == nil {
		appdrepo.NewPluginPool(context.Background(), dmstores.DMStoreManager.VapusStore, v.Logger)
	}
	emailer := pkgs.PluginServiceManager.PlatformPlugins.Emailer
	reqObj := &options.SendEmailRequest{}
	if err := json.Unmarshal(v.request.GetSpec(), reqObj); err != nil {
		v.Logger.Error().Err(err).Msg("error while unmarshalling request")
		return err
	}
	return emailer.SendRawEmail(ctx, reqObj, v.AgentId)
}

func (v *PluginActionsAgent) fileStoreAction(ctx context.Context) error {
	plQ := fmt.Sprintf("deleted_at is null and status = 'ACTIVE' AND plugin_type='%s' AND scope='PLATFORM_SCOPE'",
		mpb.IntegrationPluginTypes_FILESTORES.String())
	q := fmt.Sprintf("organization='%s' AND created_by='%s' AND plugin_type='%s'",
		v.CtxClaim[encryption.ClaimOrganizationKey],
		v.CtxClaim[encryption.ClaimUserIdKey],
		mpb.IntegrationPluginTypes_FILESTORES.String())
	var uPl *models.Plugin
	var plPL *models.Plugin
	var sourceCreds *models.DataSourceCreds
	var err error

	var wg sync.WaitGroup
	wg.Add(2)
	var errChan = make(chan error, 2)
	go func() {
		defer wg.Done()
		res, err := v.dmStore.ListPlugins(ctx, q, v.CtxClaim)
		if err != nil || len(res) == 0 {
			v.Logger.Error().Err(err).Msg("error while listing file store plguin for user")
			errChan <- dmerrors.DMError(apperr.ErrFileStorePlugin404, err)
			return
		}
		uPl = res[0]
	}()
	go func() {
		defer wg.Done()
		res, err := v.dmStore.ListPlugins(ctx, plQ, v.CtxClaim)
		if err != nil || len(res) == 0 {
			v.Logger.Error().Err(err).Msg("error while listing file store plguin for platform")
			errChan <- dmerrors.DMError(apperr.ErrFileStorePlugin404, err)
			return
		}
		plPL = res[0]
		sourceCreds, err = v.dmStore.GetDataCredsFromSecret(ctx, plPL.NetworkParams.SecretName)
		if err != nil || sourceCreds == nil {
			v.Logger.Error().Err(err).Msg("error while getting filestore creds from secret store")
			errChan <- dmerrors.DMError(apperr.ErrPlugin404, err)
		}
	}()
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while getting file store plugin")
			return dmerrors.DMError(apperr.ErrFileStorePlugin404, err)
		}
	}
	if uPl == nil {
		return dmerrors.DMError(apperr.ErrFileStorePlugin404, err)
	}
	sourceCreds.Credentials.Username = v.CtxClaim[encryption.ClaimUserIdKey]
	return nil
}

func (v *PluginActionsAgent) search(ctx context.Context) error {
	if pkgs.PluginServiceManager == nil {
		appdrepo.NewPluginPool(context.Background(), dmstores.DMStoreManager.VapusStore, v.Logger)
	}
	var searchPlugin searchengine.Search
	userPlugin, ok := pkgs.PluginServiceManager.UserPlugins.Load(v.CtxClaim[encryption.ClaimUserIdKey])
	if ok {
		searchPlugin = userPlugin.(searchengine.Search)
	}
	organizationPlugin, ok := pkgs.PluginServiceManager.OrganizationPlugins.Load(v.CtxClaim[encryption.ClaimOrganizationKey])
	if ok {
		searchPlugin = organizationPlugin.(searchengine.Search)
	}
	if searchPlugin == nil {
		searchPlugin = pkgs.PluginServiceManager.PlatformPlugins.SearchEngine
	}
	params := &options.SearchInput{}
	if err := json.Unmarshal(v.request.GetSpec(), params); err != nil {
		v.Logger.Error().Err(err).Msg("error while unmarshalling request")
		return err
	}
	params.Engine = ""
	// result
	// if params.SearchRaw {
	// 	return searchPlugin.SearchRaw(params)
	// }
	// return searchPlugin.SearchFormatted(params)
	return nil
}
