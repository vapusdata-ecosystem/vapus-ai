package services

// import (
// 	"context"
// 	"log"
// 	"slices"
// 	"sync"

// 	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
// 	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
// 	"github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
// 	utils "github.com/vapusdata-ecosystem/vapusai/aistudio/utils"
// 	aitools "github.com/vapusdata-ecosystem/vapusai/core/aistudio/tools"
// 	appdrepo "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo"
// 	aidmstore "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo/aistudio"
// 	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
// 	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
// 	dmodels "github.com/vapusdata-ecosystem/vapusai/core/models"
// 	processes "github.com/vapusdata-ecosystem/vapusai/core/process"
// 	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
// 	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
// 	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
// 	"github.com/vapusdata-ecosystem/vapusai/core/types"
// 	agents "github.com/vapusdata-ecosystem/vapusai/core/vapusagents/agents"
// 	agenthelper "github.com/vapusdata-ecosystem/vapusai/core/vapusagents/config"
// )

// type VapusAgentManagerRequest struct {
// 	ManageRequest    *pb.AgentManagerRequest
// 	GetterRequest    *pb.AgentGetterRequest
// 	Response         *pb.AgentResponse
// 	DownloadRequest  *pb.AgentDownloadFileRequest
// 	DownloadResponse *pb.AgentDownloadFileResponse
// 	StateRequest     *pb.AgentStateRequest
// 	SignalRequest    *pb.AgentSignalRequest
// }

// type VapusAgentManager struct {
// 	*VapusAgentManagerRequest
// 	AgentError error
// 	DMStore    *aidmstore.AIStudioDMStore
// 	*processes.VapusInterfaceBase
// 	stream pb.AgentService_ValidateServer
// }

// type VapusAgentOpts func(*VapusAgentManager)

// func WithVapusAgentRequest(request *VapusAgentManagerRequest) VapusAgentOpts {
// 	return func(a *VapusAgentManager) {
// 		a.VapusAgentManagerRequest = request
// 	}
// }

// func WithVapusAgentAction(action string) VapusAgentOpts {
// 	return func(a *VapusAgentManager) {
// 		a.Action = action
// 	}
// }

// func WithValidationStream(stream pb.AgentService_ValidateServer) VapusAgentOpts {
// 	return func(a *VapusAgentManager) {
// 		a.stream = stream
// 	}
// }

// func (s *AIStudioServices) NewVapusAgentManager(ctx context.Context, opt ...VapusAgentOpts) (*VapusAgentManager, error) {
// 	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
// 	if !ok {
// 		s.Logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
// 		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
// 	}
// 	agent := &VapusAgentManager{
// 		VapusAgentManagerRequest: &VapusAgentManagerRequest{},
// 		DMStore:                  s.DMStore,
// 		VapusInterfaceBase: &processes.VapusInterfaceBase{
// 			CtxClaim: vapusPlatformClaim,
// 			// Ctx:      ctx,
// 			InitAt: dmutils.GetEpochTime(),
// 		},
// 	}
// 	for _, o := range opt {
// 		o(agent)
// 	}
// 	agent.Response = &pb.AgentResponse{}
// 	agent.SetAgentId()
// 	agent.Logger = pkgs.GetSubDMLogger(types.FABRICGETTERAGENT.String(), agent.AgentId)
// 	return agent, nil
// }

// func (a *VapusAgentManager) GetManagerResponse() *pb.AgentResponse {
// 	a.FinishAt = dmutils.GetEpochTime()
// 	a.FinalLog()
// 	return a.Response
// }

// func (a *VapusAgentManager) GetDownloadFileResponse() *pb.AgentDownloadFileResponse {
// 	a.FinishAt = dmutils.GetEpochTime()
// 	a.FinalLog()
// 	return a.DownloadResponse
// }

