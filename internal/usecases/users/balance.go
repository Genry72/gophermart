package users

import (
	"context"
	"github.com/Genry72/gophermart/internal/models"
)

func (u *Users) GetUserBalance(ctx context.Context, userID int64) (*models.Balance, error) {
	return u.repo.GetUserBalance(ctx, userID)
}
