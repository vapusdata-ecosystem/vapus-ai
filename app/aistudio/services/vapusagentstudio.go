package services

// import (
// 	"context"
// 	"log"
// 	"path/filepath"
// 	"strings"

// 	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
// 	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
// 	"github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
// 	aitools "github.com/vapusdata-ecosystem/vapusai/core/aistudio/tools"
// 	aidmstore "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo/aistudio"
// 	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
// 	processes "github.com/vapusdata-ecosystem/vapusai/core/process"
// 	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
// 	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
// 	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
// 	"github.com/vapusdata-ecosystem/vapusai/core/types"
// 	crew "github.com/vapusdata-ecosystem/vapusai/core/vapusagents/agentcrew"
// 	crewmodels "github.com/vapusdata-ecosystem/vapusai/core/vapusagents/models"
// )

// type VapusAgentStudioRequest struct {
// 	DownloadRequest  *pb.AgentDownloadFileRequest
// 	DownloadResponse *pb.AgentDownloadFileResponse
// 	Stream           pb.AgentStudio_SignalServer
// 	SignalRequest    *pb.AgentSignalRequest
// }

// type VapusAgentStudio struct {
// 	*VapusAgentStudioRequest
// 	AgentError error
// 	DMStore    *aidmstore.AIStudioDMStore
// 	*processes.VapusInterfaceBase
// }

// type VapusAgentStudioOpts func(*VapusAgentStudio)

// func WithVapusAgentSignal(request *VapusAgentStudioRequest) VapusAgentStudioOpts {
// 	return func(a *VapusAgentStudio) {
// 		a.VapusAgentStudioRequest = request
// 	}
// }

// func WithVapusAgentStudioAction(action string) VapusAgentStudioOpts {
// 	return func(a *VapusAgentStudio) {
// 		a.Action = action
// 	}
// }

// func (s *AIStudioServices) NewVapusAgentStudio(ctx context.Context, opt ...VapusAgentStudioOpts) (*VapusAgentStudio, error) {
// 	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
// 	if !ok {
// 		s.Logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
// 		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
// 	}
// 	agent := &VapusAgentStudio{
// 		VapusAgentStudioRequest: &VapusAgentStudioRequest{},
// 		DMStore:                 s.DMStore,
// 		VapusInterfaceBase: &processes.VapusInterfaceBase{
// 			CtxClaim: vapusPlatformClaim,
// 			// Ctx:      ctx,
// 			InitAt: dmutils.GetEpochTime(),
// 		},
// 	}
// 	for _, o := range opt {
// 		o(agent)
// 	}
// 	agent.SetAgentId()
// 	agent.Logger = pkgs.GetSubDMLogger(types.FABRICGETTERAGENT.String(), agent.AgentId)
// 	return agent, nil
// }

// func (a *VapusAgentStudio) GetDownloadFileResponse() *pb.AgentDownloadFileResponse {
// 	a.FinishAt = dmutils.GetEpochTime()
// 	a.FinalLog()
// 	return a.DownloadResponse
// }

// func (a *VapusAgentStudio) Act(ctx context.Context) error {
// 	switch a.Action {
// 	case mpb.ResourceLcActions_DOWNLOAD.String():
// 		return a.downloadFiles(ctx)
// 	default:
// 		return a.runVapusAgent(ctx)
// 	}
// }

