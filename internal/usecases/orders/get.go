package orders

import (
	"context"
	"github.com/Genry72/gophermart/internal/models"
)

func (o *Orders) GetOrdersByUserID(ctx context.Context, userID int64) ([]*models.Order, error) {

	return o.repo.GetOrdersByUserID(ctx, userID)
}