// func (a *VapusAgentManager) Act(ctx context.Context) error {
// 	switch a.Action {
// 	case mpb.ResourceLcActions_ADD.String():
// 		return a.createVapusAgent(ctx)
// 	case mpb.ResourceLcActions_GET.String():
// 		return a.getVapusAgent(ctx)
// 	case mpb.ResourceLcActions_ARCHIVE.String():
// 		return a.archiveVapusAgent(ctx)
// 	case mpb.ResourceLcActions_UPDATE.String():
// 		return a.updateVapusAgent(ctx)
// 	case mpb.ResourceLcActions_PUBLISH.String():
// 		return a.activateVapusAgent(ctx)
// 	case mpb.ResourceLcActions_VALIDATE.String():
// 		return a.testVapusAgent(ctx)
// 	case mpb.ResourceLcActions_UNPUBLISH.String():
// 		return a.stopVapusAgent(ctx)
// 	default:
// 		return dmerrors.DMError(apperr.ErrInvalidAgentServiceAction, nil)
// 	}
// }

// func (a *VapusAgentManager) getVapusAgent(ctx context.Context) error {
// 	if a.GetterRequest == nil {
// 		a.Logger.Error().Msg("getter request is empty")
// 		return dmerrors.DMError(apperr.ErrInvalidVapusAgenttRequested, nil)
// 	}
// 	if a.GetterRequest.VapusAgentId == "" {
// 		a.Logger.Info().Msg("getting all fabric agents for the user")
// 		filter := apppkgs.ListResourceWithGovernance(a.CtxClaim)
// 		agents, err := a.DMStore.ListVapusAgents(ctx, filter, a.CtxClaim)
// 		if err != nil || agents == nil {
// 			a.Logger.Error().Msg("error while getting fabric agents from datastore")
// 			return dmerrors.DMError(apperr.ErrInvalidVapusAgentRequested, err)
// 		}
// 		a.Response.Output = utils.FbALToPbAgent(agents)
// 	} else {
// 		a.Logger.Info().Msg("getting fabric agent for the user")
// 		agents, err := a.DMStore.GetVapusAgent(ctx, a.GetterRequest.VapusAgentId, a.CtxClaim)
// 		if err != nil || agents == nil {
// 			a.Logger.Error().Msg("error while getting fabric agent from datastore")
// 			return dmerrors.DMError(apperr.ErrInvalidVapusAgentRequested, err)
// 		}
// 		a.Response.Output = utils.FbALToPbAgent([]*dmodels.VapusAgents{agents})
// 	}
// 	return nil
// }

// func (a *VapusAgentManager) createVapusAgent(ctx context.Context) error {
// 	var err error
// 	agent := &dmodels.VapusAgents{}
// 	a.ManageRequest.GetSpec().AgentId = ""
// 	agent.ConvertFromPb(a.ManageRequest.GetSpec())
// 	agent.PreSaveCreate(a.CtxClaim)
// 	agent.Organization = a.CtxClaim[encryption.ClaimOrganizationKey]
// 	agent.AssetStore = a.CtxClaim[encryption.ClaimOrganizationKey]
// 	agent.CurrentVersion = dmutils.GetVersionNumber("", mpb.VersionBumpType_PATCH.String())
// 	_, err = appdrepo.GetOrCreateBucket(ctx, a.DMStore.VapusStore, true, agent.AssetStore)
// 	if err != nil {
// 		a.Logger.Error().Msg("error while getting or creating blob storage for your request")
// 		return dmerrors.DMError(dmerrors.ErrStorageInternalError, err)
// 	}
// 	for _, v := range agent.Specs {
// 		v.Version = agent.CurrentVersion
// 		v.VersionStatus = mpb.CommonStatus_ACTIVE.String()
// 	}
// 	for _, v := range agent.Specs {
// 		v.VersionStatus = mpb.CommonStatus_ACTIVE.String()
// 	}
// 	err = a.DMStore.CreateVapusAgent(ctx, agent, a.CtxClaim)
// 	if err != nil {
// 		a.Logger.Error().Msg("error while saving thread in datastore")
// 		return dmerrors.DMError(apperr.ErrCreatingVapusAgent, err)
// 	}
// 	agentInstance, err := agents.NewAgent(ctx, &crewmodels.CrewTools{
// 		DmStore: a.DMStore.VapusStore,
// 		// DpPool:      pkgs.DataProductServerPoolManager,
// 		IntClient:   pkgs.VapusSvcInternalClientManager,
// 		SqlOps:      pkgs.SqlOps,
// 		TrinoCl:     pkgs.TrinoClient,
// 		AIModel:     agent.Model,
// 		AIModelNode: agent.ModelNode,
// 	}, &aitools.ToolCaller{
// 		Logger:      a.Logger,
// 		AIModel:     agent.Model,
// 		AIModelNode: agent.ModelNode,
// 		Client:      pkgs.VapusSvcInternalClientManager,
// 		TcDbStore:   a.DMStore.VapusStore,
// 	}, a.Logger, crew.WithAgentId(agent.VapusID))
// 	if err != nil {
// 		a.Logger.Error().Msg("error while creating fabric agent instance")
// 		return dmerrors.DMError(apperr.ErrVapusAgentInitFailed, err)
// 	}
// 	defer agentInstance.Clear()
// 	err = agentInstance.MakeReady(ctx)
// 	if err != nil {
// 		a.Logger.Error().Msg("error while analyzing and acting on fabric agent")
// 		return dmerrors.DMError(apperr.ErrVapusAgentReadyFailed, err)
// 	}
// 	a.SetCreateResponse(mpb.Resources_AIAGENTS, agent.VapusID)
// 	return nil
// }

