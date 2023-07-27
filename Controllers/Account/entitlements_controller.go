package account

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Wallet struct {
	Account_Id      string `bson:"account_id" json:"account_id,omitempty"`
	Premium_Balance int    `bson:"premium_balance" json:"premium_balance,omitempty"`
}

func CreateWallet(accountid string, db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	walletcollection := db.Collection("wallets")

	var DupeWallet Wallet

	if walletcollection == nil {
		log.Fatal("Failed To Find Wallet Collection")
		return
	}

	filter := bson.M{
		"account_id": accountid,
	}
	err := walletcollection.FindOne(ctx, filter).Decode(&DupeWallet)

	if err == mongo.ErrNoDocuments {
		newWallet := Wallet{accountid, 999999}
		insertOptions := options.InsertOne().SetBypassDocumentValidation(true)
		walletcollection.InsertOne(ctx, newWallet, insertOptions)
	} else {
		log.Fatal("Failed to Create Wallet: Duplicate")
	}
}

func GetWallet(c *gin.Context, db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	walletcollection := db.Collection("wallets")

	if walletcollection == nil {
		c.JSON(400, gin.H{"Error": "Failed To Find Wallet Collection"})
		return
	}

	var FoundWallet Wallet

	fmt.Println(strings.ToLower(c.GetHeader("X-Authenticated-Character")))
	filter := bson.M{
		"account_id": strings.ToLower(c.GetHeader("X-Authenticated-Character")),
	}

	err := walletcollection.FindOne(ctx, filter).Decode(&FoundWallet)

	if err == mongo.ErrNoDocuments {
		c.JSON(400, gin.H{"Error": "Failed To Find Wallet"})
		return
	}

	c.JSON(200, gin.H{
		"account_id":      FoundWallet.Account_Id,
		"premium_balance": FoundWallet.Premium_Balance,
	})
}
