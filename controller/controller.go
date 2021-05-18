package controller

import (
	"database/sql"
	"github.com/TechBuilder-360/Auth_Server/config"
)

type Controller struct {
	Config  config.Data
	DB		*sql.DB
}