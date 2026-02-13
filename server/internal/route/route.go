package route

import (
	"time"

	"github.com/Simpolette/HeartSteal/server/internal/bootstrap"
	"github.com/Simpolette/HeartSteal/server/internal/infrastructure"
	"github.com/Simpolette/HeartSteal/server/internal/middleware"

	"go.mongodb.org/mongo-driver/mongo"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db *mongo.Database, gin *gin.Engine) {
	gin.SetTrustedProxies(nil)

	gin.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Create email service
	emailService := infrastructure.NewMailerSendService(env.MailerSendAPIKey, env.MailerSendFromEmail)
	
	publicRouter := gin.Group("/api")
	// All Public APIs
	NewUserRouter(env, timeout, db, emailService, publicRouter)

	protectedRouter := gin.Group("/api")
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	// All Private APIs
}