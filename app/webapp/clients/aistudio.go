package clients

import (
	"github.com/labstack/echo/v4"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
)

func (x *GrpcClient) AIModelNodes(eCtx echo.Context) []*mpb.AIModelNode {
	result, err := x.AIModelClient.List(x.SetAuthCtx(eCtx), &pb.AIModelNodeGetterRequest{})
	if err != nil {
		x.logger.Err(err).Msg("error while getting AI Model Nodes")
		return []*mpb.AIModelNode{}
	}
	return result.Output.AiModelNodes
}

func (x *GrpcClient) AIModelNodesDetails(eCtx echo.Context, id string) *mpb.AIModelNode {
	result, err := x.AIModelClient.Get(x.SetAuthCtx(eCtx), &pb.AIModelNodeGetterRequest{
		AiModelNodeId: id,
	})
	if err != nil || result == nil || result.Output == nil || len(result.Output.AiModelNodes) == 0 {
		x.logger.Err(err).Msg("error while getting AI Model Nodes")
		return nil
	}
	return result.Output.AiModelNodes[0]
}

func (x *GrpcClient) AIModelPrompts(eCtx echo.Context) []*mpb.AIPrompt {
	result, err := x.AIPromptClient.List(x.SetAuthCtx(eCtx), &pb.PromptGetterRequest{})
	if err != nil {
		x.logger.Err(err).Msg("error while getting AI Prompt List")
		return []*mpb.AIPrompt{}
	}
	return result.Output
}

func (x *GrpcClient) AIModelPromptDetails(eCtx echo.Context, id string) *mpb.AIPrompt {
	result, err := x.AIPromptClient.Get(x.SetAuthCtx(eCtx), &pb.PromptGetterRequest{
		PromptId: id,
	})
	if err != nil || len(result.Output) == 0 {
		x.logger.Err(err).Msg("error while getting AI Prompt details")
		return nil
	}
	return result.Output[0]
}

func (x *GrpcClient) ListAIGuardrails(eCtx echo.Context) []*mpb.AIGuardrails {
	result, err := x.AIGurdrailsClient.List(x.SetAuthCtx(eCtx), &pb.GuardrailsGetterRequest{})
	if err != nil {
		x.logger.Err(err).Msg("error while getting AI Prompt List")
		return []*mpb.AIGuardrails{}
	}
	return result.Output
}

func (x *GrpcClient) DescribeAIGuardrail(eCtx echo.Context, id string) *mpb.AIGuardrails {
	result, err := x.AIGurdrailsClient.Get(x.SetAuthCtx(eCtx), &pb.GuardrailsGetterRequest{
		GuardrailId: id,
	})
	if err != nil || len(result.Output) == 0 {
		x.logger.Err(err).Msg("error while getting AI Prompt details")
		return nil
	}
	return result.Output[0]
}

func (x *GrpcClient) ListAIGatewayChats(eCtx echo.Context) []*pb.AIStudioChat {
	result, err := x.AIStudioConn.Getter(x.SetAuthCtx(eCtx), &pb.GetAIChatRequest{})
	if err != nil {
		x.logger.Err(err).Msg("error while getting AIStudioChat Chat history")
		return []*pb.AIStudioChat{}
	}
	return result.Output
}

func (x *GrpcClient) GetAIGatewayChat(eCtx echo.Context, id string) ([]*pb.AIStudioChat, error) {
	result, err := x.AIStudioConn.Getter(x.SetAuthCtx(eCtx), &pb.GetAIChatRequest{
		ChatId: id,
	})
	if err != nil {
		x.logger.Err(err).Msg("error while getting AIStudioChat Chat ")
		return []*pb.AIStudioChat{}, err
	}
	return result.Output, nil
}

func (x *GrpcClient) CreateAIGatewayChat(eCtx echo.Context) ([]*pb.AIStudioChat, error) {
	result, err := x.AIStudioConn.Manager(x.SetAuthCtx(eCtx), &pb.ManageAIChatRequest{
		Action: pb.AIChatAction_CREATE,
	})
	if err != nil {
		x.logger.Err(err).Msg("error while creating AIStudioChat Chat ")
		return []*pb.AIStudioChat{}, err
	}
	return result.Output, nil
}

func (x *GrpcClient) ListVapusAgents(eCtx echo.Context) []*mpb.VapusAgent {
	result, err := x.AgentServiceClient.List(x.SetAuthCtx(eCtx), &pb.AgentGetterRequest{})
	if err != nil {
		x.logger.Err(err).Msg("error while getting AI Prompt List")
		return []*mpb.VapusAgent{}
	}
	return result.Output
}

func (x *GrpcClient) VapusAgentDetail(eCtx echo.Context, id string) *mpb.VapusAgent {
	result, err := x.AgentServiceClient.Get(x.SetAuthCtx(eCtx), &pb.AgentGetterRequest{
		VapusAgentId: id,
	})
	if err != nil || len(result.Output) == 0 {
		x.logger.Err(err).Msg("error while getting AI Prompt details")
		return nil
	}
	return result.Output[0]
}

func (x *GrpcClient) InsightsList(eCtx echo.Context, aiModelNodes []string, model []string) []*mpb.ModelNodeObservability {
	result, err := x.AIModelClient.ListInsights(x.SetAuthCtx(eCtx), &pb.AIModelNodeInsightsRequest{
		AiModelNodeId: aiModelNodes,
		Model:         model,
	})
	if err != nil {
		x.logger.Err(err).Msg("error while getting AI Model Nodes")
		return []*mpb.ModelNodeObservability{}
	}
	return result.ModelNodeObservability
}
