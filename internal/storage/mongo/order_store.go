package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/MAKLUBE/AP1_Final_Project/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderStore struct {
	col *mongo.Collection
}

func NewOrderStore(db *mongo.Database) *OrderStore {
	return &OrderStore{col: db.Collection("orders")}
}

func (s *OrderStore) Insert(ctx context.Context, o *models.Order) error {
	_, err := s.col.InsertOne(ctx, o)
	return err
}

func (s *OrderStore) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Order, error) {
	var o models.Order
	err := s.col.FindOne(ctx, bson.M{"_id": id}).Decode(&o)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("not found")
	}
	return &o, err
}

func (s *OrderStore) ListByCustomer(ctx context.Context, customerID primitive.ObjectID) ([]models.Order, error) {
	cur, err := s.col.Find(ctx, bson.M{"customer_id": customerID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []models.Order
	for cur.Next(ctx) {
		var o models.Order
		if err := cur.Decode(&o); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, cur.Err()
}

func (s *OrderStore) ListByRestaurant(ctx context.Context, restaurantID primitive.ObjectID) ([]models.Order, error) {
	cur, err := s.col.Find(ctx, bson.M{"restaurant_id": restaurantID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []models.Order
	for cur.Next(ctx) {
		var o models.Order
		if err := cur.Decode(&o); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, cur.Err()
}

func (s *OrderStore) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	_, err := s.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"order_status": status,
			"updated_at":   time.Now().UTC(),
		},
	})
	return err
}
