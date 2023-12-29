package repository

import (
	"github.com/TechBuilder-360/Auth_Server/internal/database"
	"github.com/TechBuilder-360/Auth_Server/internal/database/redis"
	"gorm.io/gorm"
	"time"
)

//go:generate mockgen -destination=../mocks/repository/auth.go -package=repository github.com/TechBuilder-360/business-directory-backend/repository AuthRepository
type AuthRepository interface {
	GetToken(token string) (*string, error)
	StoreToken(key, token string, minutes uint) error
	DeleteToken(key string) error
	WithTx(tx *gorm.DB) AuthRepository
}

type DefaultAuthRepo struct {
	db *gorm.DB
	r  *redis.Client
}

func (r *DefaultAuthRepo) GetToken(token string) (*string, error) {
	return r.r.Get(token)
}

func (r *DefaultAuthRepo) DeleteToken(key string) error {
	return r.r.Delete(key)
}

func (r *DefaultAuthRepo) StoreToken(key, token string, minutes uint) error {
	return r.r.Set(key, token, time.Minute*time.Duration(minutes))
}

func NewAuthRepository() AuthRepository {
	return &DefaultAuthRepo{
		db: database.ConnectDB(),
		r:  redis.RedisClient(),
	}
}

func (r *DefaultAuthRepo) WithTx(tx *gorm.DB) AuthRepository {
	return &DefaultAuthRepo{
		db: tx,
	}
}
