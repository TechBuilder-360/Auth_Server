package controllers

import (
	"github.com/TechBuilder-360/Auth_Server/internal/common/utils"
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	Ping(router *fiber.Ctx) error
	RegisterRoutes(router *fiber.App)
}

func (c *NewController) RegisterRoutes(router *fiber.App) {
	router.Get("/", c.Ping)
}

type NewController struct {
}

func DefaultController() Controller {
	return &NewController{}
}

func (c *NewController) Ping(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON(utils.SuccessResponse{
		Status:  true,
		Message: "We are up and running 🚀",
	})
}
