package services

import (
	"errors"
	"github.com/TechBuilder-360/Auth_Server/internal/common/constant"
	"github.com/TechBuilder-360/Auth_Server/internal/common/types"
	"github.com/TechBuilder-360/Auth_Server/internal/configs"
	"github.com/TechBuilder-360/Auth_Server/internal/infrastructure/cloudinary"
	"github.com/TechBuilder-360/Auth_Server/internal/model"
	"github.com/TechBuilder-360/Auth_Server/internal/repository"
	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -destination=../mocks/services/mockService.go -package=services github.com/TechBuilder-360/business-directory-backend/services UserService
type UserService interface {
	UpgradeTierOne(body *types.UpgradeUserTierRequest, user *model.User, logger *log.Entry) error
	Update(user *model.User) error
}

type DefaultUserService struct {
	userRepo repository.UserRepository
}

func (r *DefaultUserService) Update(user *model.User) error {
	return r.userRepo.Update(user)
}

func NewUserService() UserService {
	return &DefaultUserService{userRepo: repository.NewUserRepository()}
}

func (r *DefaultUserService) UpgradeTierOne(body *types.UpgradeUserTierRequest, user *model.User, logger *log.Entry) error {
	if user.Tier > 0 {
		return errors.New("tier upgrade has already being submitted")
	}

	if configs.Instance.GetEnv() != configs.SANDBOX {
		url, err := cloudinary.ImageUpload(body.IdentityImage)
		if err != nil {
			logger.Error("identity image failed to upload: %s", err.Error())
			return errors.New("identity image upload failed")
		}

		user.IdentityImage = &url
	} else {
		user.IdentityImage = &body.IdentityImage
	}

	// Todo: Identity needs to be verified
	user.IdentityName = &body.IdentityName
	user.IdentityNumber = &body.IdentityNumber
	user.Tier = uint8(1)

	err := r.Update(user)
	if err != nil {
		logger.Error(err)
		return errors.New(constant.InternalServerError)
	}
	return nil
}
