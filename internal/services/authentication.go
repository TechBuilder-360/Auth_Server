package services

import (
	"encoding/json"
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
	"github.com/TechBuilder-360/Auth_Server/pkg/log"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

//go:generate mockgen -destination=../mocks/services/mockService.go -package=services github.com/TechBuilder-360/business-directory-backend/services UserService
type AuthService interface {
	RegisterUser(body *types.Registration, log log.Entry) (*types.RegistrationResponse, *utils.AppError)
	ActivateEmail(token string, log log.Entry) error
	Login(body *types.AuthRequest) (*types.LoginResponse, error)
	generateJWT(userID string) (*types.Authentication, error)
	ValidateToken(encodedToken string) (*jwt.RegisteredClaims, error)
	RequestToken(body *types.EmailRequest, logger log.Entry) error
	RefreshUserToken(body *types.RefreshTokenRequest, logger log.Entry) (*types.Authentication, error)
	Logout(Token string) error
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

func (d *authService) ActivateEmail(token string, logger log.Entry) error {
	uid, err := d.repo.GetToken(token)
	if err != nil {
		logger.Error("An Error occurred when validating login token. %s", err.Error())
		return errors.New("account activation failed")
	}
	if uid == nil {
		logger.Error(err.Error())
		return errors.New("activation link has expired")
	}

	user, err := d.userRepo.GetUserByID(utils.AddToStr(uid))
	if err != nil || user == nil {
		return errors.New("account not found")
	}

	if user.EmailVerified {
		return nil
	}

	user.EmailVerified = true
	user.Active = true
	user.EmailVerifiedAt = time.Now()

	if err = d.userRepo.Update(user); err != nil {
		logger.Error("An Error occurred while Activating your account, Please try again. %s", err.Error())
		return errors.New("account activation failed")
	}

	err = d.redis.Delete(token)
	if err != nil {
		logger.Error(err.Error())
	}

	return nil
}

func (d *authService) RegisterUser(body *types.Registration, log log.Entry) (*types.RegistrationResponse, *utils.AppError) {
	body.EmailAddress = utils.ToLower(body.EmailAddress)
	// Check if email address exist
	existingUser, err := d.userRepo.GetByEmail(body.EmailAddress)
	if err != nil {
		log.Error(err.Error())
		return nil, &utils.AppError{
			Message: "request failed",
			Error:   err.Error(),
		}
	}

	if existingUser != nil {
		log.Info("Email address already exist. '%s'", body.EmailAddress)
		return nil, &utils.AppError{
			Message: "account already exist",
		}
	}

	// Save user details
	user := &model.User{
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		DisplayName:  utils.AddToStr(body.DisplayName),
		EmailAddress: body.EmailAddress,
		PhoneNumber:  body.PhoneNumber,
	}

	if !configs.IsProduction() {
		user.EmailVerified = true
		user.EmailVerifiedAt = time.Now()
		user.Active = true
	}

	err = d.userRepo.Create(user)
	if err != nil {
		log.Error("error: occurred when saving new user. %s", err.Error())
		return nil, &utils.AppError{
			Message: "registration was not successful",
		}
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

		err = d.repo.StoreToken(user.ID, token, 24*60)
		if err != nil {
			log.Error("Error occurred when when token %s", err)
		}
	}

	return &types.RegistrationResponse{UserID: user.ID}, nil
}

// Login
// Handles authentication logic
func (d *authService) Login(body *types.AuthRequest) (*types.LoginResponse, error) {
	response := new(types.LoginResponse)

	user, err := d.userRepo.GetByEmail(utils.ToLower(body.EmailAddress))
	if err != nil {
		log.Error("An error occurred when fetching user profile. %s", err.Error())
		return nil, errors.New(constant.InternalServerError)
	}
	if user == nil {
		return nil, errors.New("account not found")
	}

	if !user.Active {
		return nil, errors.New("account is inactive")
	}

	// Validate OTP token
	token, err := d.repo.GetToken(user.ID)
	if err != nil {
		log.Error("An Error occurred when validating login token. %s", err.Error())
		return nil, errors.New("token validation failed")
	}

	if token == nil {
		return nil, errors.New("invalid OTP")
	}

	err = bcrypt.CompareHashAndPassword([]byte(utils.AddToStr(token)), []byte(body.Otp))
	if err != nil {
		log.Error("An Error occurred when comparing login token. %s", err.Error())
		return nil, errors.New("invalid OTP")
	}

	// Generate JWT for user
	tk, err := d.generateJWT(user.ID)
	if err != nil {
		log.Error("An error occurred when generating jwt token. %s", err.Error())
		return nil, errors.New("request failed")
	}

	err = d.repo.DeleteToken(user.ID)
	if err != nil {
		log.Error("an error occurred when removing jwt token. %s", err.Error())
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
	response.Authentication = *tk

	defer func() {
		user.LastLogin = time.Now()
		d.userRepo.Update(user)
	}()

	return response, nil
}

func (d *authService) RequestToken(body *types.EmailRequest, logger log.Entry) error {
	if !utils.ValidateEmail(body.EmailAddress) {
		return errors.New("account not found")
	}
	email := strings.ToLower(body.EmailAddress)

	// Check if email address exist
	user, err := d.userRepo.GetByEmail(email)
	if err != nil {
		logger.Error(err.Error())
		return errors.New("request failed")
	}

	if user == nil {
		logger.Error(err.Error())
		return errors.New("user not found")
	}

	duration := uint(5)
	otp := "123456"

	if configs.IsProduction() {
		otp = utils.GenerateNumericToken(6)

		mailTemplate := &sendgrid.OTPMailRequest{
			Code:     otp,
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

	hashedOTP, _ := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	err = d.repo.StoreToken(user.ID, string(hashedOTP), duration)
	if err != nil {
		logger.Error("Error occurred when sending token %s", err)
		return errors.New("request failed please try again")
	}

	return nil
}

type AuthToken struct {
	Token        string
	RefreshToken string
}

func (d *authService) generateJWT(userId string) (*types.Authentication, error) {
	refreshToken := utils.GenerateNumericToken(32)
	rt, err := d.repo.GetToken(fmt.Sprintf("auth::%s", userId))
	if err != nil {
		return nil, err
	}

	if rt != nil {
		tk := AuthToken{}
		err = json.Unmarshal([]byte(utils.AddToStr(rt)), &tk)
		if err != nil {
			return nil, err
		}
		refreshToken = tk.RefreshToken
	}

	issuedAt := time.Now()
	expireAt := issuedAt.Add(time.Hour * 24)
	claims := jwt.RegisteredClaims{
		Issuer:    configs.Instance.AppName,
		ExpiresAt: jwt.NewNumericDate(expireAt),
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ID:        userId,
	}

	//encoded string
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	at, err := token.SignedString([]byte(fmt.Sprintf("%s-%s", configs.Instance.Secret, refreshToken)))
	if err != nil {
		log.Error(err.Error())
		return nil, errors.New("token could not be generated")
	}

	// Store Refresh token to enable revoking token 30 Days
	if rt == nil {
		tk := AuthToken{
			Token:        at,
			RefreshToken: refreshToken,
		}

		var marshal []byte
		marshal, err = json.Marshal(tk)
		if err != nil {
			return nil, err
		}

		err = d.repo.StoreToken(fmt.Sprintf("auth::%s", userId), string(marshal), 30*24*60)
		if err != nil {
			return nil, err
		}
	}
	return &types.Authentication{
		AccessToken:  at,
		RefreshToken: refreshToken,
		ExpireAt:     expireAt.Unix(),
	}, nil
}

func (d *authService) ValidateToken(encodedToken string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	tkn, err := jwt.ParseWithClaims(encodedToken, claims, func(token *jwt.Token) (any, error) {
		authToken, err := d.repo.GetToken(fmt.Sprintf("auth::%s", claims.ID))
		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		tk := AuthToken{}

		err = json.Unmarshal([]byte(utils.AddToStr(authToken)), &tk)
		if err != nil {
			return nil, err
		}

		if tk.Token != encodedToken {
			return nil, errors.New("jwt token not found")
		}

		key := fmt.Sprintf("%s-%s", configs.Instance.Secret, tk.RefreshToken)
		return []byte(key), nil
	})

	if err != nil || !tkn.Valid {
		return nil, err
	}

	return claims, nil
}

func (d *authService) RefreshUserToken(body *types.RefreshTokenRequest, logger log.Entry) (*types.Authentication, error) {
	claims, err := d.ValidateToken(body.Token)
	if err != nil {
		return nil, err
	}

	var response *types.Authentication
	response, err = d.generateJWT(claims.ID)
	if err != nil {
		return nil, errors.New("token could not be generated")
	}

	return response, nil
}

func (d *authService) Logout(Token string) error {
	claims, err := d.ValidateToken(Token)
	if err != nil {
		return err
	}

	// Invalidate Refresh token
	return d.repo.DeleteToken(fmt.Sprintf("auth::%s", claims.ID))
}
