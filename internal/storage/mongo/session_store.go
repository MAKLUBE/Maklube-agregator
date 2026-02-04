package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/MAKLUBE/AP1_Final_Project/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionStore struct {
	col *mongo.Collection
}

func NewSessionStore(db *mongo.Database) *SessionStore {
	col := db.Collection("sessions")
	_ = col
	return &SessionStore{col: db.Collection("sessions")}
}

func (s *SessionStore) Create(ctx context.Context, sess *models.Session) error {
	sess.CreatedAt = time.Now().UTC()
	_, err := s.col.InsertOne(ctx, sess)
	return err
}

func (s *SessionStore) FindByToken(ctx context.Context, token string) (*models.Session, error) {
	var sess models.Session

	err := s.col.FindOne(ctx, bson.M{"token": token}).Decode(&sess)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("not found")
	}

	if time.Now().UTC().After(sess.ExpiresAt) {
		_ = s.Delete(ctx, token)
		return nil, errors.New("expired")
	}
	return &sess, err
}

func (s *SessionStore) Delete(ctx context.Context, token string) error {
	_, err := s.col.DeleteOne(ctx, bson.M{"token": token})
	return err
}
