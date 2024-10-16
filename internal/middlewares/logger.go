package middlewares

import (
	"context"
	"fmt"
	"github.com/TechBuilder-360/Auth_Server/internal/common/utils"
	"github.com/gofiber/fiber/v2"
)

const RequestID = "X-Request-ID"

func Logger(c *fiber.Ctx) error {
	// Set a custom header on all responses:
	c.Set(RequestID, utils.GenerateUUID())

	ctx := context.Background()

	// TODO: pass logger instance to user context
	// Add a User object to the context
	//user := User{Name: "John Doe"}
	//ctx = context.WithValue(ctx, "user", user)

	// Retrieve the User object from the context
	//userFromContext, ok := ctx.Value("user").(User)
	//if ok {
	//	fmt.Println(userFromContext.Name) // Output: John Doe
	//}
	c.SetUserContext(ctx)

	fmt.Println(string(c.Body()))

	c.Next()

	fmt.Println(string(c.Response().Body()))

	return nil
}
