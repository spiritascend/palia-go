package main

import (
	"context"
	"log"
	"time"

	account "palia-go/Controllers/Account"
	character "palia-go/Controllers/Character"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	paliaDB := client.Database("Palia")

	r := gin.Default()

	r.POST("/register", func(c *gin.Context) {
		account.CreateAccount(c, paliaDB)
	})

	r.POST("/api/login", func(c *gin.Context) {
		account.HandleLogin(c, paliaDB)
	})

	r.GET("/auth-proxy/api/v1/auth/validate", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})

	r.GET("/character/api/v2/characters/:cid", func(c *gin.Context) {
		cid := c.Param("cid")
		character.GetAccountCharacter(c, paliaDB, cid)
	})

	r.Run("127.0.0.1:80")
}
