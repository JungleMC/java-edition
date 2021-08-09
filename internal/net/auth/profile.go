package auth

import "github.com/google/uuid"

type Profile struct {
	ID         uuid.UUID         `json:"id"`
	Name       string            `json:"name"`
	Properties []ProfileProperty `json:"properties"`
}

type ProfileProperty struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Signed    bool   `json:"-"`
	Signature string `json:"signature"`
}

type TextureProperties struct {
	Timestamp int64 `json:"timestamp"`
	ProfileId string `json:"profileId"`
	ProfileName string `json:"profileName"`
	SignatureRequired bool `json:"signatureRequired"`
	Textures Textures `json:"textures"`
}

type Textures struct {
	Skin *Texture `json:"SKIN,omitempty"`
	Cape *Texture `json:"CAPE,omitempty"`
}

type Texture struct {
	Url string `json:"url"`
}
