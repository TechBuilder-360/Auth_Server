package main

import (
	"fmt"
	"github.com/TechBuilder-360/Auth_Server/docs"
	"github.com/TechBuilder-360/Auth_Server/internal/common/utils"
	"github.com/TechBuilder-360/Auth_Server/internal/configs"
	"github.com/TechBuilder-360/Auth_Server/internal/database"
	"github.com/TechBuilder-360/Auth_Server/internal/database/redis"
	"github.com/TechBuilder-360/Auth_Server/routers"
	"github.com/TechBuilder-360/Auth_Server/seeder"
	logrus_papertrail "github.com/polds/logrus-papertrail-hook"
	log "github.com/sirupsen/logrus"
	_ "github.com/swaggo/files"
	"net/http"
	"os"
	"strconv"
	"time"
)

// @title           Business directory API
// @version         1.0
// @description     This is the API for business directory api..

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /directory/api/v1

// @Security ApiKeyAuth
// @securityDefinitions.basic  ApiKeyAuth

func initLog() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	// Log as JSON instead of the default ASCII formatter.
	port, err := strconv.Atoi(utils.AddToStr(configs.Instance.PaperTailPort))
	if err != nil {
		return
	}

	hook, err := logrus_papertrail.NewPapertrailHook(&logrus_papertrail.Hook{
		Host:     "logs.papertrailapp.com",
		Port:     port,
		Hostname: "Auth-Server",
		Appname:  utils.AddToStr(configs.Instance.PaperTailAppName),
	})

	hook.SetLevels([]log.Level{log.ErrorLevel, log.WarnLevel})

	if err == nil {
		log.AddHook(hook)
	}
}

func main() {
	configs.Load()
	initLog()

	// Generate swagger doc information
	documentation()

	// set up redis DB
	redis.NewClient()
	dbConnection := database.ConnectDB()
	// migrate db models
	err := database.DBMigration(dbConnection)
	if err != nil {
		panic(fmt.Sprintf("Migration Failed: %s", err.Error()))
	}
	go seeder.Seed(dbConnection)

	// Setup cache
	//middlewares.ResponseCache()

	// Set up the routes
	router := routers.SetupRoutes()

	// Start the server
	log.Info("Server started on port %s:%s", configs.Instance.BASEURL, configs.Instance.Port)

	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", configs.Instance.BASEURL, configs.Instance.Port),
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err = s.ListenAndServe()
	if err != nil {
		log.Error(err.Error())
		return
	}
}

func documentation() {
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Business directory API"
	docs.SwaggerInfo.Description = "This is the API for business directory api."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s", configs.Instance.BASEURL)
	docs.SwaggerInfo.BasePath = fmt.Sprintf("/%s/api/v1", configs.Instance.BASEURL)
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}
