package catalog

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(r *gin.Engine, db *mongo.Database) {
	r.GET("/catalog/api/v1/storefront/:cid", func(c *gin.Context) {
		cid := c.Param("cid")
		GetStaticCatalogStorefront(c, db, cid)
	})

	r.POST("/catalog/api/v1/purchase/:type", func(c *gin.Context) {
		HandleStorefrontPurchase(c, db)
	})
}
