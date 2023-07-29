package entitlements

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(r *gin.Engine, db *mongo.Database) {
	r.GET("/entitlement/api/v1/wallet/:cid", func(c *gin.Context) {
		GetWallet(c, db)
	})

	r.GET("/entitlement/api/v1/entitlements/:cid", func(c *gin.Context) {
		GetEntitlements(c, db)
	})
}