// func (a *VapusAgentStudio) downloadFiles(ctx context.Context) error {
// 	if a.DownloadRequest.GetFileName() == "" {
// 		a.Logger.Error().Msg("file name is empty")
// 		return apperr.ErrInvalidFabricFileRequested
// 	}
// 	fNameL := strings.Split(a.DownloadRequest.GetFileName(), "/")
// 	if len(fNameL) < 2 {
// 		a.Logger.Error().Msg("invalid file name")
// 		return apperr.ErrInvalidFabricFileRequested
// 	}
// 	if fNameL[0] != types.VapusAgentFileKey {
// 		a.Logger.Error().Msg("invalid file name")
// 		return apperr.ErrInvalidFabricFileRequested
// 	}
// 	agent, err := a.DMStore.GetVapusAgent(ctx, fNameL[1], a.CtxClaim)
// 	if err != nil || agent == nil {
// 		a.Logger.Error().Msg("error while getting fabric agent from datastore")
// 		return dmerrors.DMError(apperr.ErrInvalidVapusAgent403, err)
// 	}
// 	bytesData, err := a.DMStore.BlobStore.DownloadObject(ctx, &processes.BlobOpsParams{
// 		BucketName: a.CtxClaim[encryption.ClaimOrganizationKey],
// 		ObjectName: a.DownloadRequest.GetFileName(),
// 	})
// 	if err != nil {
// 		a.Logger.Error().Msg("error while downloading file")
// 		return dmerrors.DMError(apperr.ErrInvalidVapusAgent404, err)
// 	}
// 	ext := strings.Replace(filepath.Ext(fNameL[len(fNameL)-1]), ".", "", -1)
// 	a.DownloadResponse = &pb.AgentDownloadFileResponse{
// 		Output: []*mpb.FileData{
// 			{
// 				Name:   fNameL[len(fNameL)-1],
// 				Format: mpb.ContentFormats(mpb.ContentFormats_value[strings.ToUpper(ext)]),
// 				Data:   bytesData,
// 				Eof:    true,
// 			},
// 		},
// 	}
// 	return nil

// }

// func (a *VapusAgentStudio) runVapusAgent(ctx context.Context) error {
// 	if a.SignalRequest != nil {
// 		a.Logger.Error().Msg("fabric agent id is empty")
// 		return dmerrors.DMError(apperr.ErrInvalidVapusAgenttRequested, nil)
// 	}
// 	agent, err := a.DMStore.GetVapusAgent(ctx, a.SignalRequest.VapusAgentId, a.CtxClaim)
// 	if err != nil || agent == nil {
// 		a.Logger.Error().Msg("error while getting fabric agent from datastore")
// 		return dmerrors.DMError(apperr.ErrInvalidVapusAgenttRequested, err)
// 	}
// 	_, err = agent.GetVersionSpec(a.SignalRequest.Version)
// 	if err != nil {
// 		a.Logger.Error().Msg("error while getting version spec")
// 		return dmerrors.DMError(apperr.ErrVapusAgentVersionSpec, err)
// 	}
// 	agentInstance, err := crew.NewAgent(ctx, &crewmodels.CrewTools{
// 		DmStore: a.DMStore.VapusStore,
// 		// DpPool:      pkgs.DataProductServerPoolManager,
// 		IntClient: pkgs.VapusSvcInternalClientManager,
// 		// SqlOps:      pkgs.SqlOps,
// 		// TrinoCl:     pkgs.TrinoClient,
// 		AIModel:     agent.Model,
// 		AIModelNode: agent.ModelNode,
// 	}, &aitools.ToolCaller{
// 		Logger:      a.Logger,
// 		AIModel:     agent.Model,
// 		AIModelNode: agent.ModelNode,
// 		Client:      pkgs.VapusSvcInternalClientManager,
// 		TcDbStore:   a.DMStore.VapusStore,
// 	}, a.Logger, crew.WithAgentSpec(agent), crew.WithAgentSignal(a.SignalRequest), crew.WithAgentRunnerStream(a.Stream))
// 	if err != nil {
// 		a.Logger.Error().Msg("error while creating fabric agent instance")
// 		return dmerrors.DMError(apperr.ErrVapusAgentInitFailed, err)
// 	}
// 	defer agentInstance.Clear()
// 	err = agentInstance.RunWithRetry(ctx, types.ClientRetryLimit, func(ctx context.Context) error {
// 		return agentInstance.Run(ctx)
// 	})
// 	if err != nil {
// 		a.Logger.Error().Msg("error while analyzing and acting on fabric agent")
// 		return dmerrors.DMError(apperr.ErrVapusAgentReadyFailed, err)
// 	}
// 	return nil
// }
