package usecases

import (
	"context"
	"github.com/Genry72/gophermart/internal/models"
	"github.com/Genry72/gophermart/internal/usecases/balance"
	"github.com/Genry72/gophermart/internal/usecases/orders"
	"github.com/Genry72/gophermart/internal/usecases/users"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Userser interface {
	AuthUser(ctx context.Context, username, password string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.UserRegister) (*models.User, error)
}

type Orderser interface {
	AddOrder(ctx context.Context, orderID int64, userID int64) (*models.Order, error)
	GetOrdersByUserID(ctx context.Context, userID int64) ([]*models.Order, error)
}

type Balancer interface {
	GetUserBalance(ctx context.Context, userID int64) (*models.Balance, error)
	Withdraw(ctx context.Context, withdraw *models.Withdraw) error
	Withdrawals(ctx context.Context, userID int64) ([]*models.Withdraw, error)
}

type Usecase struct {
	Users    Userser
	Orders   Orderser
	Balances Balancer
}

func NewUsecase(conn *sqlx.DB, log *zap.Logger) *Usecase {
	return &Usecase{
		Users:    users.NewUsers(conn, log),
		Orders:   orders.NewOrders(conn, log),
		Balances: balance.NewBalances(conn, log),
	}
}
