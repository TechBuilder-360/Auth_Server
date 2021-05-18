package main

import (
	"fmt"
	"github.com/TechBuilder-360/Auth_Server/app"
	"github.com/TechBuilder-360/Auth_Server/config"
	"github.com/TechBuilder-360/Auth_Server/database"
	"github.com/TechBuilder-360/Auth_Server/middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	APP := &app.App{}

	// APP setup
	Config := config.Data{}
	Config.Init("")
	APP.Config = Config

	APP.Router = mux.NewRouter()
	Database := &database.Database{
		Config: APP.Config,
	}

	Database.LoadDatabase()
	APP.DB = Database.DB
	APP.RegisterRoutes()


	middleware := middleware.New(handlers.CompressHandler(APP.Router), APP.DB).
		//UseClientValidation().
		//EnableCors().
		Build()

	serverAddress:=fmt.Sprintf("%s:%s",APP.Config.ServerHost, APP.Config.ServerPort)
	log.Printf("Server started on %s", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, middleware))
}


