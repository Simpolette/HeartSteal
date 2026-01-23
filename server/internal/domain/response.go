package domain

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
)

type SuccessResponse struct {
	Message string      `json:"message"`         
	Data    interface{} `json:"data,omitempty"`  
}

type ErrorResponse struct {
	Message string `json:"message"`
}