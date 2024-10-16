package middlewares

import (
	"context"
	"github.com/TechBuilder-360/Auth_Server/internal/common/utils"
	"github.com/TechBuilder-360/Auth_Server/pkg/log"
	"github.com/gofiber/fiber/v2"
)

const XRequestID = "X-Request-ID"
const RequestID = "Request-ID"

func Logger(c *fiber.Ctx) error {
	// Set a custom header on all responses:
	requestID := utils.GenerateUUID()
	c.Set(XRequestID, requestID)

	ctx := context.Background()

	logger := log.WithField(RequestID, requestID)
	ctx = context.WithValue(ctx, log.LoggerInCtx, logger)

	c.SetUserContext(ctx)

	logger.Info("Request: %s", string(c.Body()))

	return c.Next()
}
