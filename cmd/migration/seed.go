package migration

import (
	"github.com/TechBuilder-360/Auth_Server/internal/common/utils"
	"github.com/TechBuilder-360/Auth_Server/internal/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Seed the database with some data
func Seed(db *gorm.DB) {
	var errs []error
	errs = append(errs, runRolesSeeder(db))
	errs = append(errs, runCountrySeeder(db))

	for _, e := range errs {
		if e != nil {
			log.Errorf("migration error-> %v", e)
		}
	}
}

func runRolesSeeder(tx *gorm.DB) error {
	roles := []model.Role{
		{
			Name: "Owner",
		},
		{
			Name: "Organisation Admin",
		},
		{
			Name: "Branch Manager",
		},
	}

	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&roles).Error; err != nil {
		return err
	}

	return nil
}

func runCountrySeeder(tx *gorm.DB) error {
	country := []model.Country{
		{
			Name:   "Nigeria",
			Code:   "NG",
			Active: utils.ToBoolAddr(true),
		},
	}

	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&country).Error; err != nil {
		return err
	}

	return nil
}
