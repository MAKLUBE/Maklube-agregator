package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RolePartner  UserRole = "partner"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Role         UserRole           `bson:"role" json:"role"`
	Username     string             `bson:"username" json:"username"`
	Email        string             `bson:"email" json:"email"`
	Phone        string             `bson:"phone,omitempty" json:"phone,omitempty"`
	PasswordHash []byte             `bson:"password_hash" json:"-"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}
