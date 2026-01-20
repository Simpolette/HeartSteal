package route

import (
	"time"

	"github.com/Simpolette/HeartSteal/server/internal/bootstrap"
	"github.com/Simpolette/HeartSteal/server/internal/middleware"

	"go.mongodb.org/mongo-driver/mongo"
	"github.com/gin-gonic/gin"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db *mongo.Database, gin *gin.Engine) {
	publicRouter := gin.Group("/api")
	// All Public APIs
	NewUserRouter(env, timeout, db, publicRouter)

	protectedRouter := gin.Group("")
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	// All Private APIs
}