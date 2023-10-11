package repositories

import (
	"context"
	"github.com/Genry72/gophermart/internal/models"
	_ "github.com/lib/pq"
)

type Userser interface {
	AddUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserInfo(ctx context.Context, username string) (*models.User, error)
}

type Orderer interface {
	GetOrderByID(ctx context.Context, orderID int64) (*models.Order, error)
	AddOrder(ctx context.Context, orderID, userID int64) (*models.Order, error)
	GetOrdersByUserID(ctx context.Context, userID int64) ([]*models.Order, error)
}
