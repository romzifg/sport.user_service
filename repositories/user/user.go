package repositories

import (
	"context"
	"errors"
	wrapError "user-service/common/error"
	errConstant "user-service/constants/error"
	"user-service/domain/dto"
	"user-service/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type IUserRepository interface {
	Register(context.Context, *dto.RegiterRequest) (*models.User, error)
	Update(context.Context, *dto.UpdateRequest, string) (*models.User, error)
	UpdatePassword(context.Context, *dto.UpdatePasswordRequest, string) (*models.User, error)
	FindByUsername(context.Context, string) (*models.User, error)
	FindByEmail(context.Context, string) (*models.User, error)
	FindByUUID(context.Context, string) (*models.User, error)
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Register(ctx context.Context, req *dto.RegiterRequest) (*models.User, error) {
	user := &models.User{
		UUID:        uuid.New(),
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		RoleId:      req.RoleID,
	}

	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, wrapError.WrapError(errConstant.ErrSqlError)
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, req *dto.UpdateRequest, uuid string) (*models.User, error) {
	user := &models.User{
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Updates(user).Error
	if err != nil {
		return nil, wrapError.WrapError(errConstant.ErrSqlError)
	}

	return user, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, req *dto.UpdatePasswordRequest, uuid string) (*models.User, error) {
	user := &models.User{
		Password: req.NewPassword,
	}

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Updates(user).Error
	if err != nil {
		return nil, wrapError.WrapError(errConstant.ErrSqlError)
	}

	return user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).Preload("Role").Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrNotFound
		}

		return nil, wrapError.WrapError(errConstant.ErrSqlError)
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).Preload("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrNotFound
		}

		return nil, wrapError.WrapError(errConstant.ErrSqlError)
	}

	return &user, nil
}

func (r *UserRepository) FindByUUID(ctx context.Context, uuid string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).Preload("Role").Where("uuid = ?", uuid).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrNotFound
		}

		return nil, wrapError.WrapError(errConstant.ErrSqlError)
	}

	return &user, nil
}
