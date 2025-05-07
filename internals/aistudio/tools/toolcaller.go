package tools

import (
	"context"
	"fmt"
	"log"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	aipb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	aicore "github.com/vapusdata-ecosystem/vapusai/core/aistudio/core"
	aiinteface "github.com/vapusdata-ecosystem/vapusai/core/aistudio/interface"
	appdrepo "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo"
	aidmstore "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo/aistudio"
	appcl "github.com/vapusdata-ecosystem/vapusai/core/app/grpcclients"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

func NewToolCaller(
	AIModel string,
	AIModelNode string,
	Client *appcl.VapusSvcInternalClients,
	Logger zerolog.Logger,
	TcDbStore *apppkgs.VapusStore,
	ModelConnPool *appdrepo.AIModelNodeConnectionPool,
	GuardrailPool *appdrepo.GuardrailPool,
	AIDmStore *aidmstore.AIStudioDMStore,
) *ToolCaller {
	return &ToolCaller{
		AIModel:       AIModel,
		AIModelNode:   AIModelNode,
		Client:        Client,
		Logger:        Logger,
		TcDbStore:     TcDbStore,
		ModelConnPool: ModelConnPool,
		GuardrailPool: GuardrailPool,
		AIDmStore:     AIDmStore,
	}
}

func (x *ToolCaller) BuildPayload(options *Options) *aipb.ChatRequest {
	req := &aipb.ChatRequest{
		Model:       x.AIModel,
		ModelNodeId: x.AIModelNode,
	}
	if options.Temperature > 0.0 {
		req.Temperature = options.Temperature
	} else {
		req.Temperature = 0.1
	}
	if options.MaxCompletionTokens > 0 {
		req.MaxCompletionTokens = options.MaxCompletionTokens
	}
	req.Messages = []*aipb.ChatMessageObject{
		{
			Role:    aicore.USER,
			Content: options.Input,
		}, {
			Role:    aicore.SYSTEM,
			Content: options.Sysmess,
		},
	}
	if len(options.AssistantMessages) > 0 {
		for _, am := range options.AssistantMessages {
			req.Messages = append(req.Messages, &aipb.ChatMessageObject{
				Role:    aicore.ASSISTANT,
				Content: am,
			})
		}
	}
	if len(options.Tools) > 0 {
		req.Tools = options.Tools
		if options.ToolChoice != nil {
			req.ToolChoice = options.ToolChoice
		} else {
			toolName := ""
			if len(options.Tools) > 0 {
				if options.Tools[0].FunctionSchema != nil {
					toolName = options.Tools[0].FunctionSchema.Name
				}
			}
			if toolName != "" {
				req.ToolChoice = &mpb.ToolChoice{
					Type: aicore.FUNCTION,
					Function: &mpb.FunctionChoice{
						Name: toolName,
					},
				}
			}
		}
		if options.ToolChoiceParam != "" {
			req.ToolChoiceParams = options.ToolChoiceParam
		}
	}
	return req
}

func (x *ToolCaller) Tooler(ctx context.Context, param *Options) (string, error) {
	req := x.BuildPayload(param)
	tCallResp, err := x.Client.Chat(ctx, req, x.Logger, types.ClientRetryStart)
	if err != nil || len(tCallResp.GetChoices()) < 1 {
		x.Logger.Error().Msg("error while getting tool call response, or no choices found")
		return "", err
	}
	if len(tCallResp.GetChoices()[0].Messages.ToolCalls) < 1 {
		x.Logger.Error().Msg("error while getting tool call response")
		return "", err
	}

	if tCallResp.GetChoices()[0].Messages.ToolCalls[0].FunctionSchema != nil {
		return tCallResp.GetChoices()[0].Messages.ToolCalls[0].FunctionSchema.Parameters, nil
	}
	x.Logger.Error().Msg("error while getting tool call response, function schema not found")
	go func() {
		x.SaveAIToolCallLog(ctx, &models.AIToolCallLog{}, param.CtxClaim)
	}()
	return "", err
}

func (x *ToolCaller) LocalTooler(ctx context.Context, param *Options) (string, error) {
	req := x.BuildPayload(param)
	agent, err := aiinteface.NewAIGateway(ctx, x.Logger,
		aiinteface.WithGwRequest(req),
		// aiinteface.WithChatEnabled(false),
		aiinteface.WithGwAiBase(&aiinteface.AIBaseInterface{
			ModelPool:      x.ModelConnPool,
			GuardrailPool:  x.GuardrailPool,
			PlatformAIAttr: x.AIDmStore.Account.GetAiAttributes(),
			Dmstore:        x.AIDmStore,
		}),
	)
	if err != nil {
		x.Logger.Error().Err(err).Msg("error while creating new AI Agent studio thread")
		return "", err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		x.Logger.Error().Err(err).Msg("error while processing AI Agent studio request")
		return "", err
	}
	tCallResp := agent.GetResult()
	if err != nil || len(tCallResp.GetChoices()) < 1 {
		x.Logger.Error().Msg("error while getting tool call response, or no choices found")
		return "", err
	}
	var toolMess *aipb.ChatResponseChoice
	for _, choice := range tCallResp.GetChoices() {
		if choice.Messages.ToolCalls != nil && len(choice.Messages.ToolCalls) > 0 {
			toolMess = choice
			break
		}
	}
	if toolMess == nil {
		x.Logger.Error().Msg("error while getting tool call response, or no choices found")
		return "", err
	}

	if toolMess.Messages.ToolCalls[0].FunctionSchema != nil {
		log.Println("Tooler called ==========>>>>>>", toolMess.Messages.ToolCalls[0].FunctionSchema.Parameters)
		return toolMess.Messages.ToolCalls[0].FunctionSchema.Parameters, nil
	}
	x.Logger.Error().Msg("error while getting tool call response, function schema not found")
	go func() {
		x.SaveAIToolCallLog(ctx, &models.AIToolCallLog{}, param.CtxClaim)
	}()
	return "", err
}

func (x *ToolCaller) WithRetry(ctx context.Context, param *Options, toolFunc func(ctx context.Context, param *Options) (string, error)) (string, error) {
	var args string
	var err error
	param.RetryCount = 0
	if param.RetryCount == 0 {
		return toolFunc(ctx, param)
	} else {
		c := 0
		for {
			args, err = toolFunc(ctx, param)
			if err != nil || args == "" {
				x.Logger.Error().Msg("error while calling tool, retrying")
				if c > param.RetryCount {
					err = fmt.Errorf("error while calling tool, retry count exceeded")
					break
				} else {
					c++
					continue
				}
			} else {
				break
			}
		}
	}
	return args, nil
}

func (x *ToolCaller) Analyze(ctx context.Context, param *Options) (string, error) {
	return x.WithRetry(ctx, param, func(ctx context.Context, param *Options) (string, error) {
		return x.Tooler(ctx, param)
	})
}

func (x *ToolCaller) SaveAIToolCallLog(ctx context.Context, obj *models.AIToolCallLog, ctxClaim map[string]string) error {
	obj.PreSaveCreate(ctxClaim)
	_, err := x.TcDbStore.Db.PostgresClient.DB.NewInsert().
		Model(obj).
		ModelTableExpr(apppkgs.AIToolCallLogTable).Returning("id").Exec(ctx)
	if err != nil {
		x.Logger.Err(err).Ctx(ctx).Msg("error while saving fabric chat tool log in datastore")
		return err
	}
	log.Println("======================1111111", obj.Input, obj.PlainInput, obj.ID)
	upQ := fmt.Sprintf("UPDATE %s SET input = to_tsvector('english', plain_input) WHERE id = %d", apppkgs.AIToolCallLogTable, obj.ID)
	_, err = x.TcDbStore.Db.PostgresClient.DB.NewRaw(upQ).Exec(ctx)
	if err != nil {
		x.Logger.Err(err).Ctx(ctx).Msg("error while saving data source symantics and FTS vectors in datastore")
		return err
	}
	return nil
}

func (x *ToolCaller) GetAIToolCallLog(ctx context.Context, input string, actionAnalyzer, paramAnalyzer bool, ctxClaim map[string]string) (*models.AIToolCallLog, error) {
	type localVal struct {
		rank          float64 `bun:"rank"`
		output_schema string  `bun:"output_schema"`
	}
	result := []map[string]any{}
	fCon := fmt.Sprintf("param_analyzer=%t AND action_analyzer=%t", paramAnalyzer, actionAnalyzer)
	query := fmt.Sprintf(`
	SELECT rank,output_schema
	FROM (
		SELECT output_schema,ts_rank(input, plainto_tsquery('english', '%v')) AS rank
		FROM %s
		WHERE %s
	) subquery
	WHERE rank > 1.0
	ORDER BY rank DESC
	LIMIT 1;
	`,
		input, apppkgs.AIToolCallLogTable, fCon)
	log.Println("Query:>>>>>>>>>>>>><<<<<<<<<<<<<<<<<< ", query)
	err := x.TcDbStore.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		x.Logger.Err(err).Ctx(ctx).Msg("error while getting fabric chat tool log from datastore")
		return nil, err
	}
	return &models.AIToolCallLog{
		OutputSchema: result[0]["output_schema"].(string),
	}, err
}
