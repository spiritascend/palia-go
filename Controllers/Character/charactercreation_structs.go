package character

type character_name struct {
	Last  string `json:"last"`
	First string `json:"first"`
}

type AnimSet_variant struct {
}

type AnimSet struct {
	Id      string          `json:"id"`
	Variant AnimSet_variant `json:"variant"`
}

type SkinTone_variant struct {
}

type SkinTone struct {
	Id      string           `json:"id"`
	Variant SkinTone_variant `json:"variant"`
}

type Eyes_variant struct {
	IrisShape  string `json:"IrisShape"`
	PupilShape string `json:"PupilShape"`
}

type Eyes struct {
	Id      string       `json:"id"`
	Variant Eyes_variant `json:"variant"`
}

type HairStyle_variant struct {
	HairBaseColor string `json:"HairBaseColor"`
	HairTipColor  string `json:"HairTipColor"`
	EyebrowColor  string `json:"EyebrowColor"`
	HairTipMask   string `json:"HairTipMask"`
}

type HairStyle struct {
	Id      string            `json:"id"`
	Variant HairStyle_variant `json:"variant"`
}

type Head_variant struct {
}

type Head struct {
	Id      string       `json:"id"`
	Variant Head_variant `json:"variant"`
}

type FacialHair_variant struct {
}

type FacialHair struct {
	Id      string             `json:"id"`
	Variant FacialHair_variant `json:"variant"`
}

type Voice_variant struct {
}

type Voice struct {
	Id      string        `json:"id"`
	Variant Voice_variant `json:"variant"`
}

type customization_options struct {
	Head         Head       `json:"Head"`
	SkinTone     SkinTone   `json:"SkinTone"`
	Voice        Voice      `json:"Voice"`
	HairStyle    HairStyle  `json:"HairStyle"`
	AnimSet      AnimSet    `json:"AnimSet"`
	FacialHair   FacialHair `json:"FacialHair"`
	Eyes         Eyes       `json:"Eyes"`
	Body_Options *string    `json:"body_options,omitempty"`
}

type Torso_variant struct {
	TorsoBaseColor string `json:"TorsoBaseColor"`
}

type Torso struct {
	Id      string        `json:"id"`
	Variant Torso_variant `json:"variant"`
}

type Legs_variant struct {
}

type Legs struct {
	Id      string       `json:"id"`
	Variant Legs_variant `json:"variant"`
}

type Hat_variant struct {
}

type Hat struct {
	Id      string      `json:"id"`
	Variant Hat_variant `json:"variant"`
}

type FaceMask_variant struct {
}

type FaceMask struct {
	Id      string           `json:"id"`
	Variant FaceMask_variant `json:"variant"`
}

type FaceTattoo_variant struct {
}

type FaceTattoo struct {
	Id      string             `json:"id"`
	Variant FaceTattoo_variant `json:"variant"`
}

type BodyTattoo_variant struct {
}

type BodyTattoo struct {
	Id      string             `json:"id"`
	Variant BodyTattoo_variant `json:"variant"`
}

type Makeup_variant struct {
	MakeupParams string `json:"MakeupParams,omitempty"`
}

type Makeup struct {
	Id      string         `json:"id"`
	Variant Makeup_variant `json:"variant"`
}

type FaceComplexion_variant struct {
}

type FaceComplexion struct {
	Id      string                 `json:"id"`
	Variant FaceComplexion_variant `json:"variant"`
}

type BodyComplexion_variant struct {
}

type BodyComplexion struct {
	Id      string                 `json:"id"`
	Variant BodyComplexion_variant `json:"variant"`
}

type Glider_variant struct {
}

type Glider struct {
	Id      string         `json:"id"`
	Variant Glider_variant `json:"variant"`
}

type Pet_variant struct {
}

type Pet struct {
	Id      string      `json:"id"`
	Variant Pet_variant `json:"variant"`
}

type loadout_customization_options struct {
	Makeup         Makeup         `json:"Makeup"`
	Hat            Hat            `json:"Hat"`
	Torso          Torso          `json:"Torso"`
	BodyTattoo     BodyTattoo     `json:"BodyTattoo"`
	Glider         Glider         `json:"Glider"`
	FaceComplexion FaceComplexion `json:"FaceComplexion"`
	Legs           Legs           `json:"Legs"`
	Pet            Pet            `json:"Pet"`
	BodyComplexion BodyComplexion `json:"BodyComplexion"`
	FaceMask       FaceMask       `json:"FaceMask"`
	FaceTattoo     FaceTattoo     `json:"FaceTattoo"`
}

type loadout struct {
	Loadout_id            string                        `json:"loadout_id,omitempty"`
	Name                  string                        `json:"name"`
	Customization_options loadout_customization_options `json:"customization_options"`
	Set_current_loadout   *bool                         `json:"set_current_loadout,omitempty"`
}
type characterCreationPayload struct {
	Account_Id            string                `json:"account_id"`
	Character_Name        character_name        `json:"character_name"`
	Body_Type             int                   `json:"body_type"`
	Customization_Options customization_options `json:"customization_options"`
	Body_options          string                `json:"body_options,omitempty"`
	Loadout               loadout               `json:"loadout"`
}

type characterCreationRoot struct {
	Account_id            string                `json:"account_id"`
	Character_id          string                `json:"character_id"`
	Character_Name        character_name        `json:"character_name"`
	Body_Type             int                   `json:"body_type"`
	Customization_Options customization_options `json:"customization_options"`
	Current_loadout       string                `json:"current_loadout"`
	Loadouts              []loadout             `json:"loadouts"`
}

type characterCreationResponse struct {
	Characters []characterCreationRoot `json:",omitempty"`
}
