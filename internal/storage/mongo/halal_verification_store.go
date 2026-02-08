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

type HalalVerificationStore struct {
	col *mongo.Collection
}

func NewHalalVerificationStore(db *mongo.Database) *HalalVerificationStore {
	return &HalalVerificationStore{col: db.Collection("halal_verifications")}
}

func (s *HalalVerificationStore) Insert(ctx context.Context, v *models.HalalVerification) error {
	_, err := s.col.InsertOne(ctx, v)
	return err
}

func (s *HalalVerificationStore) FindByID(ctx context.Context, id primitive.ObjectID) (*models.HalalVerification, error) {
	var v models.HalalVerification
	err := s.col.FindOne(ctx, bson.M{"_id": id}).Decode(&v)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("not found")
	}
	return &v, err
}

func (s *HalalVerificationStore) ListPending(ctx context.Context) ([]models.HalalVerification, error) {
	cur, err := s.col.Find(ctx, bson.M{"status": "pending"})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []models.HalalVerification
	for cur.Next(ctx) {
		var v models.HalalVerification
		if err := cur.Decode(&v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, cur.Err()
}

func (s *HalalVerificationStore) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string, note string) error {
	_, err := s.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"status":      status,
			"reviewed_at": time.Now().UTC(),
			"review_note": note,
		},
	})
	return err
}
