package mongo

import (
	"context"
	"time"

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

func (s *ReviewStore) ListAll(ctx context.Context) ([]models.Review, error) {
	cur, err := s.col.Find(ctx, bson.M{})
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

func (s *ReviewStore) SetHidden(ctx context.Context, id primitive.ObjectID, hidden bool, moderatorID primitive.ObjectID, note string) error {
	_, err := s.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"is_hidden":       hidden,
			"moderated_by":    moderatorID,
			"moderated_at":    time.Now().UTC(),
			"moderate_reason": note,
		},
	})
	return err
}
