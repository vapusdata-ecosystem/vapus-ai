package services

import (
	"context"

	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	aiinteface "github.com/vapusdata-ecosystem/vapusai/core/aistudio/interface"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

func (v *AIStudioServices) ChatAgent(stream pb.AIStudio_BidiChatServer) error {
	ctx := stream.Context()
	agent, err := aiinteface.NewAIAgentChat(ctx, v.Logger,
		aiinteface.WithAgentChatGrpcStream(stream),
		aiinteface.WithAgentChatAiBase(&aiinteface.AIBaseInterface{
			ModelPool:      pkgs.AIModelNodeConnectionPoolManager,
			GuardrailPool:  pkgs.GuardrailPoolManager,
			PlatformAIAttr: v.DMStore.Account.GetAiAttributes(),
			Dmstore:        v.DMStore,
			VapusSvcClient: pkgs.VapusSvcInternalClientManager,
		}),
	)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while creating new AI Agent studio thread")
		return err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Chat(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while processing AI Agent studio request")
		return err
	}
	return nil
}

func (v *AIStudioServices) ChatStream(req *pb.ChatRequest, stream pb.AIStudio_ChatServer) error {
	ctx := stream.Context()
	agent, err := aiinteface.NewAIChatGateway(ctx, v.Logger,
		aiinteface.WithChatGatewayGrpcStream(stream),
		aiinteface.WithChatGatewayRequest(req),
		aiinteface.WithChatEnabled(true),
		aiinteface.WithChatGatewayAiBase(&aiinteface.AIBaseInterface{
			ModelPool:      pkgs.AIModelNodeConnectionPoolManager,
			GuardrailPool:  pkgs.GuardrailPoolManager,
			PlatformAIAttr: v.DMStore.Account.GetAiAttributes(),
			Dmstore:        v.DMStore,
			VapusSvcClient: pkgs.VapusSvcInternalClientManager,
		}),
	)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while creating new AI Agent studio thread")
		return err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()

	err = agent.Act(ctx)

	if err != nil {
		v.Logger.Error().Err(err).Msg("error while processing AI Agent studio request")
		return err
	}
	return nil
}

func (v *AIStudioServices) PointChat(ctx context.Context, req *pb.ChatRequest) (*pb.ChatResponse, error) {
	agent, err := aiinteface.NewAIGateway(ctx, v.Logger,
		aiinteface.WithGwRequest(req),
		aiinteface.WithGwAiBase(&aiinteface.AIBaseInterface{
			ModelPool:      pkgs.AIModelNodeConnectionPoolManager,
			GuardrailPool:  pkgs.GuardrailPoolManager,
			PlatformAIAttr: v.DMStore.Account.GetAiAttributes(),
			Dmstore:        v.DMStore,
			VapusSvcClient: pkgs.VapusSvcInternalClientManager,
		}),
	)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while creating new AI Agent studio thread")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while processing AI Agent studio request")
		return nil, err
	}
	return agent.GetResult(), nil
}
