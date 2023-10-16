package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Genry72/gophermart/internal/models"
)

func (u *UserStorage) GetUserBalance(ctx context.Context, userID int64) (*models.Balance, error) {
	query := `
select
sum(accrual) - COALESCE(sum(w.points), 0) as current,
COALESCE(sum(w.points), 0) as withdrawn
from orders o
          left join withdraw w on o.order_id = w.order_id
where o.user_id = $1
group by o.user_id
`
	fmt.Println(query)
	result := &models.Balance{}

	row := u.conn.QueryRowxContext(ctx, query, userID)

	if err := row.StructScan(result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			result.Current = 0
			result.Withdrawn = 0

			return result, nil
		}

		return nil, fmt.Errorf("row.StructScan:%w", err)
	}

	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("row.Err: %w", err)
	}

	return result, nil
}
