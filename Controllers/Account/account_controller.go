package account

import (
	"context"
	"fmt"
	entitlements "palia-go/Controllers/Entitlements"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Account_Schema struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Email     string             `bson:"email" json:"email,omitempty"`
	AccountID string             `bson:"id" json:"account_id,omitempty"`
	Username  string             `bson:"username" json:"username,omitempty"`
}

type Account struct {
	Email    string `bson:"email" json:"email,omitempty"`
	Username string `bson:"username" json:"username,omitempty"`
}

func firstEmptyField(a *Account) (bool, *string) {
	v := reflect.ValueOf(*a)

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := v.Type().Field(i)

		if reflect.DeepEqual(fieldValue.Interface(), reflect.Zero(fieldType.Type).Interface()) {
			return true, &fieldType.Name
		}
	}

	return false, nil
}
func CreateAccount(c *gin.Context, db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var requestPayload Account

	if err := c.BindJSON(&requestPayload); err != nil {
		c.JSON(500, gin.H{"error": "Failed to parse JSON data"})
		return
	}

	isempty, emptyfield := firstEmptyField(&requestPayload)

	if isempty {
		EmptyErrorMessage := fmt.Sprintf("Failed to create account because field {%s} is missing", *emptyfield)
		c.JSON(400, gin.H{"error": EmptyErrorMessage})
		return
	}

	accountscollection := db.Collection("accounts")

	if accountscollection == nil {
		c.JSON(500, gin.H{"error": "Failed to access database collection"})
		return
	}

	filter := bson.M{
		"$or": []bson.M{
			{"email": requestPayload.Email},
			{"username": requestPayload.Username},
		},
	}

	err := accountscollection.FindOne(ctx, filter).Decode(nil)

	if err != mongo.ErrNoDocuments {
		c.JSON(403, gin.H{"error": "Email or Username Duplicate"})
		return
	} else {

		var newaccount Account_Schema
		newaccount.Email = requestPayload.Email
		newaccount.Username = requestPayload.Username
		newaccount.AccountID = uuid.New().String()

		result, err := accountscollection.InsertOne(ctx, newaccount)
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
		entitlements.CreateEntitlementProfile(newaccount.AccountID, db)
		c.JSON(201, "Success")
	}
}
