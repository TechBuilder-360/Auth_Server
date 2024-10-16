package types

import "time"

type UserProfile struct {
	ID            string    `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	DisplayName   string    `json:"display_name"`
	EmailAddress  string    `json:"email_address"`
	PhoneNumber   *string   `json:"phone_number"`
	EmailVerified bool      `json:"email_verified"`
	LastLogin     time.Time `json:"last_login"`
}
