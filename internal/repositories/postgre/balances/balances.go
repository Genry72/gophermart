package balances

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Genry72/gophermart/internal/models"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type BalanceStorage struct {
	conn *sqlx.DB
	log  *zap.Logger
}

func NewBalanceStorage(conn *sqlx.DB, log *zap.Logger) *BalanceStorage {
	return &BalanceStorage{conn: conn, log: log}
}

func (u *BalanceStorage) GetUserBalance(ctx context.Context, userID int64) (*models.Balance, error) {
	query := `
select
sum(accrual) - COALESCE(sum(w.points), 0) as current,
COALESCE(sum(w.points), 0) as withdrawn
from orders o
          left join withdraw w on o.user_id = w.user_id
where o.user_id = $1
group by o.user_id, w.date
order by w.date
`
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
