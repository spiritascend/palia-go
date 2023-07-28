package account

import (
	"context"
	entitlements "palia-go/Controllers/Entitlements"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Account struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Email     string             `bson:"email" json:"email,omitempty"`
	AccountID string             `bson:"id" json:"account_id,omitempty"`
	Username  string             `bson:"username" json:"username,omitempty"`
}

func CreateAccount(c *gin.Context, db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var createAccountPayload Account

	if err := c.BindJSON(&createAccountPayload); err != nil {
		c.JSON(500, gin.H{"error": "Failed to parse JSON data"})
		return
	}

	accountscollection := db.Collection("accounts")

	if accountscollection == nil {
		c.JSON(500, gin.H{"error": "Failed to access database collection"})
		return
	}

	var DupeResult Account

	filter := bson.M{
		"$or": []bson.M{
			{"email": createAccountPayload.Email},
			{"username": createAccountPayload.Username},
		},
	}

	err := accountscollection.FindOne(ctx, filter).Decode(&DupeResult)

	if err != mongo.ErrNoDocuments {
		c.JSON(403, gin.H{"error": "Email or Username Duplicate"})
		return
	} else {
		createAccountPayload.AccountID = uuid.New().String()

		result, err := accountscollection.InsertOne(ctx, createAccountPayload)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		insertedID := result.InsertedID.(primitive.ObjectID)

		var insertedAccount Account
		err = accountscollection.FindOne(ctx, bson.M{"_id": insertedID}).Decode(&insertedAccount)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		entitlements.CreateWallet(insertedAccount.AccountID, db)
		c.JSON(201, insertedAccount)
	}
}
