package services

import (
	"github.com/TechBuilder-360/Auth_Server/internal/common/types"
	"github.com/TechBuilder-360/Auth_Server/internal/model"
	"github.com/TechBuilder-360/Auth_Server/internal/repository"
)

//go:generate mockgen -destination=../mocks/services/mockService.go -package=services github.com/TechBuilder-360/business-directory-backend/services UserService
type UserService interface {
	Update(user *model.User) error
	GetUserByID(id string) (*types.UserProfile, error)
	GetUserByEmail(email string) (*types.UserProfile, error)
}

type DefaultUserService struct {
	userRepo repository.UserRepository
}

func (r *DefaultUserService) GetUserByID(id string) (*types.UserProfile, error) {
	user, err := r.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return &types.UserProfile{
		ID:            user.ID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		DisplayName:   user.DisplayName,
		EmailAddress:  user.EmailAddress,
		PhoneNumber:   user.PhoneNumber,
		EmailVerified: user.EmailVerified,
		LastLogin:     user.LastLogin,
	}, nil
}

func (r *DefaultUserService) GetUserByEmail(email string) (*types.UserProfile, error) {
	user, err := r.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return &types.UserProfile{
		ID:            user.ID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		DisplayName:   user.DisplayName,
		EmailAddress:  user.EmailAddress,
		PhoneNumber:   user.PhoneNumber,
		EmailVerified: user.EmailVerified,
		LastLogin:     user.LastLogin,
	}, nil
}

func (r *DefaultUserService) Update(user *model.User) error {
	return r.userRepo.Update(user)
}

func NewUserService() UserService {
	return &DefaultUserService{userRepo: repository.NewUserRepository()}
}
