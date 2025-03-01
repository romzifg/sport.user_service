package dto

import "github.com/google/uuid"

type LoginRequest struct {
	Username string `json:"usernmae" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	PhoneNumber string    `json:"phone_number"`
}

type LoginResponse struct {
	User  UserResponse
	Token string `json:"token"`
}

type RegiterRequest struct {
	Username        string `json:"usernmae" validate:"required"`
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	PhoneNumber     string `json:"phone_number"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
	RoleID          uint
}

type RegiterResponse struct {
	User UserResponse `json:"user"`
}

type UpdateRequest struct {
	Username    string `json:"usernmae" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number"`
	RoleID      uint
}

type UpdatePasswordRequest struct {
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}
