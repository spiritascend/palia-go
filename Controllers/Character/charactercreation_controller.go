package character

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateUserCharacter(c *gin.Context, db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	characterCollection := db.Collection("characters")

	var requestPayload characterCreationPayload

	if err := c.BindJSON(&requestPayload); err != nil {
		c.JSON(500, gin.H{"error": "Failed to parse Request Payload"})
		return
	}

	newCharacter := characterCreationRoot{
		Account_id:            requestPayload.Account_Id,
		Character_id:          uuid.NewString(),
		Character_Name:        requestPayload.Character_Name,
		Body_Type:             requestPayload.Body_Type,
		Customization_Options: requestPayload.Customization_Options,
		Current_loadout:       uuid.NewString(),
	}

	newCharacter.Loadouts = append(newCharacter.Loadouts, requestPayload.Loadout)
	newCharacter.Loadouts[0].Loadout_id = uuid.NewString()

	var resp_struct characterCreationResponse
	resp_struct.Characters = append(resp_struct.Characters, newCharacter)
	newCharacter.Loadouts[0].Set_current_loadout = nil

	insertOptions := options.InsertOne().SetBypassDocumentValidation(true)

	_, err := characterCollection.InsertOne(ctx, resp_struct, insertOptions)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response, err := json.Marshal(resp_struct)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.String(200, string(response))

}
