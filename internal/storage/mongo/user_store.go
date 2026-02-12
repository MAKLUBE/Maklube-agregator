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

	res, err := s.col.InsertOne(ctx, u)
	if err != nil {
		return err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("insertedID is not ObjectID")
	}
	u.ID = oid

	return nil
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

func (s *UserStore) SetVerification(ctx context.Context, userID primitive.ObjectID, code string, expiresAt time.Time) error {
	res, err := s.col.UpdateOne(ctx,
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{
			"verification_code":       code,
			"verification_expires_at": expiresAt,
			"verified":                false,
		}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("user not found for verification update")
	}
	return nil
}

func (s *UserStore) VerifyByEmailAndCode(ctx context.Context, email string, code string, now time.Time) (*models.User, error) {
	var u models.User
	err := s.col.FindOne(ctx, bson.M{
		"email":                   email,
		"verification_code":       code,
		"verified":                false,
		"verification_expires_at": bson.M{"$gt": now},
	}).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("invalid code")
	}
	if err != nil {
		return nil, err
	}

	_, err = s.col.UpdateOne(ctx,
		bson.M{"_id": u.ID},
		bson.M{"$set": bson.M{"verified": true}, "$unset": bson.M{
			"verification_code":       "",
			"verification_expires_at": "",
		}},
	)
	if err != nil {
		return nil, err
	}

	u.Verified = true
	u.VerificationCode = ""
	return &u, nil
}
