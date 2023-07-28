package chat

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/chat/api/v1/connection_info", func(c *gin.Context) {
		GetConnectionInfo(c)
	})
}
