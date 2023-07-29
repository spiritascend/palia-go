package entitlements

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Wallet struct {
	Premium_Balance int `bson:"premium_balance" json:"premium_balance"`
}

type Entitlement_Schema struct {
	Account_Id   string   `bson:"account_id" json:"account_id"`
	Wallet       Wallet   `bson:"wallet" json:"wallet"`
	Entitlements []string `bson:"entitlements" json:"entitlements"`
}

func GetEntitlementProfile(ctx context.Context, accountid string, entitlementcollection *mongo.Collection) (*Entitlement_Schema, error) {

	if entitlementcollection == nil {
		log.Fatal("Failed To Find Entitlement Collection")
		return nil, mongo.ErrNilValue
	}

	filter := bson.M{
		"account_id": strings.ToLower(accountid),
	}

	var ret Entitlement_Schema
	err := entitlementcollection.FindOne(ctx, filter).Decode(&ret)

	if err == nil {
		fmt.Println("Found Profile")
		return &ret, nil
	} else {
		fmt.Println("Didn't Find Profile")
		return nil, err
	}
}
func CreateEntitlementProfile(accountid string, db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	entitlementcollection := db.Collection("entitlements")

	if entitlementcollection == nil {
		log.Fatal("Failed To Find Entitlement Collection")
	}

	_, err := GetEntitlementProfile(ctx, accountid, entitlementcollection)

	if err == mongo.ErrNoDocuments {
		neweprofile := &Entitlement_Schema{accountid, Wallet{9999999}, []string{}}
		entitlementcollection.InsertOne(ctx, *neweprofile)
	} else {
		log.Fatal("Failed to Create Wallet: Duplicate")
	}
}

func AddEntitlementToProfile(ctx context.Context, entitlementcollection *mongo.Collection, accountid string, entitlementid string) bool {
	if entitlementcollection == nil {
		log.Fatal("Failed To Find Entitlement Collection")
	}

	filter := bson.M{
		"account_id": strings.ToLower(accountid),
		"entitlements": bson.M{
			"$nin": []string{entitlementid},
		},
	}

	update := bson.M{
		"$push": bson.M{
			"entitlements": entitlementid,
		},
	}

	result := entitlementcollection.FindOneAndUpdate(ctx, filter, update)

	if result.Err() != nil {
		fmt.Println(result.Err().Error())
	}
	return result.Err() == nil
}

func UpdatePremiumBalance(ctx context.Context, entitlementcollection *mongo.Collection, accountid string, amount int) bool {
	if entitlementcollection == nil {
		log.Fatal("Failed To Find Entitlement Collection")
	}

	filter := bson.M{
		"account_id": strings.ToLower(accountid),
	}

	update := bson.M{"$set": bson.M{"wallet.premium_balance": amount}}
	_, err := entitlementcollection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		fmt.Println(err)
	}

	return err == nil
}

func GetPremiumBalance(ctx context.Context, entitlementcollection *mongo.Collection, accountid string) int {
	entitlementprofile, err := GetEntitlementProfile(ctx, accountid, entitlementcollection)
	if err != nil {
		log.Fatal("Failed to Get Premium Balance (Entitlement Profile Error)")
	}
	return entitlementprofile.Wallet.Premium_Balance
}

func GetWallet(c *gin.Context, db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	entitlementcollection := db.Collection("entitlements")

	if entitlementcollection == nil {
		c.JSON(400, gin.H{"Error": "Failed To Find Entitlements Collection"})
		return
	}

	FoundProfile, err := GetEntitlementProfile(ctx, c.GetHeader("X-Authenticated-Character"), entitlementcollection)

	if err == mongo.ErrNoDocuments {
		c.JSON(400, gin.H{"Error": "Failed To Find Entitlements Profile"})
		return
	}

	c.JSON(200, gin.H{
		"account_id":      FoundProfile.Account_Id,
		"premium_balance": FoundProfile.Wallet.Premium_Balance,
	})
}

func GetEntitlements(c *gin.Context, db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	entitlementcollection := db.Collection("entitlements")

	if entitlementcollection == nil {
		c.JSON(400, gin.H{"Error": "Failed To Find Entitlements Collection"})
		return
	}

	FoundProfile, err := GetEntitlementProfile(ctx, c.GetHeader("X-Authenticated-Character"), entitlementcollection)

	if err == mongo.ErrNoDocuments {
		c.JSON(400, gin.H{"Error": "Failed To Find Entitlements Profile"})
		return
	}

	c.JSON(200, gin.H{
		"account_id":   FoundProfile.Account_Id,
		"entitlements": FoundProfile.Entitlements,
	})
}
