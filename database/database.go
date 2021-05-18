package database

import (
	"database/sql"
	"fmt"
	"github.com/TechBuilder-360/Auth_Server/config"
	"github.com/TechBuilder-360/Auth_Server/logger"
	_ "github.com/lib/pq"
)

//Database : database struct
type Database struct {
	Config 	config.Data
	DB      *sql.DB
}

func (d Database) LoadDatabase()  {

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=verify-full",
		d.Config.PGUser, d.Config.PGPassword, d.Config.PGHost, d.Config.PGPort, d.Config.PGDatabase)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Error("Auth DB failed to start")
	} else{
		logger.Info("Auth DB started successfully")
	}

	d.DB = db
}