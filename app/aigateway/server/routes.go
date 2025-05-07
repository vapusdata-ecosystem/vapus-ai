package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vapusdata-ecosystem/vapusai/aigateway/services"
)

func chatRouter(x fiber.Router) {
	x.Options("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	x.Post("/v1/chat/completions", services.AIGatewayServicesManager.ChatCompletion)
	// x.Post("/v1/tasks/logs", services.AIGatewayServicesManager.RunningTaskStream)
}
