package main

import (
	"fmt"
	"github.com/TechBuilder-360/Auth_Server/docs"
	"github.com/TechBuilder-360/Auth_Server/internal/common/utils"
	"github.com/TechBuilder-360/Auth_Server/internal/configs"
	"github.com/TechBuilder-360/Auth_Server/internal/database"
	"github.com/TechBuilder-360/Auth_Server/internal/database/redis"
	"github.com/TechBuilder-360/Auth_Server/routers"
	logrus_papertrail "github.com/polds/logrus-papertrail-hook"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

// @title           Authentication API
// @version         1.0
// @description     This is the API for Authentication api..

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /auth/v1

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
		panic(fmt.Sprintf("DB migration failed: %s", err.Error()))
	}

	//go seeder.Seed(dbConnection)

	// Set up the routes
	router := routers.SetupRoutes()

	// Start the server
	log.Info(fmt.Sprintf("Server started on %s:%s", configs.Instance.BASEURL, configs.Instance.Port))

	//s := &http.Server{
	//	Addr:           fmt.Sprintf("%s:%s", configs.Instance.BASEURL, configs.Instance.Port),
	//	Handler:        router,
	//	ReadTimeout:    30 * time.Second,
	//	WriteTimeout:   30 * time.Second,
	//	MaxHeaderBytes: 1 << 20,
	//}
	err = router.Listen(fmt.Sprintf("%s:%s", configs.Instance.BASEURL, configs.Instance.Port))
	if err != nil {
		log.Error(err.Error())
		return
	}
}

func documentation() {
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Authentication API"
	docs.SwaggerInfo.Description = "This is the API for Authentication api."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = configs.Instance.BASEURL
	docs.SwaggerInfo.BasePath = "/auth/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}
