package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID             uuid.UUID
	Username       string
	Role           string
	EmailAddr      string
	Password       string
	Token          string
	IsActive       bool
	IsLocked       bool
	PasswordMaxTry uint
	PasswordTries  uint
	Created        time.Time
	LastLogin      time.Time
}

type SocialUser struct {
	ID 		  uuid.UUID
	UserID    string    // Email address
	Provider  string    // Google
	Created   time.Time
}
