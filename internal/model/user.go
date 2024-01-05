package model

import (
	"time"
)

// User ...
type User struct {
	Base

	FirstName       string    `json:"first_name" gorm:"not null"`
	LastName        string    `json:"last_name" gorm:"not null"`
	DisplayName     string    `json:"display_name" gorm:"not null"`
	EmailAddress    string    `json:"email_address" gorm:"not null"`
	PhoneNumber     *string   `json:"phone_number" gorm:"null"`
	Avatar          *string   `json:"avatar" gorm:"null"`
	Active          bool      `json:"active" gorm:"default:false"`
	EmailVerified   bool      `json:"email_verified" gorm:"default:false"`
	EmailVerifiedAt time.Time `json:"email_verified_at"`
	LastLogin       time.Time `json:"last_login" gorm:"null"`
}
