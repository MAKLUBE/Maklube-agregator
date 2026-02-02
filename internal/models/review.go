package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RestaurantID primitive.ObjectID `bson:"restaurant_id" json:"restaurant_id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	Rating       int                `bson:"rating" json:"rating"`
	Comment      string             `bson:"comment,omitempty" json:"comment,omitempty"`
	IsHidden     bool               `bson:"is_hidden" json:"is_hidden"`
	ModeratedBy  primitive.ObjectID `bson:"moderated_by,omitempty" json:"moderated_by,omitempty"`
	ModeratedAt  time.Time          `bson:"moderated_at,omitempty" json:"moderated_at,omitempty"`
	ModerateNote string             `bson:"moderate_reason,omitempty" json:"moderate_reason,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}
