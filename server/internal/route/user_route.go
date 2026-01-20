package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Simpolette/HeartSteal/server/internal/bootstrap"
	"github.com/Simpolette/HeartSteal/server/internal/domain"
	"github.com/Simpolette/HeartSteal/server/internal/handler"
	// "github.com/Simpolette/HeartSteal/server/internal/middleware"
	"github.com/Simpolette/HeartSteal/server/internal/repository"
	"github.com/Simpolette/HeartSteal/server/internal/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserRouter(env *bootstrap.Env, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)
	uc := usecase.NewUserUseCase(ur, timeout, env.AccessTokenSecret, env.AccessTokenExpiryHour)
	h := handler.NewUserHandler(uc)

	// Public Routes
	group.POST("/signup", h.Signup)
	group.POST("/login", h.Login)

	// Private Routes
	// protected := group.Group("/users")
	
	// protected.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	
	// protected.GET("/profile", h.GetProfile)
}