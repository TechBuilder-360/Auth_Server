package types

type (
	LoginResponse struct {
		Authentication Authentication `json:"authentication"`
		Profile        UserProfile    `json:"profile"`
	}

	OrganisationMember struct {
		ID string `json:"id"`
	}

	Authentication struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpireAt     int64  `json:"expire_at"`
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
