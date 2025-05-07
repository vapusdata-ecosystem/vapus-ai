package processes

import (
	"context"
	fmt "fmt"

	guuid "github.com/google/uuid"
	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusdata/core/pkgs/pbtools"
)

// Non DB model
type VapusInterfaceBase struct {
	CtxClaim       map[string]string  `json:"ctx_claim" yaml:"ctx_claim"`
	Action         string             `json:"action" yaml:"action"`
	Logger         zerolog.Logger     `json:"logger" yaml:"logger"`
	AgentId        string             `json:"agent_id" yaml:"agent_id"`
	InitAt         int64              `json:"init_at" yaml:"init_at"`
	FinishAt       int64              `json:"finish_at" yaml:"finish_at"`
	Status         string             `json:"status" yaml:"status"`
	AgentType      string             `json:"agent_type" yaml:"agent_type"`
	MetaData       map[string]any     `json:"meta_data" yaml:"meta_data"`
	ContextCancel  context.CancelFunc `json:"context_cancel" yaml:"context_cancel"`
	createResponse *mpb.VapusCreateResponse
}

func (x *VapusInterfaceBase) FinalLog() {
	x.Logger.Info().Msgf("%v - %v action %v started at %v and finished at %v with status %v", x.AgentType, x.AgentId, x.Action, x.InitAt, x.FinishAt, x.Status)
	x.SetAgentLog(x.Logger.Info(), "finalLog", fmt.Sprintf("%v - %v action %v started at %v and finished at %v with status %v", x.AgentType, x.AgentId, x.Action, x.InitAt, x.FinishAt, x.Status))
}

func (x *VapusInterfaceBase) GetAgentLogs() map[string]any {
	return x.MetaData
}

func (x *VapusInterfaceBase) SetCreateResponse(resource mpb.Resources, resourceId string) {
	if x.createResponse == nil {
		x.createResponse = &mpb.VapusCreateResponse{
			Result: &mpb.VapusCreateResponse_ResourceInfo{
				Resource:   resource,
				ResourceId: resourceId,
			},
		}
	} else {
		x.createResponse.Result = &mpb.VapusCreateResponse_ResourceInfo{
			Resource:   resource,
			ResourceId: resourceId,
		}
	}
}

func (x *VapusInterfaceBase) InitCreateResponse() *mpb.VapusCreateResponse {
	x.createResponse = &mpb.VapusCreateResponse{
		Result: &mpb.VapusCreateResponse_ResourceInfo{},
	}
	return x.createResponse
}

func (x *VapusInterfaceBase) GetCreateResponse() *mpb.VapusCreateResponse {
	return x.createResponse
}

func (x *VapusInterfaceBase) SetAgentLog(loggerEvent *zerolog.Event, key, val string) {
	loggerEvent.Msgf("%s -- %v", key, val)
	if x.MetaData == nil {
		x.MetaData = make(map[string]any)
	}
	x.MetaData[key] = val
}

func (x *VapusInterfaceBase) GetAction() string {
	return x.Action
}

func (x *VapusInterfaceBase) SetAgentId() {
	if x.AgentId == "" {
		x.AgentId = guuid.New().String()
	}
}

// func (x *VapusInterfaceBase) SetNewGoRCtx(withAuth bool) error {
// 	if withAuth {
// 		nCtx, err := pbtools.SwapNewContextWithAuthToken(x.Ctx)
// 		if err != nil {
// 			x.Logger.Error().Msgf("error swapping context with auth token: %v", err)
// 			return err
// 		}
// 		x.Ctx = nCtx
// 	} else {
// 		x.Ctx = context.Background()
// 	}
// 	return nil
// }

func (x *VapusInterfaceBase) GetNewGoRCtx(ctx context.Context, withAuth bool) context.Context {
	if withAuth {
		return pbtools.SwapNewContextWithAuthToken(ctx)
	} else {
		return context.Background()
	}
}
