package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Simpolette/HeartSteal/server/internal/domain"
)

type signupRequest struct {
	UserName    string `json:"username" binding:"required"`
	Email    	string `json:"email" binding:"required,email"`
	Password 	string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    	string `json:"email" binding:"required,email"`
	Password 	string `json:"password" binding:"required"`
}

type UserHandler struct {
	UserUseCase domain.UserUsecase
}

func NewUserHandler(usecase domain.UserUsecase) *UserHandler {
	return &UserHandler{
		UserUseCase: usecase,
	}
}

func (h *UserHandler) Signup(c *gin.Context) {
	var req signupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user := &domain.User{
		UserName:     	req.UserName,
		Email:    		req.Email,
		Password: 		req.Password,
	}

	err := h.UserUseCase.Register(c.Request.Context(), user)
	if err != nil {
		if err == domain.ErrEmailExists {
			c.JSON(http.StatusConflict, domain.ErrorResponse{Message: "Email already existed"})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, domain.SuccessResponse{Message: "User registered successfully"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	accessToken, err := h.UserUseCase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Login successfully",
		Data: gin.H{
			accessToken: accessToken,
		},
	})
}