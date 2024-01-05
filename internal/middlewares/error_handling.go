package middlewares

import (
	"errors"
	"fmt"
	"github.com/TechBuilder-360/Auth_Server/internal/common/utils"
	"github.com/gofiber/fiber/v2"
)

// Default error handler
var DefaultErrorHandler = func(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError
	message := err.Error()

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Error()
	}

	if code == fiber.StatusNotFound {
		message = fmt.Sprintf("route not found %s '%s'", c.Route().Method, c.Request().URI().Path())
	}

	// Set Content-Type: text/plain; charset=utf-8
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Return status code with error message
	return c.Status(code).JSON(utils.ErrorResponse{
		Status:  false,
		Message: "request failed",
		Error:   message,
	})
}
