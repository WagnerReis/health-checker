package errors

import "errors"

var (
	// User Errors
	ErrUserNotFound           = errors.New("user not found")
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrUserNameRequired       = errors.New("name is required")
	ErrUserPasswordRequired   = errors.New("password is required")
	ErrUserPasswordTooShort   = errors.New("password must be at least 8 characters long")
	ErrUserEmailAlreadyExists = errors.New("email already exists")
	ErrUserInvalidCredentials = errors.New("invalid credentials")

	// Refresh Token Errors
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)
