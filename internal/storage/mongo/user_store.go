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

type UserStore struct {
	col *mongo.Collection
}

func (s *UserStore) CreateRestaurant(ctx context.Context, restaurant *models.Restaurant) error {
	_, err := s.col.InsertOne(ctx, restaurant)
	if err != nil {
		return err
	}
	return nil
}

func NewUserStore(db *mongo.Database) *UserStore {
	return &UserStore{col: db.Collection("users")}
}

func (s *UserStore) Insert(ctx context.Context, u *models.User) error {
	u.CreatedAt = time.Now().UTC()
	_, err := s.col.InsertOne(ctx, u)
	return err
}

func (s *UserStore) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := s.col.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("not found")
	}
	return &u, err
}

func (s *UserStore) FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var u models.User
	err := s.col.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("not found")
	}
	return &u, err
}
