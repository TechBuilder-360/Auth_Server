package middlewares

import (
	"errors"
	"github.com/TechBuilder-360/Auth_Server/internal/model"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// AuthorizeUserJWT authorise user JWT
//func AuthorizeUserJWT(ctx *fiber.Ctx) error {
//	var user *model.User
//	tokenString := ExtractBearerToken(ctx)
//	if tokenString == "" {
//		return ctx.Status(http.StatusUnauthorized).JSON(utils.ErrorResponse{
//			Status:  false,
//			Message: "missing authentication token",
//		})
//	}
//
//	token, err := services.NewAuthService().ValidateToken(tokenString)
//	if err != nil {
//		log.Error(err)
//		return ctx.Status(http.StatusUnauthorized).JSON(utils.ErrorResponse{
//			Status:  false,
//			Message: "authentication failed",
//		})
//	}
//	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//		userId := claims["user_id"].(string)
//		user, err = repository.NewUserRepository().GetUserByID(userId)
//		if err != nil {
//			log.Error(err)
//			return ctx.Status(http.StatusUnauthorized).JSON(utils.ErrorResponse{
//				Status:  false,
//				Message: "account not found",
//			})
//		}
//
//		ctx.Locals(AuthUserContextKey, user)
//
//	} else {
//		log.Error(err)
//		return ctx.Status(http.StatusUnauthorized).JSON(utils.ErrorResponse{
//			Status:  false,
//			Message: "unauthorized",
//		})
//	}
//
//	// Serve the next handler
//	return ctx.Next()
//}

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
