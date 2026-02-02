package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CustomerID   primitive.ObjectID `bson:"customer_id" json:"customer_id"`
	RestaurantID primitive.ObjectID `bson:"restaurant_id" json:"restaurant_id"`

	OrderStatus  string  `bson:"order_status" json:"order_status"`
	DeliveryType string  `bson:"delivery_type" json:"delivery_type"`
	DeliveryAddr string  `bson:"delivery_address,omitempty" json:"delivery_address,omitempty"`
	Subtotal     float64 `bson:"subtotal" json:"subtotal"`
	Total        float64 `bson:"total" json:"total"`
	DeliveryFee  float64 `bson:"delivery_fee" json:"delivery_fee"`
	PaymentType  string  `bson:"payment_type,omitempty" json:"payment_type,omitempty"`
	CustomerNote string  `bson:"customer_comment,omitempty" json:"customer_comment,omitempty"`

	Items []OrderItem `bson:"items" json:"items"`

	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
	ConfirmedAt time.Time `bson:"confirmed_at,omitempty" json:"confirmed_at,omitempty"`
	CompletedAt time.Time `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
}

type OrderItem struct {
	MenuItemID primitive.ObjectID `bson:"menu_item_id" json:"menu_item_id"`
	Name       string             `bson:"name" json:"name"`
	Price      float64            `bson:"price" json:"price"`
	Qty        int                `bson:"qty" json:"qty"`
	LineTotal  float64            `bson:"line_total" json:"line_total"`
}
