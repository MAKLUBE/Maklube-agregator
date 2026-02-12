package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IngredientProof struct {
	Ingredient string   `bson:"ingredient" json:"ingredient"`
	ProofType  string   `bson:"proof_type,omitempty" json:"proof_type,omitempty"`
	ProofURLs  []string `bson:"proof_urls,omitempty" json:"proof_urls,omitempty"`
}

type HalalVerification struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RestaurantID     primitive.ObjectID `bson:"restaurant_id" json:"restaurant_id"`
	RequestedBy      primitive.ObjectID `bson:"requested_by" json:"requested_by"`
	Status           string             `bson:"status" json:"status"`
	MainProofType    string             `bson:"main_proof_type,omitempty" json:"main_proof_type,omitempty"`
	MainProofURLs    []string           `bson:"main_proof_urls,omitempty" json:"main_proof_urls,omitempty"`
	IngredientProofs []IngredientProof  `bson:"ingredient_proofs,omitempty" json:"ingredient_proofs,omitempty"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	ReviewedAt       time.Time          `bson:"reviewed_at,omitempty" json:"reviewed_at,omitempty"`
	ReviewNote       string             `bson:"review_note,omitempty" json:"review_note,omitempty"`
}
