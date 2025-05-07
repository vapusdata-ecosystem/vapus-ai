package aiface

import (
	"context"
	"encoding/json"

	"github.com/pgvector/pgvector-go"
	"github.com/rs/zerolog"
	aiutilitypb "github.com/vapusdata-ecosystem/apis/protos/vapus-aiutilities/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusdata/core/aistudio/prompts"
	appdrepo "github.com/vapusdata-ecosystem/vapusdata/core/app/datarepo"
	aidmstore "github.com/vapusdata-ecosystem/vapusdata/core/app/datarepo/aistudio"
	appcl "github.com/vapusdata-ecosystem/vapusdata/core/app/grpcclients"
	models "github.com/vapusdata-ecosystem/vapusdata/core/models"
	"github.com/vapusdata-ecosystem/vapusdata/core/pkgs/pbtools"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
)

type AIBaseInterface struct {
	ModelPool      *appdrepo.AIModelNodeConnectionPool
	GuardrailPool  *appdrepo.GuardrailPool
	PlatformAIAttr *models.AccountAIAttributes
	logger         zerolog.Logger
	Dmstore        *aidmstore.AIStudioDMStore
	VapusSvcClient *appcl.VapusSvcInternalClients
}

// func GenerateSummary(VapusSvcClient *appcl.VapusSvcInternalClients, ctx context.Context,payload *prompts.AIEmbeddingPayload,ctxClaim map[string]string) ([]string , error) {

// 	fmt.Println("+++++++ entered Generate SUmmary")
// 	req := &aiutilitypb.SummarizerRequest{
// 		Text: []byte(payload.Input),
// 	}
// 	resp,err := VapusSvcClient.AIUtilityServerClient.Summarizer(ctx,req)

// 	if (err!= nil ){
// 		fmt.Println(err)
// 		return nil,err
// 	}

// 	decodedData := []string{}

// 	for _,data := range resp.Data {
// 		decodedData = append(decodedData, string(data))
// 	}

// 	return decodedData,nil
// }

func (s *AIBaseInterface) logRequest(ctx context.Context, parsedIp, parsedOp string, obj *models.AIStudioLog, usagelog *models.AIStudioUsages, ctxClaim map[string]string) error {
	obj.UpdatedAt = dmutils.GetEpochTime()
	embeddingModel, err := s.ModelPool.GetorSetNodeObject(s.PlatformAIAttr.EmbeddingModelNode, nil, false)
	if err != nil {
		s.logger.Error().Err(err).Msg("error while getting embedding model node")
		return err
	}
	modelConn, err := s.ModelPool.GetorSetConnection(embeddingModel, true, false)
	ip := map[string]string{
		"User Request":    parsedIp,
		"Studio Response": parsedOp,
	}
	obj.Output = append(obj.Output, obj.ParsedOutput...)
	obj.ParsedOutput = make([]*models.MessageLog, 0)
	bbytes, err := json.Marshal(ip)
	if err != nil {
		s.logger.Error().Err(err).Msg("error while marshalling studio payload log")
		bbytes = []byte(parsedIp)
	}
	payload := &prompts.AIEmbeddingPayload{
		Dimensions:     1536,
		EmbeddingModel: s.PlatformAIAttr.EmbeddingModel,
		Input:          string(bbytes),
	}
	err = modelConn.GenerateEmbeddings(ctx, payload)
	if err != nil {
		s.logger.Error().Err(err).Msgf("error while generating studio payload log embeddings from model %v", s.PlatformAIAttr.EmbeddingModel)
		return err
	}
	obj.LogEmbeddings = pgvector.NewVector(payload.Embeddings.Vectors32)

	req := &aiutilitypb.SummarizerRequest{
		Text:      bbytes,
		Sentences: 4,
	}
	resp, err := s.VapusSvcClient.GenerateSummary(ctx, req, s.logger, 3)
	decodedData := []string{}
	if err != nil {
		s.logger.Error().Err(err).Msg("error while generating summary")
	} else {
		for _, data := range resp.Data {
			decodedData = append(decodedData, string(data))
		}
		obj.Summary = decodedData
	}

	return s.Dmstore.SaveAIInterfaceLog(ctx, obj, usagelog, ctxClaim)
}

func (s *AIBaseInterface) logGuardrailRequest(nctx context.Context, parsedIp string, obj *models.AIGuardrailsLog, usage *models.AIStudioUsages, ctxClaim map[string]string) error {
	ctx, ContextCancel := pbtools.NewInCancelCtxWithAuthToken(nctx)
	go func() {
		defer ContextCancel()
		obj.UpdatedAt = dmutils.GetEpochTime()
		embeddingModel, err := s.ModelPool.GetorSetNodeObject(s.PlatformAIAttr.EmbeddingModelNode, nil, false)
		if err != nil {
			s.logger.Error().Err(err).Msg("error while getting embedding model node")
			return
		}
		modelConn, err := s.ModelPool.GetorSetConnection(embeddingModel, true, false)
		payload := &prompts.AIEmbeddingPayload{
			Dimensions:     1536,
			EmbeddingModel: s.PlatformAIAttr.EmbeddingModel,
			Input:          parsedIp,
		}
		err = modelConn.GenerateEmbeddings(ctx, payload)
		if err != nil {
			s.logger.Error().Err(err).Msgf("error while generating studio payload log embeddings from model %v", s.PlatformAIAttr.EmbeddingModel)
			return
		}
		obj.InputEmbeddings = pgvector.NewVector(payload.Embeddings.Vectors32)
		_ = s.Dmstore.SaveAIGuardrailLog(ctx, obj, usage, ctxClaim)
	}()
	return nil
}
