package controllers

import (
	"github.com/TechBuilder-360/Auth_Server/internal/common/types"
	"github.com/TechBuilder-360/Auth_Server/internal/common/utils"
	"github.com/TechBuilder-360/Auth_Server/internal/middlewares"
	"github.com/TechBuilder-360/Auth_Server/internal/services"
	"github.com/TechBuilder-360/Auth_Server/internal/validation"
	"github.com/TechBuilder-360/Auth_Server/pkg/log"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type AuthController interface {
	Registration(ctx *fiber.Ctx) error
	ActivateEmail(ctx *fiber.Ctx) error
	Authenticate(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	RefreshUserToken(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
	ValidateToken(ctx *fiber.Ctx) error
	RegisterRoutes(router *fiber.App)
}

type NewAuthController struct {
	as services.AuthService
}

func (c *NewAuthController) RegisterRoutes(router *fiber.App) {
	apis := router.Group("auth")

	apis.Use(middlewares.Logger)

	apis.Post("/registration", c.Registration)
	apis.Get("/activate", c.ActivateEmail)
	apis.Post("/authentication", c.Authenticate)
	apis.Post("/login", c.Login)
	apis.Get("/validate-token", c.ValidateToken)
	apis.Post("/refresh", c.RefreshUserToken)
	apis.Put("/logout", c.Logout)
}

func DefaultAuthController() AuthController {
	return &NewAuthController{
		as: services.NewAuthService(),
	}
}

func (c *NewAuthController) Authenticate(ctx *fiber.Ctx) error {
	logger := log.LoggerInContext(ctx.UserContext())
	logger.Info("Authenticate")

	body := new(types.EmailRequest)
	err := ctx.BodyParser(body)
	if err != nil {
		return err
	}
	logger.Info("Request data: %+v", body)

	if err, ok := validation.ValidateStruct(body, logger); !ok {
		return ctx.Status(http.StatusBadRequest).JSON(utils.ValidationResponse(err))
	}

	err = c.as.RequestToken(body, logger)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(utils.ErrorResponse{
			Status:  false,
			Message: "request failed",
			Error:   err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(utils.SuccessResponse{
		Status:  true,
		Message: "Success",
	})
}

func (c *NewAuthController) Login(ctx *fiber.Ctx) error {
	logger := log.LoggerInContext(ctx.UserContext())
	logger.Info("Verify User email and send login token.")

	body := new(types.AuthRequest)

	err := ctx.BodyParser(body)
	if err != nil {
		return err
	}

	if err, ok := validation.ValidateStruct(body, logger); !ok {
		return ctx.Status(http.StatusBadRequest).JSON(utils.ValidationResponse(err))
	}

	response, err := c.as.Login(body)
	if err != nil {
		logger.Error(err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(utils.ErrorResponse{
			Status:  false,
			Message: err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(utils.SuccessResponse{
		Status:  true,
		Message: "Successful",
		Data:    response,
	})
}

func (c *NewAuthController) Registration(ctx *fiber.Ctx) error {
	logger := log.LoggerInContext(ctx.UserContext())
	logger.Info("Registration Request")

	body := new(types.Registration)

	err := ctx.BodyParser(body)
	if err != nil {
		return err
	}

	if err, ok := validation.ValidateStruct(body, logger); !ok {
		return ctx.Status(http.StatusBadRequest).JSON(utils.ValidationResponse(err))
	}

	resp, e := c.as.RegisterUser(body, logger)
	if e != nil {
		logger.Error("Message: %s, Error: %s", e.Error, e.Message)
		return ctx.Status(http.StatusBadRequest).JSON(utils.ErrorResponse{
			Status:  false,
			Message: e.Message,
			Error:   e.Error,
		})
	}

	return ctx.Status(http.StatusCreated).JSON(utils.SuccessResponse{
		Status:  true,
		Message: "Successful",
		Data:    resp,
	})
}

func (c *NewAuthController) ActivateEmail(ctx *fiber.Ctx) error {
	logger := log.LoggerInContext(ctx.UserContext())
	logger.Info("Activating User")

	token := ctx.Get("token")

	err := c.as.ActivateEmail(token, logger)
	if err != nil {
		logger.Error(err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(utils.ErrorResponse{
			Status:  false,
			Message: err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(utils.SuccessResponse{
		Status:  true,
		Message: "account activation successful",
	})
}

func (c *NewAuthController) RefreshUserToken(ctx *fiber.Ctx) error {
	logger := log.LoggerInContext(ctx.UserContext())
	logger.Info("refreshing user token")

	body := new(types.RefreshTokenRequest)

	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	if err, ok := validation.ValidateStruct(body, logger); !ok {
		return ctx.Status(http.StatusBadRequest).JSON(utils.ValidationResponse(err))
	}

	tk, err := c.as.RefreshUserToken(body, logger)
	if err != nil {
		logger.Error(err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(utils.ErrorResponse{
			Status:  false,
			Message: err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(utils.SuccessResponse{
		Status:  true,
		Message: "success",
		Data:    tk,
	})
}

func (c *NewAuthController) Logout(ctx *fiber.Ctx) error {
	logger := log.LoggerInContext(ctx.UserContext())
	logger.Info("Logout")

	err := c.as.Logout(middlewares.ExtractBearerToken(ctx))
	if err != nil {
		logger.Error(err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(utils.ErrorResponse{
			Status:  false,
			Message: "request failed",
		})
	}

	return ctx.Status(http.StatusOK).JSON(utils.SuccessResponse{
		Status:  true,
		Message: "success",
	})
}

func (c *NewAuthController) ValidateToken(ctx *fiber.Ctx) error {
	logger := log.LoggerInContext(ctx.UserContext())
	logger.Info("Validate Token")

	_, err := c.as.ValidateToken(middlewares.ExtractBearerToken(ctx))
	if err != nil {
		logger.Error(err.Error())
		return ctx.SendStatus(http.StatusUnauthorized)
	}

	return ctx.SendStatus(http.StatusOK)
}
