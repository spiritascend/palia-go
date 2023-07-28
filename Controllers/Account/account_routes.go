package account

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(r *gin.Engine, db *mongo.Database) {
	r.POST("/register", func(c *gin.Context) {
		CreateAccount(c, db)
	})

	r.POST("/api/login", func(c *gin.Context) {
		HandleLogin(c, db)
	})

	r.GET("/auth-proxy/api/v1/auth/validate", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})
}
