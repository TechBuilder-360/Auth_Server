package types

import "time"

type (
	LoginResponse struct {
		Authentication Authentication `json:"authentication"`
		Profile        UserProfile    `json:"profile"`
	}

	OrganisationMember struct {
		ID string `json:"id"`
	}

	UserProfile struct {
		ID            string    `json:"id"`
		FirstName     string    `json:"first_name"`
		LastName      string    `json:"last_name"`
		DisplayName   string    `json:"display_name"`
		EmailAddress  string    `json:"email_address"`
		PhoneNumber   *string   `json:"phone_number"`
		EmailVerified bool      `json:"email_verified"`
		LastLogin     time.Time `json:"last_login"`
	}

	Authentication struct {
		AccessToken string `json:"access_token"`
		ExpireAt    int64  `json:"expire_at"`
	}

	RefreshTokenRequest struct {
		Token        string `json:"token" validate:"required"`
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	Activate struct {
		Status bool `json:"status"`
	}

	Query struct {
		Page     int    `json:"page" schema:"page"`
		PageSize int    `json:"page_size" schema:"page_size"`
		Search   string `json:"search" schema:"search"`
	}

	PaginatedResponse struct {
		Page    int         `json:"page"`
		PerPage int         `json:"per_page"`
		Total   int64       `json:"total"`
		Data    interface{} `json:"-"`
	}
)

type OrganisationSize string
type RoleType string
type Directory string
type Hash string
