package types

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
	Avatar       *string `json:"avatar" validate:"url"`
	FirstName    string  `json:"first_name" validate:"required"`
	LastName     string  `json:"last_name" validate:"required"`
	DisplayName  *string `json:"display_name"`
	PhoneNumber  *string `json:"phone_number" validate:"e164"`
}

type UpgradeUserTierRequest struct {
	IdentityNumber string `json:"identity_number" validate:"required" `
	IdentityName   string `json:"identity_name" validate:"required"`
	IdentityImage  string `json:"identity_image" validate:"required" `
}
