package error

import "errors"

var (
	ErrInternalServerError = errors.New("internal Server Error")
	ErrSqlError            = errors.New("database error")
	ErrToManyRequest       = errors.New("to many requests")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInvalidToken        = errors.New("invalid token")
	ErrForbidden           = errors.New("forbidden")
)

var GeneralErrors = []error{
	ErrInternalServerError, ErrSqlError, ErrToManyRequest, ErrUnauthorized, ErrInvalidToken, ErrForbidden,
}
