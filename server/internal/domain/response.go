package domain

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists  = errors.New("email already exists")

	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid access token")
	
	ErrInternalServerError = errors.New("internal server error")
)

type SuccessResponse struct {
	Message string      `json:"message"`         
	Data    interface{} `json:"data,omitempty"`  
}

type ErrorResponse struct {
	Message string `json:"message"`
}