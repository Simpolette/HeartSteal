package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Simpolette/HeartSteal/server/internal/bootstrap"
	"github.com/Simpolette/HeartSteal/server/internal/domain"
	"github.com/Simpolette/HeartSteal/server/internal/handler"
	"github.com/Simpolette/HeartSteal/server/internal/middleware"
	"github.com/Simpolette/HeartSteal/server/internal/repository"
	"github.com/Simpolette/HeartSteal/server/internal/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserRouter(env *bootstrap.Env, timeout time.Duration, db *mongo.Database, emailService domain.EmailService, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)
	sr := repository.NewSessionRepository(db, domain.CollectionSession)
	pr := repository.NewVerificationPINRepository(db, domain.CollectionVerificationPIN)

	uc := usecase.NewUserUseCase(
		ur, sr, pr, emailService, timeout,
		env.AccessTokenSecret, env.AccessTokenExpiryHour,
		env.RefreshTokenSecret, env.RefreshTokenExpiryHour,
	)
	h := handler.NewUserHandler(uc)

	// Public Routes
	group.POST("/signup", h.Signup)
	group.POST("/login", h.Login)
	group.POST("/forgot-password", h.ForgotPassword)
	group.POST("/verify-pin", h.VerifyPIN)
	group.POST("/reset-password", h.ResetPassword)

	// Protected Routes
	protected := group.Group("")
	protected.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	protected.POST("/signout", h.SignOut)
}