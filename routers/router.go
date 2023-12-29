package routers

import (
	"github.com/TechBuilder-360/Auth_Server/internal/configs"
	"github.com/TechBuilder-360/Auth_Server/internal/controllers"
	"github.com/TechBuilder-360/Auth_Server/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	log "github.com/sirupsen/logrus"
)

func SetupRoutes() *fiber.App {
	router := fiber.New(fiber.Config{
		CaseSensitive: true,
		ErrorHandler:  middlewares.DefaultErrorHandler,
	})

	var (
		authController  = controllers.DefaultAuthController()
		usersController = controllers.DefaultUserController()
		controller      = controllers.DefaultController()
	)

	//*******************************************
	//******* Middlewares **********************
	//*******************************************
	router.Use(recover.New())

	//*******************************************
	//******* Controller **********************
	//*******************************************
	controller.RegisterRoutes(router)

	//*******************************************
	//******* Authentication **********************
	//*******************************************
	authController.RegisterRoutes(router)

	//*************************************
	//******* USERS **********************
	//*************************************
	usersController.RegisterRoutes(router)

	if configs.IsSandBox() {
		router.Get("/swagger/*", swagger.HandlerDefault)
	}

	log.Info("Routes have been initialized")

	return router
}
