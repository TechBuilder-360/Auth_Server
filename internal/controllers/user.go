package controllers

import (
	"github.com/TechBuilder-360/Auth_Server/internal/common/constant"
	"github.com/TechBuilder-360/Auth_Server/internal/common/utils"
	"github.com/TechBuilder-360/Auth_Server/internal/services"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type IUserController interface {
	GetUser(ctx *fiber.Ctx) error
	GetUserByEmail(ctx *fiber.Ctx) error
	RegisterRoutes(router *fiber.App)
}

type UserController struct {
	as services.UserService
}

func (c *UserController) RegisterRoutes(router *fiber.App) {
	users := router.Group("/users")

	users.Get("", c.GetUserByEmail)
	users.Get("/:id", c.GetUser)

}

func DefaultUserController() IUserController {
	return &UserController{
		as: services.NewUserService(),
	}
}

// GetUser
// @Summary      Validate Token
// @Description  Validate Token
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200      {object}  utils.SuccessResponse{Data=types.UserProfile}
// @Router       /users/{id} [get]
func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	logger := log.WithFields(log.Fields{constant.RequestIdentifier: utils.GenerateUUID()})
	logger.Info("Get User")

	userId := ctx.Params("id")

	profile, err := c.as.GetUserByID(userId)
	if err != nil {
		logger.Error("error fetching user profile %s", err.Error())
		return ctx.Status(http.StatusOK).JSON(utils.ErrorResponse{
			Status:  false,
			Message: "request failed",
			Error:   err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(utils.SuccessResponse{
		Status:  true,
		Message: "success",
		Data:    profile,
	})
}

// GetUserByEmail
// @Summary      Get User by email
// @Description  Get User by email
// @Tags         Users
// @Produce      json
// @Success      200      {object}  utils.SuccessResponse{Data=types.UserProfile}
// @Router       /users [get]
func (c *UserController) GetUserByEmail(ctx *fiber.Ctx) error {
	logger := log.WithFields(log.Fields{constant.RequestIdentifier: utils.GenerateUUID()})
	logger.Info("Get User")

	email := ctx.Query("email")

	profile, err := c.as.GetUserByEmail(email)
	if err != nil {
		logger.Error("error fetching user profile %s", err.Error())
		return ctx.Status(http.StatusOK).JSON(utils.ErrorResponse{
			Status:  false,
			Message: "request failed",
			Error:   err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(utils.SuccessResponse{
		Status:  true,
		Message: "success",
		Data:    profile,
	})
}
