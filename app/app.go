package app

import (
	"database/sql"
	"github.com/TechBuilder-360/Auth_Server/config"
	"github.com/TechBuilder-360/Auth_Server/logger"
	"github.com/gorilla/mux"
	"sync"
)

type App struct {
	Config config.Data
	DB 	   *sql.DB
	Router *mux.Router
}

var (
	once sync.Once
)

func (app *App) RegisterRoutes() {

	once.Do(func() {

		//controller := controller.Controller{
		//	Service: service.NewService(repo.NewRepositoryDb(app.DB, app.Config), app.Config),
		//	Config:  app.Config,
		//	DB:      app.DB,
		//}
		//
		//subRouter := app.Router.PathPrefix("/api/v1").Subrouter()
	})

	logger.Info("App routes registered successfully!")
}
