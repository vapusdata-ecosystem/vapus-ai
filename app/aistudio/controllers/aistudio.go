package dmcontrollers

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/aistudio/pkgs"
	"github.com/vapusdata-ecosystem/vapusdata/aistudio/services"
	dmsvc "github.com/vapusdata-ecosystem/vapusdata/aistudio/services"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
)

type AIStudio struct {
	pb.UnimplementedAIStudioServer
	validator  *dmutils.DMValidator
	DMServices *dmsvc.AIStudioServices
	Logger     zerolog.Logger
}

var AIModelInterfaceManager *AIStudio

func NewAIStudio() *AIStudio {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "AIModelInterface")
	validator, err := dmutils.NewDMValidator()
	if err != nil {
		l.Panic().Err(err).Msg("Failed to initialize validator")
	}

	l.Info().Msg("AIModelInterface Controller initialized")
	return &AIStudio{
		validator:  validator,
		Logger:     l,
		DMServices: dmsvc.AIStudioServiceManager,
	}
}

func InitNewAIStudio() {
	if AIModelInterfaceManager == nil {
		AIModelInterfaceManager = NewAIStudio()
	}
}

func (v *AIStudio) GenerateEmbeddings(ctx context.Context, req *pb.EmbeddingsInterface) (*pb.EmbeddingsResponse, error) {
	agent, err := v.DMServices.NewEmbeddingAgent(ctx, req)
	if err != nil {
		v.Logger.Error().Err(err).Msg("Failed to create EmbeddingAgent")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("Failed to execute EmbeddingAgent")
		return nil, err
	}
	return &pb.EmbeddingsResponse{
		Output: agent.GetEmbeddings(),
	}, nil
}

func (v *AIStudio) Getter(ctx context.Context, req *pb.GetAIChatRequest) (*pb.AIChatResponse, error) {
	agent, err := v.DMServices.NewAIStudioChatManager(ctx, &services.AIStudioChatAgentRequest{
		GetterRequest: req,
	})

	if err != nil {
		v.Logger.Error().Err(err).Msg("Failed to get AIStudioChatManager")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("Failed to execute AIStudioChatManager getter request")
		return nil, err
	}
	return agent.GetResponse(), nil
}

func (v *AIStudio) Manager(ctx context.Context, req *pb.ManageAIChatRequest) (*pb.AIChatResponse, error) {
	agent, err := v.DMServices.NewAIStudioChatManager(ctx, &services.AIStudioChatAgentRequest{
		ManageRequest: req,
	})
	if err != nil {
		v.Logger.Error().Err(err).Msg("Failed to get AIStudioChatManager")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("Failed to execute AIStudioChatManager manager request")
		return nil, err
	}
	return agent.GetResponse(), nil
}

func (v *AIStudio) Completions(ctx context.Context, req *pb.ChatRequest) (*pb.ChatResponse, error) {
	return v.DMServices.PointChat(ctx, req)
}

func (v *AIStudio) Chat(req *pb.ChatRequest, stream pb.AIStudio_ChatServer) error {
	fmt.Println("I am in Chat Stream: ============") // sabse Phale chat request ayah pr ayega....
	fmt.Println(req.ModelNodeId)
	return v.DMServices.ChatStream(req, stream)
}

func (v *AIStudio) BidiChat(stream pb.AIStudio_BidiChatServer) error {
	return v.DMServices.ChatAgent(stream)
}
