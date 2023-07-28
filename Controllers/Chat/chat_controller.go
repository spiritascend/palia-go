package chat

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type ChatSDK struct {
	Sdk_Api_Id int    `json:"sdk_api_id"`
	Signature  string `json:"signature"`
	User_Id    string `json:"user_id"`
}

type GConnectionInfoResp struct {
	Account_Id string `json:"account_id"`
}

func GetConnectionInfo(c *gin.Context) {
	var requestPayload GConnectionInfoResp
	if err := c.BindJSON(&requestPayload); err != nil {
		c.JSON(404, gin.H{"error": "Failed to parse Request Payload"})
		return
	}
	c.JSON(200, ChatSDK{70000083, "Signature", strings.ToLower(requestPayload.Account_Id)})
}
