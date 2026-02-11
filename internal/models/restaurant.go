package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Restaurant struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	Description   string             `bson:"description,omitempty" json:"description,omitempty"`
	Phone         string             `bson:"phone,omitempty" json:"phone,omitempty"`
	HalalStatus   string             `bson:"halal_status" json:"halal_status"`
	RatingAverage float64            `bson:"rating_average" json:"rating_average"`
	KaspiNumber   string             `bson:"kaspi_number,omitempty" json:"kaspi_number"`
	Address       struct {
		AddressText string `bson:"address_text" json:"address_text"`
		City        string `bson:"city" json:"city"`
		District    string `bson:"district,omitempty" json:"district,omitempty"`
		MapURL      string `bson:"map_url,omitempty" json:"map_url,omitempty"`
	} `bson:"address" json:"address"`

	OwnerUserID primitive.ObjectID `bson:"owner_user_id" json:"owner_user_id"`
	Instagram   string             `bson:"instagram,omitempty" json:"instagram,omitempty"`

	HalalProof primitive.ObjectID `bson:"halal_proof,omitempty" json:"halal_proof,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}
