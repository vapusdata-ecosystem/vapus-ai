package services

import (
	"github.com/gofiber/fiber/v2"
	aicore "github.com/vapusdata-ecosystem/vapusai/core/aistudio/core"
)

func SendAIGatewayError(c *fiber.Ctx, status int, errObj *aicore.AiGatewayError) error {
	return c.Status(status).JSON(errObj)
}
