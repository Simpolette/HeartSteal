package bootstrap

import (
	"context"
	"fmt"
	"log"
	"time"

	"net/url"

	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongoDatabase(env *Env) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	dbUser := url.QueryEscape(env.DBUser)
	dbPass := url.QueryEscape(env.DBPass)

	mongodbURI := fmt.Sprintf("mongodb+srv://%s:%s@clusterdev.0ham19e.mongodb.net/?appName=ClusterDev", dbUser, dbPass)

	clientOptions := options.Client().ApplyURI(mongodbURI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary()) // Check Primary is reachable
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func CloseMongoDBConnection(client *mongo.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}