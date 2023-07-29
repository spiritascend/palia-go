package catalog

import (
	"context"
	"log"
	entitlements "palia-go/Controllers/Entitlements"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserPurchase_Resp struct {
	AccountID     string   `json:"account_id"`
	RecipientID   string   `json:"recipient_id"`
	Price         int      `json:"price"`
	TransactionID string   `json:"transaction_id"`
	Contents      []string `json:"contents"`
}

type UserPurchase_Payload struct {
	Purchase Purchase `json:"purchase"`
	Contents Contents `json:"contents"`
}
type Purchase struct {
	AccountID   string `json:"account_id"`
	RecipientID string `json:"recipient_id"`
	Price       int    `json:"price"`
}
type Contents struct {
	EntitlementID string `json:"entitlement_id"`
}

func HandleStorefrontPurchase(c *gin.Context, db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var requestPayload UserPurchase_Payload

	if err := c.BindJSON(&requestPayload); err != nil {
		c.JSON(500, gin.H{"error": "Failed to parse Request Payload"})
		return
	}

	entitlementcollection := db.Collection("entitlements")

	if entitlementcollection == nil {
		log.Fatal("Failed To Find Entitlement Collection")
	}

	addedentitlement := entitlements.AddEntitlementToProfile(ctx, entitlementcollection, strings.ToLower(requestPayload.Purchase.AccountID), requestPayload.Contents.EntitlementID)

	if addedentitlement {

		response := UserPurchase_Resp{
			AccountID:     requestPayload.Purchase.AccountID,
			RecipientID:   requestPayload.Purchase.RecipientID,
			Price:         requestPayload.Purchase.Price,
			TransactionID: uuid.NewString(),
			Contents:      []string{requestPayload.Contents.EntitlementID},
		}
		previousamount := entitlements.GetPremiumBalance(ctx, entitlementcollection, requestPayload.Purchase.AccountID)
		updatebalance := entitlements.UpdatePremiumBalance(ctx, entitlementcollection, requestPayload.Purchase.AccountID, previousamount-requestPayload.Purchase.Price)
		if updatebalance {
			c.JSON(200, response)
		} else {
			c.JSON(400, gin.H{"error": "Failed to Update Account PremiumBalance"})
		}
		return
	} else {
		c.JSON(400, gin.H{"error": "Item Already Purchased"})
		return
	}

}
