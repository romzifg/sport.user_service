package services

import (
	"context"
	"strings"
	"time"
	"user-service/config"
	"user-service/constants"
	errorConstant "user-service/constants/error"
	"user-service/domain/dto"
	"user-service/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository repositories.IRepositoryRegistry
}

type IUserService interface {
	Login(context.Context, *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(context.Context, *dto.RegiterRequest) (*dto.RegiterResponse, error)
	Update(context.Context, *dto.UpdateRequest, string) (*dto.UserResponse, error)
	UpdatePassword(context.Context, *dto.UpdatePasswordRequest, string) (*dto.UserResponse, error)
	GetUserLogin(context.Context) (*dto.UserResponse, error)
	GetUserByUUID(context.Context, string) (*dto.UserResponse, error)
}

type Claims struct {
	User *dto.UserResponse
	jwt.RegisteredClaims
}

func NewUserService(repository repositories.IRepositoryRegistry) IUserService {
	return &UserService{
		repository: repository,
	}
}

func (u *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.repository.GetUser().FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	expirationTime := time.Now().Add(time.Duration(config.Config.JwtExpirationTime) * time.Minute).Unix()
	data := &dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        strings.ToLower(user.Role.Code),
	}

	claims := &Claims{
		User: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime, 0)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.JwtSecret))
	if err != nil {
		return nil, err
	}

	response := &dto.LoginResponse{
		User:  *data,
		Token: tokenString,
	}

	return response, nil
}

func (u *UserService) IsUsernameExists(ctx context.Context, username string) bool {
	user, err := u.repository.GetUser().FindByUsername(ctx, username)
	if user != nil || err != nil {
		return true
	}

	return false
}

func (u *UserService) IsEmailExists(ctx context.Context, email string) bool {
	user, err := u.repository.GetUser().FindByEmail(ctx, email)
	if user != nil || err != nil {
		return true
	}

	return false
}

func (u *UserService) Register(ctx context.Context, req *dto.RegiterRequest) (*dto.RegiterResponse, error) {
	if req.Password != req.ConfirmPassword {
		return nil, errorConstant.ErrPasswordIsNotMatch
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if u.IsUsernameExists(ctx, req.Username) {
		return nil, errorConstant.ErrUsernameExists
	}

	if u.IsEmailExists(ctx, req.Email) {
		return nil, errorConstant.ErrEmailExists
	}

	user, err := u.repository.GetUser().Register(ctx, &dto.RegiterRequest{
		Username:    req.Username,
		Name:        req.Name,
		Password:    string(hashedPassword),
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		RoleID:      constants.Customer,
	})
	if err != nil {
		return nil, err
	}

	response := &dto.RegiterResponse{
		User: dto.UserResponse{
			UUID:        user.UUID,
			Name:        user.Name,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
		},
	}

	return response, nil
}

func (u *UserService) Update(ctx context.Context, req *dto.UpdateRequest, uuid string) (*dto.UserResponse, error) {
	user, err := u.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	usernameExists := u.IsUsernameExists(ctx, req.Username)
	if usernameExists && user.Username != req.Username {
		checkUsername, err := u.repository.GetUser().FindByUsername(ctx, req.Username)
		if err != nil {
			return nil, err
		}

		if checkUsername != nil {
			return nil, errorConstant.ErrUsernameExists
		}
	}

	emailExists := u.IsEmailExists(ctx, req.Email)
	if emailExists && user.Email != req.Email {
		checkEmail, err := u.repository.GetUser().FindByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}

		if checkEmail != nil {
			return nil, errorConstant.ErrEmailExists
		}
	}

	newUser, err := u.repository.GetUser().Update(ctx, &dto.UpdateRequest{
		Username:    req.Username,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
	}, uuid)
	if err != nil {
		return nil, err
	}

	response := &dto.UserResponse{
		UUID:        newUser.UUID,
		Name:        newUser.Name,
		Username:    newUser.Username,
		Email:       newUser.Email,
		PhoneNumber: newUser.PhoneNumber,
	}

	return response, nil
}

func (u *UserService) UpdatePassword(ctx context.Context, req *dto.UpdatePasswordRequest, uuid string) (*dto.UserResponse, error) {
	if req.NewPassword != req.ConfirmPassword {
		return nil, errorConstant.ErrPasswordIsNotMatch
	}
	user, err := u.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	_, err = u.repository.GetUser().UpdatePassword(ctx, &dto.UpdatePasswordRequest{
		NewPassword: string(hashedPassword),
	}, uuid)
	if err != nil {
		return nil, err
	}

	response := &dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}

	return response, nil
}

func (u *UserService) GetUserLogin(ctx context.Context) (*dto.UserResponse, error) {
	var (
		userLogin = ctx.Value(constants.UserLogin).(*dto.UserResponse)
		data      dto.UserResponse
	)

	data = dto.UserResponse{
		UUID:        userLogin.UUID,
		Name:        userLogin.Name,
		Username:    userLogin.Username,
		Email:       userLogin.Email,
		PhoneNumber: userLogin.PhoneNumber,
		Role:        userLogin.Role,
	}

	return &data, nil
}

func (u *UserService) GetUserByUUID(ctx context.Context, uuid string) (*dto.UserResponse, error) {
	user, err := u.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	response := &dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}

	return response, nil
}
