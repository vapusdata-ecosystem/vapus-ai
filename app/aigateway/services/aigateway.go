package services

import (
	"bufio"
	"fmt"
	"log"
	"strconv"

	json "github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusdata/aigateway/pkgs"
	aicore "github.com/vapusdata-ecosystem/vapusdata/core/aistudio/core"
	aiinteface "github.com/vapusdata-ecosystem/vapusdata/core/aistudio/interface"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
)

func (x *AIGatewayServices) ChatCompletion(c *fiber.Ctx) error {
	x.logger.Info().Msg("Chat service called")
	payload := &pb.ChatRequest{}
	nodeKey := c.Request().Header.Peek(types.AIMODEL_HEADER_KEY)
	if len(nodeKey) == 0 {
		x.logger.Error().Msg("Node key not found in the request")
		return SendAIGatewayError(c, fiber.StatusBadRequest, &aicore.AiGatewayError{
			Error: aicore.AiGatewayErrorDetail{
				Code:    strconv.Itoa(fiber.StatusBadRequest),
				Message: "Model Node key not found in the request",
				Param:   types.AIMODEL_HEADER_KEY,
			},
		})
	}
	x.logger.Info().Msgf("Node key: %s", nodeKey)

	// if err := x.unProtojson.Unmarshal(c.Request().Body(), requestObj); err != nil {
	// 	x.logger.Error().Err(err).Msg("Error while parsing the request")
	// 	return SendAIGatewayError(c, fiber.StatusBadRequest, &aicore.AiGatewayError{
	// 		Error: aicore.AiGatewayErrorDetail{
	// 			Code:    strconv.Itoa(fiber.StatusBadRequest),
	// 			Message: err.Error(),
	// 		},
	// 	})
	// }
	log.Println(string(c.Request().Body()))
	reqobj := &aicore.AIGatewayChatRequest{}
	if err := dmutils.Unmarshall(c.Request().Body(), reqobj); err != nil {
		x.logger.Error().Err(err).Msg("Error while parsing the request")
		return SendAIGatewayError(c, fiber.StatusBadRequest, &aicore.AiGatewayError{
			Error: aicore.AiGatewayErrorDetail{
				Code:    strconv.Itoa(fiber.StatusBadRequest),
				Message: err.Error(),
			},
		})
	}
	log.Println("reqobj", reqobj.Messages[0])
	payload = reqobj.ConvertToPb()
	log.Println("requestObj", payload.Messages[0], "|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||")
	payload.ModelNodeId = string(nodeKey)
	x.logger.Info().Msg("Request parsed successfully")
	x.logger.Info().Msgf("Request: %v", payload)
	validationErrors := x.Validate(payload)
	if len(validationErrors) > 0 {
		x.logger.Error().Msg("Validation failed")
		resMap := fiber.Map{}
		for _, err := range validationErrors {
			resMap[err.FailedField] = err.Tag
		}
		return c.Status(fiber.StatusBadRequest).JSON(resMap)
	}
	x.logger.Info().Msg("Validation passed")
	processOpts := []aiinteface.GwOpts{
		aiinteface.WithGwRequest(payload),
		aiinteface.WithGwAiBase(&aiinteface.AIBaseInterface{
			ModelPool:      pkgs.AIModelNodeConnectionPoolManager,
			GuardrailPool:  pkgs.GuardrailPoolManager,
			PlatformAIAttr: x.dmstores.Account.GetAiAttributes(),
			Dmstore:        x.dmstores,
			VapusSvcClient: pkgs.VapusSvcInternalClientManager,
		}),
	}
	switch payload.Stream {
	case true:
		if payload.StreamOptions == nil {
			payload.StreamOptions = &pb.StreamOptions{
				IncludeUsage: true,
			}
		}
		return x.p2stream(c, processOpts...)
	default:
		return x.p2p(c, processOpts...)
	}
}

func (x *AIGatewayServices) p2p(c *fiber.Ctx, opts ...aiinteface.GwOpts) error {
	fmt.Println("I am inside P2P========?>>>>>")
	agent, err := aiinteface.NewAIGateway(c.UserContext(), x.logger, opts...)
	if err != nil {
		x.logger.Error().Err(err).Msg("error while creating new AI Agent studio thread")
		return SendAIGatewayError(c, fiber.StatusInternalServerError, &aicore.AiGatewayError{
			Error: aicore.AiGatewayErrorDetail{
				Code:    strconv.Itoa(fiber.StatusInternalServerError),
				Message: err.Error(),
			},
		})
	}
	pCtx := c.UserContext()
	err = agent.Act(pCtx)
	if err != nil {
		x.logger.Error().Err(err).Msg("error while processing AI Agent studio request")
		return SendAIGatewayError(c, fiber.StatusInternalServerError, &aicore.AiGatewayError{
			Error: aicore.AiGatewayErrorDetail{
				Code:    strconv.Itoa(fiber.StatusInternalServerError),
				Message: err.Error(),
			},
		})
	}
	result := agent.GetResult()
	respBytes, err := json.Marshal(result)
	if err != nil {
		x.logger.Error().Err(err).Msg("error while marshalling the response")
		return SendAIGatewayError(c, fiber.StatusInternalServerError, &aicore.AiGatewayError{
			Error: aicore.AiGatewayErrorDetail{
				Code:    strconv.Itoa(fiber.StatusInternalServerError),
				Message: err.Error(),
			},
		})
	}
	return c.Status(fiber.StatusOK).Send(respBytes)
}

func (x *AIGatewayServices) p2stream(c *fiber.Ctx, opts ...aiinteface.GwOpts) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")
	gwChan := make(chan *aicore.GwChatCompletionChunk, 100)
	opts = append(opts, aiinteface.WithGwStream(true), aiinteface.WithGwGatewayChannel(gwChan))
	agent, err := aiinteface.NewAIGateway(c.UserContext(), x.logger, opts...)
	if err != nil {
		x.logger.Error().Err(err).Msg("error while creating new AI Agent studio thread")
		return SendAIGatewayError(c, fiber.StatusInternalServerError, &aicore.AiGatewayError{
			Error: aicore.AiGatewayErrorDetail{
				Code:    strconv.Itoa(fiber.StatusInternalServerError),
				Message: err.Error(),
			},
		})
	}
	pCtx := c.UserContext()
	go func() {
		if err := agent.Act(pCtx); err != nil {
			x.logger.Error().Err(err).Msg("error while processing AI Agent studio request")
		}
		close(gwChan)
	}()

	c.Status(fiber.StatusOK)

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		for obj := range gwChan {
			respBytes, err := json.Marshal(obj)
			if err != nil {
				x.logger.Error().Err(err).Msg("error while marshalling the response")
				continue
			}
			sseEvent := "data: " + string(respBytes) + "\n\n"
			if _, err := w.WriteString(sseEvent); err != nil {
				x.logger.Error().Err(err).Msg("error while writing the response")
				continue
			}
			if err := w.Flush(); err != nil {
				x.logger.Error().Err(err).Msg("error while flushing the response")
				continue
			}
		}
		endEvent := "data: [DONE]\n\n"
		if _, err := w.WriteString(endEvent); err != nil {
			x.logger.Error().Err(err).Msg("error while writing the end response")
		}
		if err := w.Flush(); err != nil {
			x.logger.Error().Err(err).Msg("error while flushing the end response")
		}
	}))
	return nil
}
