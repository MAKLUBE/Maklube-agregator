package mongo

import (
	"context"

	"github.com/MAKLUBE/AP1_Final_Project/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewStore struct {
	col *mongo.Collection
}

func NewReviewStore(db *mongo.Database) *ReviewStore {
	return &ReviewStore{col: db.Collection("reviews")}
}

func (s *ReviewStore) ListByRestaurant(ctx context.Context, restaurantID primitive.ObjectID) ([]models.Review, error) {
	cur, err := s.col.Find(ctx, bson.M{"restaurant_id": restaurantID, "is_hidden": false})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []models.Review
	for cur.Next(ctx) {
		var r models.Review
		if err := cur.Decode(&r); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, cur.Err()
}

func (s *ReviewStore) Insert(ctx context.Context, r *models.Review) error {
	_, err := s.col.InsertOne(ctx, r)
	return err
}
