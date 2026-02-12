package main

import (
	"time"
	"log"

	route "github.com/Simpolette/HeartSteal/server/internal/route"
	"github.com/Simpolette/HeartSteal/server/internal/bootstrap"
	"github.com/gin-gonic/gin"
)

func main() {

	app := bootstrap.App()

	env := app.Env

	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	gin := gin.Default()

	route.Setup(env, timeout, db, gin)

	if err := gin.Run(env.ServerAddress); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}