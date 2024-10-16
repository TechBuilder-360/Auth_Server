package model

import "github.com/TechBuilder-360/Auth_Server/internal/common/types"

const (
	OWNER             = "Owner"
	OrganisationAdmin = "Organisation Admin"
	BranchManager     = "Branch Manager"
)

type Role struct {
	Base

	Name types.RoleType `gorm:"unique"`
}
