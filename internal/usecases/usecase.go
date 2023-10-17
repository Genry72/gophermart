package usecases

import (
	"github.com/Genry72/gophermart/internal/usecases/balance"
	"github.com/Genry72/gophermart/internal/usecases/orders"
	"github.com/Genry72/gophermart/internal/usecases/users"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Usecase struct {
	Users    *users.Users
	Orders   *orders.Orders
	Balances *balance.Balances
}

func NewUsecase(conn *sqlx.DB, log *zap.Logger) *Usecase {
	return &Usecase{
		Users:    users.NewUsers(conn, log),
		Orders:   orders.NewOrders(conn, log),
		Balances: balance.NewBalances(conn, log),
	}
}
