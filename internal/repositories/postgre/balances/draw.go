package balances

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Genry72/gophermart/internal/models"
	"github.com/Genry72/gophermart/internal/models/myerrors"
	"go.uber.org/zap"
)

// Withdraw Списание средств
func (u *BalanceStorage) Withdraw(ctx context.Context, withdraw *models.Withdraw) error {

	tx, err := u.conn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("r.conn.BeginTx: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			u.log.Error("tx.Rollback", zap.Error(err))
		}
	}()

	// Получаем количество записей по списаниям пользователя
	var countDraw int64

	countDrawRow := tx.QueryRowxContext(ctx, "select count (*) from withdraw WHERE user_id = $1;", withdraw.UserID)

	if err := countDrawRow.Scan(&countDraw); err != nil {
		return fmt.Errorf("countDrawRow.Scan: %w", err)
	}

	// Если списаний не было, то блокируем всю таблицу. Операция дорогая, но иначе можем получить повторное списание
	// при одновременных запросах.
	if countDraw == 0 {
		fmt.Println("lock")
		if _, err := tx.ExecContext(ctx, "LOCK TABLE withdraw IN SHARE MODE;"); err != nil {
			return fmt.Errorf("LOCK TABLE withdraw: %w", err)
		}
	} else { // Записи есть, блокируем только их
		lockDrawQuery := "SELECT * FROM withdraw  WHERE user_id = $1 FOR UPDATE;"
		if _, err := tx.ExecContext(ctx, lockDrawQuery, withdraw.UserID); err != nil {
			return fmt.Errorf("lockDtawQuery: %w", err)
		}
	}

	// Выполняем списание, возвращаем результат операции
	resultQuery := `
WITH _ AS (
    INSERT INTO withdraw (user_id, order_id, points, date)
    VALUES ($1, $2, $3, NOW())
    RETURNING user_id
)
SELECT
    SUM(accrual) - COALESCE(SUM(w.points), 0) AS current,
    COALESCE(SUM(w.points), 0) AS withdrawn
FROM
    orders o
LEFT JOIN
    withdraw w ON o.user_id = w.user_id
WHERE
    o.user_id = $1
GROUP BY
    o.user_id;
`

	resultRow := tx.QueryRowxContext(ctx, resultQuery, withdraw.UserID, withdraw.Order, withdraw.Points)

	var result models.Balance

	if err := resultRow.StructScan(&result); err != nil {
		return fmt.Errorf("resultRow.StructScan: %w", err)
	}

	if result.Current < 0 {
		return myerrors.ErrNoMoney
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
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
