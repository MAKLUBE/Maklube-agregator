package mongo

import (
	"context"
	"errors"

	"github.com/MAKLUBE/AP1_Final_Project/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RestaurantStore struct {
	col *mongo.Collection
}

func NewRestaurantStore(db *mongo.Database) *RestaurantStore {
	return &RestaurantStore{col: db.Collection("restaurants")}
}

func (s *RestaurantStore) List(ctx context.Context) ([]models.Restaurant, error) {
	cur, err := s.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []models.Restaurant
	for cur.Next(ctx) {
		var r models.Restaurant
		if err := cur.Decode(&r); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, cur.Err()
}

func (s *RestaurantStore) FindByID(ctx context.Context, id any) (*models.Restaurant, error) {
	var r models.Restaurant
	err := s.col.FindOne(ctx, bson.M{"_id": id}).Decode(&r)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("not found")
	}
	return &r, err
}
