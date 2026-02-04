package app

import (
	"context"

	"github.com/MAKLUBE/AP1_Final_Project/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStore interface {
	Insert(ctx context.Context, u *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	CreateRestaurant(ctx context.Context, restaurant *models.Restaurant) error
}

type SessionStore interface {
	Create(ctx context.Context, s *models.Session) error
	FindByToken(ctx context.Context, token string) (*models.Session, error)
	Delete(ctx context.Context, token string) error
}

type RestaurantStore interface {
	List(ctx context.Context) ([]models.Restaurant, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Restaurant, error)
	Insert(ctx context.Context, restaurant *models.Restaurant) error
}

type MenuItemStore interface {
	ListByRestaurant(ctx context.Context, restaurantID primitive.ObjectID) ([]models.MenuItem, error)
}

type ReviewStore interface {
	ListByRestaurant(ctx context.Context, restaurantID primitive.ObjectID) ([]models.Review, error)
	Insert(ctx context.Context, r *models.Review) error
}