// func (a *VapusAgentManager) archiveVapusAgent(ctx context.Context) error {
// 	if a.GetterRequest == nil || a.GetterRequest.VapusAgentId == "" {
// 		a.Logger.Error().Msg("fabric agent id is empty")
// 		return dmerrors.DMError(apperr.ErrInvalidVapusAgenttRequested, nil)
// 	}
// 	agent, err := a.DMStore.GetVapusAgent(ctx, a.GetterRequest.VapusAgentId, a.CtxClaim)
// 	if err != nil || agent == nil {
// 		a.Logger.Error().Msg("error while getting fabric agent from datastore")
// 		return dmerrors.DMError(apperr.ErrInvalidVapusAgenttRequested, err)
// 	}
// 	if !agent.IsOwner(a.CtxClaim) {
// 		a.Logger.Error().Msg("error while validating fabric agent ownership")
// 		return dmerrors.DMError(apperr.ErrVapusAgentOwner403, nil)
// 	}
// 	agent.PreSaveDelete(a.CtxClaim)
// 	err = a.DMStore.PutVapusAgent(ctx, agent, a.CtxClaim)
// 	if err != nil {
// 		a.Logger.Error().Msg("error while archiving fabric agent in datastore")
// 		return dmerrors.DMError(apperr.ErrArchivingAgentFailed, err)
// 	}
// 	a.Response.Output = utils.FbALToPbAgent([]*dmodels.VapusAgents{agent})
// 	return nil
// }

