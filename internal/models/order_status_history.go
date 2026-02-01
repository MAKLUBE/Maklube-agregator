package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatusHistory struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderID   primitive.ObjectID `bson:"order_id" json:"order_id"`
	Status    string             `bson:"status" json:"status"`
	ChangedBy primitive.ObjectID `bson:"changed_by" json:"changed_by"`
	Comment   string             `bson:"comment,omitempty" json:"comment,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
