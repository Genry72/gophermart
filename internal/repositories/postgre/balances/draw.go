package balances

import (
	"context"
	"fmt"
	"github.com/Genry72/gophermart/internal/models"
)

// Withdraw Списание средств
func (u *BalanceStorage) Withdraw(ctx context.Context, withdraw *models.Withdraw) error {
	query := `
INSERT INTO withdraw (user_id, order_id, points, date)
VALUES (:user_id, :order_id, :points, now())
returning user_id, order_id, points;
`

	_, err := u.conn.NamedExecContext(ctx, query, withdraw)
	if err != nil {
		return fmt.Errorf("u.conn.NamedExecContext:%w", err)
	}

	return nil
}

// Withdrawals Получение информации о выводе средств
func (u *BalanceStorage) Withdrawals(ctx context.Context, userID int64) ([]*models.Withdraw, error) {
	query := `
select order_id, points, date
from withdraw
where user_id = $1;
`

	result := make([]*models.Withdraw, 0)

	if err := u.conn.SelectContext(ctx, &result, query, userID); err != nil {
		return nil, fmt.Errorf("u.conn.SelectContext:%w", err)
	}

	return result, nil
}
