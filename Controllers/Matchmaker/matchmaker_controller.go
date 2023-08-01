package matchmaker

import (
	"github.com/gin-gonic/gin"
)

func JoinMatchmaker(c *gin.Context) {
	c.JSON(200, gin.H{"number": 1})
}

type MatchmakerPlayerPayload struct {
	Player Player `json:"player"`
}
type Player struct {
	AccountID   string `json:"account_id"`
	CharacterID string `json:"character_id"`
}

type MatchmakingTicket struct {
	Ticket struct {
		TicketID string `json:"ticket_id"`
		Player   struct {
			AccountID   string `json:"account_id"`
			CharacterID string `json:"character_id"`
		} `json:"player"`
		Server struct {
			ServerID       string `json:"server_id"`
			ServerType     string `json:"server_type"`
			Version        string `json:"version"`
			Status         string `json:"status"`
			ConnectionInfo struct {
				Addr       string `json:"addr"`
				Port       int    `json:"port"`
				BeaconPort int    `json:"beacon_port"`
			} `json:"connection_info"`
		} `json:"server"`
		Expiry int `json:"expiry"`
	} `json:"ticket"`
	Signature string `json:"signature"`
}

func JoinMatchmakerStatus(c *gin.Context) {

	var requestPayload MatchmakerPlayerPayload

	if err := c.BindJSON(&requestPayload); err != nil {
		c.JSON(404, gin.H{"error": "Failed to parse Request Payload"})
		return
	}

	ticket := &MatchmakingTicket{
		Ticket: struct {
			TicketID string `json:"ticket_id"`
			Player   struct {
				AccountID   string `json:"account_id"`
				CharacterID string `json:"character_id"`
			} `json:"player"`
			Server struct {
				ServerID       string `json:"server_id"`
				ServerType     string `json:"server_type"`
				Version        string `json:"version"`
				Status         string `json:"status"`
				ConnectionInfo struct {
					Addr       string `json:"addr"`
					Port       int    `json:"port"`
					BeaconPort int    `json:"beacon_port"`
				} `json:"connection_info"`
			} `json:"server"`
			Expiry int `json:"expiry"`
		}{
			TicketID: "0",
			Player: struct {
				AccountID   string `json:"account_id"`
				CharacterID string `json:"character_id"`
			}{
				AccountID:   requestPayload.Player.AccountID,
				CharacterID: requestPayload.Player.CharacterID,
			},
			Server: struct {
				ServerID       string `json:"server_id"`
				ServerType     string `json:"server_type"`
				Version        string `json:"version"`
				Status         string `json:"status"`
				ConnectionInfo struct {
					Addr       string `json:"addr"`
					Port       int    `json:"port"`
					BeaconPort int    `json:"beacon_port"`
				} `json:"connection_info"`
			}{
				ServerID:   "Palia-Go",
				ServerType: "None",
				Version:    "1.0",
				Status:     "online",
				ConnectionInfo: struct {
					Addr       string `json:"addr"`
					Port       int    `json:"port"`
					BeaconPort int    `json:"beacon_port"`
				}{
					Addr:       "127.0.0.1",
					Port:       1337,
					BeaconPort: 1337,
				},
			},
			Expiry: 1630000000,
		},
		Signature: "0",
	}
	c.JSON(200, *ticket)
}
