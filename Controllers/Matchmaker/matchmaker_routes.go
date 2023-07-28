package matchmaker

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/matchmaker/api/v1/join", func(c *gin.Context) {
		JoinMatchmaker(c)
	})

	r.POST("/matchmaker/api/v1/join/status", func(c *gin.Context) {
		JoinMatchmakerStatus(c)
	})
}
