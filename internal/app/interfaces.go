package app

import (
	"context"
	"time"

	"github.com/MAKLUBE/AP1_Final_Project/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStore interface {
	Insert(ctx context.Context, u *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	CreateRestaurant(ctx context.Context, restaurant *models.Restaurant) error
	SetVerification(ctx context.Context, userID primitive.ObjectID, code string, expiresAt time.Time) error
	VerifyByEmailAndCode(ctx context.Context, email string, code string, now time.Time) (*models.User, error)
}

type SessionStore interface {
	Create(ctx context.Context, s *models.Session) error
	FindByToken(ctx context.Context, token string) (*models.Session, error)
	Delete(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context, now time.Time) error
}

type RestaurantStore interface {
	List(ctx context.Context) ([]models.Restaurant, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Restaurant, error)
	Insert(ctx context.Context, restaurant *models.Restaurant) error
	ListByOwner(ctx context.Context, ownerID primitive.ObjectID) ([]models.Restaurant, error)
	Update(ctx context.Context, r *models.Restaurant) error
}

type MenuItemStore interface {
	ListByRestaurant(ctx context.Context, restaurantID primitive.ObjectID) ([]models.MenuItem, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.MenuItem, error)
	Insert(ctx context.Context, m *models.MenuItem) error
	Update(ctx context.Context, m *models.MenuItem) error
}

type ReviewStore interface {
	ListByRestaurant(ctx context.Context, restaurantID primitive.ObjectID) ([]models.Review, error)
	ListAll(ctx context.Context) ([]models.Review, error)
	Insert(ctx context.Context, restaurant *models.Review) error
	SetHidden(ctx context.Context, id primitive.ObjectID, hidden bool, moderatorID primitive.ObjectID, note string) error
}

type OrderStore interface {
	Insert(ctx context.Context, o *models.Order) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Order, error)
	ListByCustomer(ctx context.Context, customerID primitive.ObjectID) ([]models.Order, error)
	ListByRestaurant(ctx context.Context, restaurantID primitive.ObjectID) ([]models.Order, error)
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error
}

type OrderStatusHistoryStore interface {
	Insert(ctx context.Context, h *models.OrderStatusHistory) error
	ListByOrder(ctx context.Context, orderID primitive.ObjectID) ([]models.OrderStatusHistory, error)
}

type HalalVerificationStore interface {
	Insert(ctx context.Context, v *models.HalalVerification) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.HalalVerification, error)
	ListPending(ctx context.Context) ([]models.HalalVerification, error)
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status string, note string) error
}
