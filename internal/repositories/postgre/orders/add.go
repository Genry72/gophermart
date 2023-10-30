package orders

import (
	"context"
	"fmt"
	"github.com/Genry72/gophermart/internal/models"
)

func (o *OrderStorage) AddOrder(ctx context.Context, orderID, userID int64) (*models.Order, error) {
	query := `
INSERT INTO orders (order_id, user_id, status, accrual, created_at, updated_at)
VALUES ($1, $2, $3, 0, DEFAULT, DEFAULT)
returning 
    user_id,
    order_id,
    user_id,
    status,
    accrual,
    created_at,
    updated_at;
`

	var result models.Order

	row := o.conn.QueryRowxContext(ctx, query, orderID, userID, models.OrderStatusNew)

	if err := row.StructScan(&result); err != nil {
		return nil, fmt.Errorf("row.StructScan: %w", err)
	}

	return &result, nil
}
