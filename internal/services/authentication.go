package services

import (
	"errors"
	"fmt"
	"github.com/TechBuilder-360/Auth_Server/internal/common/constant"
	"github.com/TechBuilder-360/Auth_Server/internal/common/types"
	"github.com/TechBuilder-360/Auth_Server/internal/common/utils"
	"github.com/TechBuilder-360/Auth_Server/internal/configs"
	"github.com/TechBuilder-360/Auth_Server/internal/database/redis"
	"github.com/TechBuilder-360/Auth_Server/internal/infrastructure/sendgrid"
	"github.com/TechBuilder-360/Auth_Server/internal/model"
	"github.com/TechBuilder-360/Auth_Server/internal/repository"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

//go:generate mockgen -destination=../mocks/services/mockService.go -package=services github.com/TechBuilder-360/business-directory-backend/services UserService
type AuthService interface {
	RegisterUser(body *types.Registration, log *log.Entry) error
	ActivateEmail(token string, uid string, log *log.Entry) error
	Login(body *types.AuthRequest) (*types.LoginResponse, error)
	GenerateJWT(userID string) (*types.Authentication, error)
	ValidateToken(encodedToken string) (*jwt.Token, error)
	RequestToken(body *types.EmailRequest, logger *log.Entry) error
	RefreshUserToken(body types.RefreshTokenRequest, token string, logger *log.Entry) (*types.Authentication, error)
}

type authService struct {
	repo     repository.AuthRepository
	userRepo repository.UserRepository
	redis    *redis.Client
}

func NewAuthService() AuthService {
	return &authService{
		repo:     repository.NewAuthRepository(),
		userRepo: repository.NewUserRepository(),
		redis:    redis.RedisClient(),
	}
}

func (d *authService) ActivateEmail(token string, uid string, logger *log.Entry) error {
	user, err := d.userRepo.GetUserByID(uid)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	if user.EmailVerified {
		return errors.New("account is already active")
	}

	valid, err := d.repo.IsTokenValid(user.ID, token)
	if err != nil {
		logger.Error("An Error occurred when validating login token. %s", err.Error())
		return errors.New("account activation failed")
	}
	if !valid {
		logger.Error(err.Error())
		return errors.New("invalid activation link")
	}

	user.EmailVerified = true
	user.EmailVerifiedAt = time.Now()

	if err = d.userRepo.Update(user); err != nil {
		logger.Error("An Error occurred while Activating your account, Please try again. %s", err.Error())
		return errors.New("account activation failed")
	}

	err = d.redis.Delete(user.ID)
	if err != nil {
		logger.Error(err.Error())
	}

	return nil
}

func (d *authService) RegisterUser(body *types.Registration, log *log.Entry) error {
	body.EmailAddress = utils.ToLower(body.EmailAddress)
	// Check if email address exist
	existingUser, err := d.userRepo.GetByEmail(body.EmailAddress)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if existingUser != nil {
		log.Info("Email address already exist. '%s'", body.EmailAddress)
		return errors.New("account already exist")
	}

	// Save user details
	user := &model.User{
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		DisplayName:  utils.AddToStr(body.DisplayName),
		EmailAddress: body.EmailAddress,
		Tier:         uint8(0),
		PhoneNumber:  utils.AddToStr(body.PhoneNumber),
	}

	err = d.userRepo.Create(user)
	if err != nil {
		log.Error("error: occurred when saving new user. %s", err.Error())
		return errors.New("registration was not successful")
	}

	// OTP token
	var token string
	if configs.IsProduction() {
		token = utils.GenerateRandomString(20)
		// Send Activate email
		mailTemplate := &sendgrid.ActivationMailRequest{
			ToMail:   body.EmailAddress,
			ToName:   fmt.Sprintf("%s %s", body.LastName, body.FirstName),
			FullName: fmt.Sprintf("%s %s", body.LastName, body.FirstName),
			Token:    token,
			UID:      user.ID,
		}
		err = sendgrid.SendActivateMail(mailTemplate)
		if err != nil {
			log.Error("Error occurred when sending activation email. %s", err.Error())
		}

		err = d.redis.Set(user.ID, token, time.Hour*24)
		if err != nil {
			log.Error("Error occurred when when token %s", err)
		}
	} else {
		user.EmailVerified = true
		user.EmailVerifiedAt = time.Now()

		err = d.userRepo.Update(user)
		if err != nil {
			return err
		}
	}

	return nil
}

// Login
// Handles authentication logic
func (d *authService) Login(body *types.AuthRequest) (*types.LoginResponse, error) {
	response := &types.LoginResponse{}

	user, err := d.userRepo.GetByEmail(strings.ToLower(body.EmailAddress))
	if err != nil {
		log.Error("An error occurred when fetching user profile. %s", err.Error())
		return nil, errors.New(constant.InternalServerError)
	}
	if user == nil {
		return nil, errors.New("account not found")
	}

	// Validate user token
	token, err := d.redis.Get(user.ID)
	if err != nil {
		log.Error("An Error occurred when validating login token. %s", err.Error())
		return nil, errors.New("token validation failed")
	}

	if token == nil || utils.AddToStr(token) != body.Otp {
		return nil, errors.New("invalid OTP")
	}

	// Generate JWT for user
	jwt, err := d.GenerateJWT(user.ID)
	if err != nil {
		log.Error("An error occurred when generating jwt token. %s", err.Error())
		return nil, errors.New("authentication failed")
	}

	profile := types.UserProfile{
		ID:            user.ID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		DisplayName:   user.DisplayName,
		EmailAddress:  user.EmailAddress,
		PhoneNumber:   user.PhoneNumber,
		EmailVerified: user.EmailVerified,
		LastLogin:     user.LastLogin,
	}

	response.Profile = profile
	response.Authentication = *jwt

	return response, nil

}

func (d *authService) RequestToken(body *types.EmailRequest, logger *log.Entry) error {
	if !utils.ValidateEmail(body.EmailAddress) {
		return errors.New("invalid email address")
	}
	email := strings.ToLower(body.EmailAddress)

	// Check if email address exist
	user, err := d.userRepo.GetByEmail(email)
	if err != nil {
		logger.Error(err.Error())
		return errors.New(constant.InternalServerError)
	}

	if user == nil {
		logger.Error(err.Error())
		return errors.New("user not found")
	}

	var token string
	var duration uint = 5
	if configs.Instance.GetEnv() != configs.SANDBOX {
		token = utils.GenerateNumericToken(4)
		var duration uint = 5
		err = d.redis.Set(user.ID, token, time.Minute*time.Duration(duration))
		if err != nil {
			logger.Error("Error occurred when sending token %s", err)
			return errors.New("request failed please try again")

		}
	} else {
		token = "1234"
		err = d.redis.Set(user.ID, token, time.Minute*time.Duration(duration))
		if err != nil {
			logger.Error("Error occurred when sending token %s", err)
			return errors.New("request failed please try again")
		}
	}

	if configs.Instance.GetEnv() != configs.SANDBOX {
		mailTemplate := &sendgrid.OTPMailRequest{
			Code:     token,
			ToMail:   user.EmailAddress,
			ToName:   user.LastName + " " + user.FirstName,
			Name:     user.DisplayName,
			Duration: duration,
		}
		err = sendgrid.SendOTPMail(mailTemplate)
		if err != nil {
			logger.Error("Error occurred when sending otp email. %s", err.Error())
		}
	}

	return nil
}

type authCustomClaims struct {
	UserId string `json:"user_id"`
	RToken string `json:"token"`
	Token  jwt.StandardClaims
}

func (a authCustomClaims) Valid() error {
	//TODO implement me
	panic("implement me")
}

func (d *authService) GenerateJWT(userId string) (*types.Authentication, error) {
	refreshToken := utils.GenerateNumericToken(32)
	expireAt := time.Now().Add(time.Hour * 24)
	claims := &authCustomClaims{
		UserId: userId,
		RToken: refreshToken,
		Token: jwt.StandardClaims{
			ExpiresAt: expireAt.Unix(),
			Issuer:    configs.Instance.AppName,
			IssuedAt:  time.Now().Unix(),
		},
	}

	//encoded string
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	at, err := token.SignedString([]byte(fmt.Sprintf("%s-%s", configs.Instance.Secret, refreshToken)))
	if err != nil {
		log.Error(err.Error())
		return nil, errors.New("token could not be generated")
	}

	return &types.Authentication{
		AccessToken: at,
		ExpireAt:    expireAt.Second(),
	}, nil
}

func (d *authService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {

		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, errors.New("invalid token")
		}
		return []byte(configs.Instance.Secret), nil
	})
}

