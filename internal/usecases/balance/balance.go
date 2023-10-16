package balance

import (
	"context"
	"github.com/Genry72/gophermart/internal/models"
	"github.com/Genry72/gophermart/internal/repositories"
	"github.com/Genry72/gophermart/internal/repositories/postgre"
	"github.com/Genry72/gophermart/internal/repositories/postgre/balances"
	"go.uber.org/zap"
)

type Balances struct {
	log  *zap.Logger
	repo repositories.Balancer
}

func NewBalances(repo *postgre.PGStorage, log *zap.Logger) *Balances {
	return &Balances{
		log:  log,
		repo: balances.NewBalanceStorage(repo.Conn, log),
	}
}

func (u *Balances) GetUserBalance(ctx context.Context, userID int64) (*models.Balance, error) {
	return u.repo.GetUserBalance(ctx, userID)
}

func (u *Balances) Withdraw(ctx context.Context, withdraw *models.Withdraw) error {
	return u.repo.Withdraw(ctx, withdraw)
}

func (u *Balances) Withdrawals(ctx context.Context, userID int64) ([]*models.Withdraw, error) {
	return u.repo.Withdrawals(ctx, userID)
}
