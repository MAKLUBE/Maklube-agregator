package mongo

import (
	"context"
	"errors"

	"github.com/MAKLUBE/AP1_Final_Project/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MenuItemStore struct {
	col *mongo.Collection
}

func NewMenuItemStore(db *mongo.Database) *MenuItemStore {
	return &MenuItemStore{col: db.Collection("menu_items")}
}

func (s *MenuItemStore) ListByRestaurant(ctx context.Context, restaurantID primitive.ObjectID) ([]models.MenuItem, error) {
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

func (s *MenuItemStore) FindByID(ctx context.Context, id primitive.ObjectID) (*models.MenuItem, error) {
	var m models.MenuItem
	err := s.col.FindOne(ctx, bson.M{"_id": id}).Decode(&m)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("not found")
	}
	return &m, err
}

func (s *MenuItemStore) Insert(ctx context.Context, m *models.MenuItem) error {
	_, err := s.col.InsertOne(ctx, m)
	return err
}

func (s *MenuItemStore) Update(ctx context.Context, m *models.MenuItem) error {
	_, err := s.col.UpdateOne(ctx, bson.M{"_id": m.ID}, bson.M{
		"$set": bson.M{
			"name":          m.Name,
			"description":   m.Description,
			"category":      m.Category,
			"price":         m.Price,
			"photo_url":     m.PhotoURL,
			"prep_time_min": m.PrepTimeMin,
			"updated_at":    m.UpdatedAt,
			"is_available":  m.IsAvailable,
		},
	})
	return err
}
