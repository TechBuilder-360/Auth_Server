package routers

import (
	"github.com/TechBuilder-360/Auth_Server/internal/controllers"
	"github.com/TechBuilder-360/Auth_Server/internal/middlewares"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/swaggo/http-swagger"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	var (
		authController  = controllers.DefaultAuthController()
		usersController = controllers.DefaultUserController()
		controller      = controllers.DefaultController()
	)

	router.Use(middlewares.Recovery)

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

	router.PathPrefix("/documentation/").Handler(httpSwagger.WrapHandler)

	log.Info("Routes have been initialized")

	return router
}
