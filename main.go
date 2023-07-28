package main

import (
	"context"
	"log"
	"time"

	account "palia-go/Controllers/Account"
	character "palia-go/Controllers/Character"
	chat "palia-go/Controllers/Chat"
	entitlements "palia-go/Controllers/Entitlements"
	matchmaker "palia-go/Controllers/Matchmaker"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	r := gin.Default()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	paliaDB := client.Database("Palia")

	account.RegisterRoutes(r, paliaDB)
	entitlements.RegisterRoutes(r, paliaDB)
	character.RegisterRoutes(r, paliaDB)
	matchmaker.RegisterRoutes(r)
	chat.RegisterRoutes(r)

	r.Run("127.0.0.1:80")
}
