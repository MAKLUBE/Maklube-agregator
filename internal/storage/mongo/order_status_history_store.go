package mongo

import (
	"context"

	"github.com/MAKLUBE/AP1_Final_Project/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderStatusHistoryStore struct {
	col *mongo.Collection
}

func NewOrderStatusHistoryStore(db *mongo.Database) *OrderStatusHistoryStore {
	return &OrderStatusHistoryStore{col: db.Collection("order_status_history")}
}

func (s *OrderStatusHistoryStore) Insert(ctx context.Context, h *models.OrderStatusHistory) error {
	_, err := s.col.InsertOne(ctx, h)
	return err
}

func (s *OrderStatusHistoryStore) ListByOrder(ctx context.Context, orderID primitive.ObjectID) ([]models.OrderStatusHistory, error) {
	cur, err := s.col.Find(ctx, bson.M{"order_id": orderID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []models.OrderStatusHistory
	for cur.Next(ctx) {
		var h models.OrderStatusHistory
		if err := cur.Decode(&h); err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	return out, cur.Err()
}
