package controllers

import (
	"github.com/TechBuilder-360/Auth_Server/internal/services"
	"github.com/gofiber/fiber/v2"
)

type IUserController interface {
	RegisterRoutes(router *fiber.App)
}

type UserController struct {
	as services.UserService
}

func (c *UserController) RegisterRoutes(router *fiber.App) {
	_ = router.Group("/users")

}

func DefaultUserController() IUserController {
	return &UserController{
		as: services.NewUserService(),
	}
}
