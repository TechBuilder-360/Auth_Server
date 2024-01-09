package types

import "time"

// JWTResponse ...
type JWTResponse struct {
	AccessToken string      `json:"access_token"`
	Profile     UserProfile `json:"profile"`
}

type MailTemplate struct {
	Header   string
	Code     string
	Token    string
	ToEmail  string
	ToName   string
	Subject  string
	Duration time.Duration
}

// AuthRequest ...
type AuthRequest struct {
	EmailAddress string `json:"email_address" validate:"required,email"`
	Otp          string `json:"otp" validate:"required,len=6"`
}

// EmailRequest ...
type EmailRequest struct {
	EmailAddress string `json:"email_address" validate:"required,email"`
}

// Registration ...
type Registration struct {
	EmailAddress string  `json:"email_address" validate:"required,email"`
	Avatar       *string `json:"avatar"`
	FirstName    string  `json:"first_name" validate:"required"`
	LastName     string  `json:"last_name" validate:"required"`
	DisplayName  *string `json:"display_name"`
	PhoneNumber  *string `json:"phone_number" validate:"e164"`
}
