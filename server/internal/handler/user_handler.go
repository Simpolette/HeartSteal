package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Simpolette/HeartSteal/server/internal/domain"
)

type UserHandler struct {
	UserUseCase domain.UserUsecase
}

func NewUserHandler(usecase domain.UserUsecase) *UserHandler {
	return &UserHandler{
		UserUseCase: usecase,
	}
}

func (h *UserHandler) Signup(c *gin.Context) {
	var req domain.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
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
	var req domain.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	accessToken, refreshToken, err := h.UserUseCase.Login(c.Request.Context(), req.Username, req.Password)
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
		Data: domain.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var req domain.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	err := h.UserUseCase.ForgotPassword(c.Request.Context(), req.Email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			// Return 200 even if user not found to prevent email enumeration
			c.JSON(http.StatusOK, domain.SuccessResponse{Message: "If the email exists, a PIN code has been sent"})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "If the email exists, a PIN code has been sent"})
}

func (h *UserHandler) VerifyPIN(c *gin.Context) {
	var req domain.VerifyPinRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	resetToken, err := h.UserUseCase.VerifyPIN(c.Request.Context(), req.Email, req.PinCode)
	if err != nil {
		if err == domain.ErrInvalidPIN || err == domain.ErrPINExpired {
			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "PIN verified successfully",
		Data: domain.VerifyPinResponse{
			ResetToken: resetToken,
		},
	})
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req domain.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	err := h.UserUseCase.ResetPassword(c.Request.Context(), req.ResetToken, req.NewPassword)
	if err != nil {
		if err == domain.ErrInvalidResetToken {
			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Password reset successfully"})
}

func (h *UserHandler) SignOut(c *gin.Context) {
	userID, exists := c.Get("x-user-id")
	if !exists {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Not authorized"})
		return
	}

	err := h.UserUseCase.SignOut(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Signed out successfully"})
}