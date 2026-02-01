package mongo

import (
	"context"

	"github.com/MAKLUBE/AP1_Final_Project/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MenuItemStore struct {
	col *mongo.Collection
}

func NewMenuItemStore(db *mongo.Database) *MenuItemStore {
	return &MenuItemStore{col: db.Collection("menu_items")}
}

func (s *MenuItemStore) ListByRestaurant(ctx context.Context, restaurantID any) ([]models.MenuItem, error) {
	cur, err := s.col.Find(ctx, bson.M{"restaurant_id": restaurantID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []models.MenuItem
	for cur.Next(ctx) {
		var m models.MenuItem
		if err := cur.Decode(&m); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, cur.Err()
}
