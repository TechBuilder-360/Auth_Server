package database

import (
	"fmt"
	"github.com/TechBuilder-360/Auth_Server/internal/configs"
	"github.com/TechBuilder-360/Auth_Server/internal/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", configs.Instance.DbHost, configs.Instance.DbUser, configs.Instance.DbPass, configs.Instance.DbName, configs.Instance.DbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to DB. %s", err.Error())
	}
	return db
}

func DBMigration(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
	)

	return err
}
