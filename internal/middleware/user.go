package middlewares

import (
	"context"
	"errors"
	"github.com/TechBuilder-360/Auth_Server/internal/common/utils"
	"github.com/TechBuilder-360/Auth_Server/internal/model"
	"github.com/TechBuilder-360/Auth_Server/internal/repository"
	"github.com/TechBuilder-360/Auth_Server/internal/services"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// AuthJWT authorise user JWT
func AuthJWT(c *fiber.Ctx) error {
	var user *model.User
	tokenString := ExtractBearerToken(c)
	if tokenString == "" {
		return c.Status(http.StatusUnauthorized).JSON(utils.ErrorResponse{
			Status:  false,
			Message: "missing authentication token",
		})
	}

	token, err := services.NewAuthService().ValidateToken(tokenString)
	if err != nil || token == nil {
		return c.Status(http.StatusUnauthorized).JSON(utils.ErrorResponse{
			Status:  false,
			Message: "authentication failed",
		})
	}

	user, err = repository.NewUserRepository().GetUserByID(token.ID)
	if err != nil {
		return err
	}

	ctx := context.Background()

	ctx = context.WithValue(ctx, AuthUserContextKey, user)

	c.SetUserContext(ctx)

	// Serve the next handler
	return c.Next()
}

func ExtractBearerToken(ctx *fiber.Ctx) string {
	const BearerSchema = "Bearer"
	authHeader := ctx.Get(fiber.HeaderAuthorization)
	if authHeader == "" {
		return ""
	}
	tokenString := authHeader[len(BearerSchema)+1:]
	return tokenString
}

func UserFromContext(r *http.Request) (*model.User, error) {
	u := r.Context().Value(AuthUserContextKey)

	if u == nil {
		return nil, errors.New("no user in context")
	}

	user := u.(*model.User)

	return user, nil
}
