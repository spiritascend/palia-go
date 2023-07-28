package account

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoginPayload struct {
	ApplicationId string `bson:"applicationId"`
	LoginId       string `bson:"loginId"`
	Password      string `bson:"password"`
}

type LoginUserData struct {
	Country    string `json:"country"`
	Subscribed bool   `json:"subscribed"`
}

type LoginUserMembership struct {
	GroupId       string `json:"groupId"`
	Id            string `json:"id"`
	InsertInstant int    `json:"insertInstant"`
}

type LoginUserRegistration struct {
	ApplicationId      string   `json:"applicationId"`
	Id                 string   `json:"id"`
	InsertInstant      int      `json:"insertInstant"`
	LastLoginInstant   int      `json:"lastLoginInstant"`
	LastUpdateInstant  int      `json:"lastUpdateInstant"`
	PreferredLanguages []string `json:"perferredLanguages"`
	Roles              []string `json:"roles"`
	UsernameStatus     string   `json:"usernameStatus"`
	Verified           bool     `json:"verified"`
}
type LoginUser struct {
	Active                             bool                    `json:"active"`
	BirthDate                          string                  `json:"birthDate"`
	BreachedPasswordLastCheckedInstant int                     `json:"breachedPasswordLastCheckedInstant"`
	BreachedPasswordStatus             string                  `json:"breachedPasswordStatus"`
	ConnectorId                        string                  `json:"connectorId"`
	Data                               LoginUserData           `json:"data"`
	Email                              string                  `json:"Email"`
	Id                                 string                  `json:"id"`
	InsertInstant                      int                     `json:"insertInstant"`
	LastLoginInstant                   int                     `json:"lastLoginInstant"`
	LastUpdateInstant                  int                     `json:"lastUpdateInstant"`
	Memberships                        []LoginUserMembership   `json:"memberships"`
	PasswordChangeRequired             bool                    `json:"passwordChangeRequired"`
	PasswordLastUpdateInstant          int                     `json:"passwordLastUpdateInstant"`
	PreferredLanguages                 []string                `json:"preferredLanguages"`
	Registrations                      []LoginUserRegistration `json:"registrations"`
	TenantId                           string                  `json:"tenantId"`
	// Two Factor should go here but its not yet implemented in palia so ima ignore it
	UniqueUsername string `json:"uniqueUsername"`
	Username       string `json:"username"`
	UsernameStatus string `json:"usernameStatus"`
	Verified       bool   `json:"verified"`
}

type LoginResponse struct {
	RefreshToken           string    `json:"refreshToken"`
	RefreshTokenId         string    `json:"refreshTokenId"`
	Token                  string    `json:"token"`
	TokenExpirationInstant int       `json:"tokenExpirationInstant"`
	User                   LoginUser `json:"user"`
}

func HandleLogin(c *gin.Context, db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var requestPayload LoginPayload

	if err := c.BindJSON(&requestPayload); err != nil {
		c.JSON(500, gin.H{"error": "Failed to parse Request Payload"})
		return
	}

	accountscollection := db.Collection("accounts")

	if accountscollection == nil {
		c.JSON(500, gin.H{"error": "Failed to access accounts database"})
		return
	}

	filter := bson.M{"email": requestPayload.LoginId}

	var FetchedAccount Account_Schema
	err := accountscollection.FindOne(ctx, filter).Decode(&FetchedAccount)

	if err == mongo.ErrNoDocuments {
		c.JSON(404, gin.H{"error": "Account not Found"})
		return
	}

	token := base64.StdEncoding.EncodeToString([]byte(FetchedAccount.AccountID))

	loginuserregistration := []LoginUserRegistration{
		{
			ApplicationId:      requestPayload.ApplicationId,
			Id:                 "b3f6ade3-65ce-46a0-b546-13941b005a9c",
			InsertInstant:      0,
			LastLoginInstant:   0,
			LastUpdateInstant:  0,
			PreferredLanguages: []string{"en"},
			Roles:              []string{"palia_access", "player"},
			UsernameStatus:     "ACTIVE",
			Verified:           true,
		},
	}

	loginuser := LoginUser{
		Active:                             true,
		BirthDate:                          "0000-00-00",
		BreachedPasswordLastCheckedInstant: 0,
		BreachedPasswordStatus:             "None",
		ConnectorId:                        uuid.New().String(),
		Data: LoginUserData{
			Country:    "US",
			Subscribed: false,
		},
		Email:                     FetchedAccount.Email,
		Id:                        FetchedAccount.AccountID,
		InsertInstant:             1690375631449,
		LastLoginInstant:          1690375631449,
		LastUpdateInstant:         1690375631449,
		Memberships:               []LoginUserMembership{},
		PasswordChangeRequired:    false,
		PasswordLastUpdateInstant: 1690375631449,
		PreferredLanguages:        []string{"en"},
		Registrations:             loginuserregistration,
		TenantId:                  uuid.New().String(),
		UniqueUsername:            FetchedAccount.Username,
		Username:                  FetchedAccount.Username,
		UsernameStatus:            "ACTIVE",
		Verified:                  true,
	}

	loginresponse := &LoginResponse{
		RefreshToken:           token,
		RefreshTokenId:         uuid.New().String(),
		Token:                  token,
		TokenExpirationInstant: 1690375631449,
		User:                   loginuser,
	}

	response, err := json.Marshal(loginresponse)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to Login"})
		return
	}

	c.String(200, string(response))
}
