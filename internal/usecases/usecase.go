package usecases

import (
	"github.com/Genry72/gophermart/internal/repositories/postgre"
	"github.com/Genry72/gophermart/internal/usecases/balance"
	"github.com/Genry72/gophermart/internal/usecases/orders"
	"github.com/Genry72/gophermart/internal/usecases/users"
	"go.uber.org/zap"
)

type Usecase struct {
	Users    *users.Users
	Orders   *orders.Orders
	Balances *balance.Balances
}

func NewUsecase(repo *postgre.PGStorage, log *zap.Logger) *Usecase {
	return &Usecase{
		Users:    users.NewUsers(repo, log),
		Orders:   orders.NewOrders(repo, log),
		Balances: balance.NewBalances(repo, log),
	}
}
