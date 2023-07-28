package character

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(r *gin.Engine, db *mongo.Database) {
	r.GET("/character/api/v2/characters/:cid", func(c *gin.Context) {
		cid := c.Param("cid")
		GetAccountCharacter(c, db, cid)
	})

	r.POST("/character/api/v2/characters", func(c *gin.Context) {
		CreateUserCharacter(c, db)
	})
}