// func (a *VapusAgentManager) updateVapusAgent(ctx context.Context) error {
// 	if a.ManageRequest == nil {
// 		a.Logger.Error().Msg("fabric agent id is empty")
// 		return dmerrors.DMError(apperr.ErrInvalidVapusAgenttRequested, nil)
// 	}
// 	agent, err := a.DMStore.GetVapusAgent(ctx, a.ManageRequest.VapusAgentId, a.CtxClaim)
// 	if err != nil || agent == nil {
// 		a.Logger.Error().Msg("error while getting fabric agent from datastore")
// 		return dmerrors.DMError(apperr.ErrInvalidVapusAgenttRequested, err)
// 	}
// 	newAgent := &dmodels.VapusAgents{}
// 	newAgent.ConvertFromPb(a.ManageRequest.GetSpec())
// 	agent.CurrentVersion = dmutils.GetVersionNumber(agent.CurrentVersion, a.ManageRequest.GetUpgradeType().String())
// 	agent.Specs = append(agent.Specs, newAgent.Specs[len(newAgent.Specs)-1])
// 	agent.Specs[len(agent.Specs)-1].Version = agent.CurrentVersion
// 	agent.PreSaveUpdate(a.CtxClaim[encryption.ClaimUserIdKey])
// 	err = a.DMStore.PutVapusAgent(ctx, agent, a.CtxClaim)
// 	if err != nil {
// 		a.Logger.Error().Msg("error while updating fabric agent in datastore")
// 		return dmerrors.DMError(apperr.ErrUpdatingVapusAgent, err)
// 	}
// 	ct := &crewmodels.CrewTools{
// 		DmStore:     a.DMStore.VapusStore,
// 		IntClient:   pkgs.VapusSvcInternalClientManager,
// 		SqlOps:      pkgs.SqlOps,
// 		TrinoCl:     pkgs.TrinoClient,
// 		AIModel:     agent.Model,
// 		AIModelNode: agent.ModelNode,
// 	}
// 	at := &aitools.ToolCaller{
// 		Logger:      a.Logger,
// 		AIModel:     agent.Model,
// 		AIModelNode: agent.ModelNode,
// 		Client:      pkgs.VapusSvcInternalClientManager,
// 		TcDbStore:   a.DMStore.VapusStore,
// 	}
// 	var wg sync.WaitGroup
// 	var errChan = make(chan error, len(agent.Specs))
// 	for _, v := range agent.Specs {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			if v.VersionStatus != mpb.CommonStatus_ACTIVE.String() {
// 				return
// 			}
// 			agentInstance, err := crew.NewAgent(ctx, ct, at, a.Logger, crew.WithAgentId(agent.VapusID), crew.WithAgentVersion(v.Version))
// 			if err != nil {
// 				a.Logger.Error().Msg("error while creating fabric agent instance")
// 				errChan <- dmerrors.DMError(apperr.ErrVapusAgentInitFailed, err)
// 			}
// 			defer agentInstance.Clear()
// 			err = agentInstance.MakeReady(ctx)
// 			if err != nil {
// 				a.Logger.Error().Msg("error while analyzing and acting on fabric agent")
// 				errChan <- dmerrors.DMError(apperr.ErrVapusAgentReadyFailed, err)
// 			}
// 		}()
// 	}
// 	wg.Wait()
// 	close(errChan)
// 	for err := range errChan {
// 		if err != nil {
// 			a.Logger.Error().Msg("error while updating fabric agent in datastore")
// 			return dmerrors.DMError(apperr.ErrUpdatingVapusAgent, err)
// 		}
// 	}
// 	a.Response.Output = utils.FbALToPbAgent([]*dmodels.VapusAgents{agent})
// 	return nil
// }

// func (a *VapusAgentManager) stopVapusAgent(ctx context.Context) error {
// 	if a.StateRequest == nil {
// 		a.Logger.Error().Msg("fabric agent id is empty")
// 		return dmerrors.DMError(apperr.ErrInvalidVapusAgenttRequested, nil)
// 	}
// 	agent, err := a.DMStore.GetVapusAgent(ctx, a.StateRequest.VapusAgentId, a.CtxClaim)
// 	if err != nil || agent == nil {
// 		a.Logger.Error().Msg("error while getting fabric agent from datastore")
// 		return dmerrors.DMError(apperr.ErrInvalidVapusAgenttRequested, err)
// 	}
// 	if agent.Status == mpb.CommonStatus_READY.String() && agent.IsOwner(a.CtxClaim) {
// 		agent.Status = mpb.CommonStatus_STOPPED.String()
// 	} else {
// 		a.Logger.Error().Msg("error while stopping fabric agent")
// 		return dmerrors.DMError(apperr.ErrVapusAgentOwner403, nil)
// 	}
// 	agent.PreSaveUpdate(a.CtxClaim[encryption.ClaimUserIdKey])
// 	return a.DMStore.PutVapusAgent(ctx, agent, a.CtxClaim)
// }

