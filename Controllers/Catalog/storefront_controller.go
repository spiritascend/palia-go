package catalog

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	entitlements "palia-go/Controllers/Entitlements"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
)

type Storefront struct {
	CampaignID     string    `json:"campaign_id"`
	Metadata       Metadata  `json:"metadata"`
	Title          Title     `json:"title"`
	LocalizedTitle string    `json:"localized_title"`
	Widgets        []Widgets `json:"widgets"`
}
type Metadata struct {
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
}
type Title struct {
	EnUs string `json:"en_us"`
}
type Name struct {
	EnUs string `json:"en_us"`
}
type Items struct {
	ItemID        string `json:"item_id"`
	InternalName  string `json:"internal_name"`
	Name          Name   `json:"name"`
	LocalizedName string `json:"localized_name"`
	Tag           string `json:"tag"`
}
type Price struct {
	Price           int  `json:"price"`
	DiscountedPrice int  `json:"discounted_price"`
	Owned           bool `json:"owned"`
}
type Variant struct {
	VariantID     string  `json:"variant_id"`
	InternalName  string  `json:"internal_name"`
	Name          Name    `json:"name"`
	LocalizedName string  `json:"localized_name"`
	Items         []Items `json:"items"`
	Price         Price   `json:"price"`
}
type Variant_Contents struct {
	Variant Variant `json:"variant"`
}
type Set struct {
	SetID         string             `json:"set_id"`
	InternalName  string             `json:"internal_name"`
	Name          Name               `json:"name"`
	LocalizedName string             `json:"localized_name"`
	Contents      []Variant_Contents `json:"contents"`
	Price         Price              `json:"price"`
}
type Set_Contents struct {
	Set Set `json:"set"`
}
type Headline struct {
	EnUs string `json:"en_us"`
}
type Description struct {
	EnUs string `json:"en_us"`
}
type ForegroundURL struct {
	Masculine string `json:"masculine"`
	Feminine  string `json:"feminine"`
}
type Assets struct {
	BackgroundURL string        `json:"background_url"`
	ForegroundURL ForegroundURL `json:"foreground_url"`
}
type Widgets struct {
	WidgetType           string       `json:"widget_type"`
	WidgetSize           string       `json:"widget_size"`
	Contents             Set_Contents `json:"contents"`
	Price                Price        `json:"price"`
	Headline             Headline     `json:"headline"`
	Name                 Name         `json:"name"`
	Description          Description  `json:"description"`
	LocalizedHeadline    string       `json:"localized_headline"`
	LocalizedName        string       `json:"localized_name"`
	LocalizedDescription string       `json:"localized_description"`
	Assets               Assets       `json:"assets"`
}

func GetStaticCatalogStorefront(c *gin.Context, db *mongo.Database, accountid string) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var ItemShopResponse Storefront
	jsonData, err := os.ReadFile("itemshop.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Itemshop Static Not Found"})
		return
	}

	if err := json.Unmarshal(jsonData, &ItemShopResponse); err != nil {
		c.JSON(500, gin.H{"error": "Failed to parse JSON data"})
		return
	}

	entitlementcollection := db.Collection("entitlements")

	if entitlementcollection == nil {
		c.JSON(400, gin.H{"error": "Failed To Find Entitlement Collection"})
		return
	}

	entitlements, err := entitlements.GetEntitlementProfile(ctx, accountid, entitlementcollection)

	if err == mongo.ErrNoDocuments {
		c.JSON(400, gin.H{"error": "Failed to Find Entitlement Collection"})
	}

	for Widget := range ItemShopResponse.Widgets {

		ItemShopResponse.Widgets[Widget].Contents.Set.Price.Owned = slices.Contains(entitlements.Entitlements, strings.ToUpper(ItemShopResponse.Widgets[Widget].Contents.Set.SetID))

		for Variant := range ItemShopResponse.Widgets[Widget].Contents.Set.Contents {
			ItemShopResponse.Widgets[Widget].Contents.Set.Contents[Variant].Variant.Price.Owned = slices.Contains(entitlements.Entitlements, strings.ToUpper(ItemShopResponse.Widgets[Widget].Contents.Set.Contents[Variant].Variant.VariantID))
		}
	}

	response, err := json.Marshal(ItemShopResponse)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.Data(http.StatusOK, gin.MIMEJSON, response)
}
