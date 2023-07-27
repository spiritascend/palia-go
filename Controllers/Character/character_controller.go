package character

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAccountCharacter(c *gin.Context, db *mongo.Database, cid string) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	characterCollection := db.Collection("characters")

	filter := bson.M{
		"characters": bson.M{
			"$elemMatch": bson.M{
				"account_id": cid,
			},
		},
	}

	var FetchedCharacterProfile characterCreationResponse
	err := characterCollection.FindOne(ctx, filter).Decode(&FetchedCharacterProfile)

	if err == mongo.ErrNoDocuments {
		c.JSON(200, make([]string, 0))
		return
	}

	response, err := json.Marshal(FetchedCharacterProfile.Characters)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.String(200, string(response))

}
