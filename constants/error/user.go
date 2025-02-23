package error

import "errors"

var (
	ErrNotFound           = errors.New("user not found")
	ErrInvalidPassword    = errors.New("password invalid")
	ErrUsernameExists     = errors.New("username already exists")
	ErrPasswordIsNotMatch = errors.New("password does not match")
)

var UserErrors = []error{
	ErrNotFound, ErrInvalidPassword, ErrUsernameExists, ErrPasswordIsNotMatch,
}
