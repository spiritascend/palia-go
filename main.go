package main

import (
	"context"
	"log"
	"os"
	"time"

	account "palia-go/Controllers/Account"
	catalog "palia-go/Controllers/Catalog"
	character "palia-go/Controllers/Character"
	chat "palia-go/Controllers/Chat"
	entitlements "palia-go/Controllers/Entitlements"
	matchmaker "palia-go/Controllers/Matchmaker"
	launcher "palia-go/Launcher"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"
)

func StartServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	//gin.SetMode(gin.ReleaseMode)
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
	catalog.RegisterRoutes(r, paliaDB)

	r.Run("127.0.0.1:80")
}

func main() {
	args := os.Args

	refuseserver := slices.Contains(args, "-noserver")

	if !refuseserver {
		go func() {
			StartServer()
		}()
	}
	launcher.IntiializeLauncher()
}
