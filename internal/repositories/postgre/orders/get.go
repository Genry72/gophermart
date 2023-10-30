package orders

import (
	"context"
	"fmt"
	"github.com/Genry72/gophermart/internal/models"
)

func (o *OrderStorage) GetOrderByID(ctx context.Context, orderID int64) (*models.Order, error) {
	query := `
select order_id, user_id, status, accrual, created_at, updated_at
from orders where order_id = $1
`

	row := o.conn.QueryRowxContext(ctx, query, orderID)

	var result models.Order

	if err := row.StructScan(&result); err != nil {
		return nil, fmt.Errorf("row.StructScan: %w", err)
	}

	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("row.Err: %w", err)
	}

	return &result, nil
}

func (o *OrderStorage) GetOrdersByUserID(ctx context.Context, userID int64) ([]*models.Order, error) {
	query := `
select order_id, user_id, status, accrual, created_at, updated_at
from orders where user_id = $1
order by created_at
`

	result := make([]*models.Order, 0)

	if err := o.conn.SelectContext(ctx, &result, query, userID); err != nil {
		return nil, fmt.Errorf("o.conn.SelectContext: %w", err)
	}

	return result, nil
}
