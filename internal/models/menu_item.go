package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MenuItem struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RestaurantID primitive.ObjectID `bson:"restaurant_id" json:"restaurant_id"`
	Name         string             `bson:"name" json:"name"`
	Description  string             `bson:"description,omitempty" json:"description,omitempty"`
	Category     string             `bson:"category,omitempty" json:"category,omitempty"`
	Price        float64            `bson:"price" json:"price"`
	PhotoURL     string             `bson:"photo_url,omitempty" json:"photo_url,omitempty"`
	PrepTimeMin  int                `bson:"prep_time_min,omitempty" json:"prep_time_min,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
	IsAvailable  bool               `bson:"is_available" json:"is_available"`
}