func (d *authService) RefreshUserToken(body types.RefreshTokenRequest, token string, logger *log.Entry) (*types.Authentication, error) {
	var userId string
	rToken, err := NewAuthService().ValidateToken(body.RefreshToken)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			logger.Error(jwt.ErrSignatureInvalid)
			return nil, errors.New(constant.InternalServerError)
		}
	}

	tk, err := NewAuthService().ValidateToken(token)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			logger.Error(jwt.ErrSignatureInvalid)
			return nil, errors.New(constant.InternalServerError)
		}
	}

	if claims, ok := rToken.Claims.(jwt.MapClaims); ok && rToken.Valid {
		userId = claims["user_id"].(string)

		if claims, ok = tk.Claims.(jwt.MapClaims); ok {
			uid := claims["user_id"].(string)

			if uid != userId {
				return nil, errors.New("invalid refresh token")
			}

			// Todo: check refresh or JWT has been revoked
		}
	} else {
		return nil, errors.New("refresh token has expired")
	}

	// Todo: Invalidate old refresh and JWT
	var response *types.Authentication
	response, err = NewAuthService().GenerateJWT(userId)
	if err != nil {
		log.Error(http.StatusInternalServerError)
		return nil, errors.New(constant.InternalServerError)
	}
	return response, nil
}
