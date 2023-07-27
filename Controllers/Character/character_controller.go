package character

import (
	"context"
	account "palia-go/Controllers/Account"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAccountCharacter(c *gin.Context, db *mongo.Database, cid string) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	characterCollection := db.Collection("characters")

	filter := bson.M{"account_id": cid}

	var FetchedAccount account.Account
	err := characterCollection.FindOne(ctx, filter).Decode(&FetchedAccount)

	if err == mongo.ErrNoDocuments {
		c.JSON(200, make([]string, 0))
	}
}