// func (a *VapusAgentManager) testVapusAgent(ctx context.Context) error {
// 	if a.SignalRequest == nil {
// 		a.Logger.Error().Msg("fabric agent id is empty")
// 		return dmerrors.DMError(apperr.ErrInvalidVapusAgenttRequested, nil)
// 	}
// 	agent, err := a.DMStore.GetVapusAgent(ctx, a.SignalRequest.VapusAgentId, a.CtxClaim)
// 	if err != nil || agent == nil {
// 		a.Logger.Error().Msg("error while getting fabric agent from datastore")
// 		return dmerrors.DMError(apperr.ErrInvalidVapusAgenttRequested, err)
// 	}
// 	agentInstance, err := crew.NewAgent(ctx, &crewmodels.CrewTools{
// 		DmStore:     a.DMStore.VapusStore,
// 		IntClient:   pkgs.VapusSvcInternalClientManager,
// 		SqlOps:      pkgs.SqlOps,
// 		TrinoCl:     pkgs.TrinoClient,
// 		AIModel:     agent.Model,
// 		AIModelNode: agent.ModelNode,
// 	}, &aitools.ToolCaller{
// 		Logger:      a.Logger,
// 		AIModel:     agent.Model,
// 		AIModelNode: agent.ModelNode,
// 		Client:      pkgs.VapusSvcInternalClientManager,
// 		TcDbStore:   a.DMStore.VapusStore,
// 	}, a.Logger,
// 		crew.WithAgentSpec(agent),
// 		crew.WithAgentSignal(a.SignalRequest),
// 		crew.WithAgentTestMode(true),
// 		crew.WithAgentValidationStream(a.stream))
// 	if err != nil {
// 		a.Logger.Error().Msg("error while creating fabric agent instance")
// 		return dmerrors.DMError(apperr.ErrVapusAgentInitFailed, err)
// 	}
// 	err = agentInstance.RunWithRetry(ctx, types.ClientRetryLimit, func(ctx context.Context) error {
// 		return agentInstance.Run(ctx)
// 	})
// 	if err != nil {
// 		a.Logger.Error().Msg("error while analyzing and acting on fabric agent")
// 		return dmerrors.DMError(apperr.ErrVapusAgentReadyFailed, err)
// 	}
// 	agent.Status = mpb.CommonStatus_VALIDATED.String()
// 	agent.PreSaveUpdate(a.CtxClaim[encryption.ClaimUserIdKey])
// 	return a.DMStore.PutVapusAgent(ctx, agent, a.CtxClaim)
// }

// func (a *VapusAgentManager) activateVapusAgent(ctx context.Context) error {
// 	if a.StateRequest == nil {
// 		a.Logger.Error().Msg("fabric agent request is empty")
// 		return dmerrors.DMError(apperr.ErrInvalidVapusAgenttRequested, nil)
// 	}
// 	agent, err := a.DMStore.GetVapusAgent(ctx, a.StateRequest.VapusAgentId, a.CtxClaim)
// 	if err != nil || agent == nil {
// 		a.Logger.Error().Msg("error while getting fabric agent from datastore")
// 		return dmerrors.DMError(apperr.ErrInvalidVapusAgenttRequested, err)
// 	}
// 	if slices.Contains([]string{mpb.CommonStatus_VALIDATED.String(), mpb.CommonStatus_ACTIVE_READY.String()}, agent.Status) && agent.IsOwner(a.CtxClaim) {
// 		agent.Status = mpb.CommonStatus_READY.String()
// 	} else {
// 		a.Logger.Error().Msg("error while stopping fabric agent")
// 		return dmerrors.DMError(apperr.ErrVapusAgentOwner403, nil)
// 	}
// 	agent.PreSaveUpdate(a.CtxClaim[encryption.ClaimUserIdKey])
// 	return a.DMStore.PutVapusAgent(ctx, agent, a.CtxClaim)
// }
