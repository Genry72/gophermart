package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Genry72/gophermart/internal/models"
	"go.uber.org/zap"
)

func (u *UserStorage) AddUser(ctx context.Context, user *models.User) (*models.User, error) {

	tx, err := u.conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("r.conn.BeginTx: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			u.log.Error("tx.Rollback", zap.Error(err))
		}
	}()

	// Добавляем пользователя
	queryAddUser := `
INSERT INTO users (username,
                   password_hash,
                   email,
                   first_name,
                   last_name,
                   phone,
                   created_at,
                   updated_at)
VALUES (:username,
        :password_hash,
        :email,
        :first_name,
        :last_name,
        :phone,
        DEFAULT,
        DEFAULT)
returning user_id,
    username,
    password_hash,
    email,
    first_name,
    last_name,
    phone,
    created_at,
    updated_at;
`

	rows, err := tx.NamedQuery(queryAddUser, user)
	if err != nil {
		return nil, fmt.Errorf("u.conn.NamedQueryContext: %w", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			u.log.Error("rows.Close()", zap.Error(err))
			return
		}
	}()

	var addedUser *models.User

	for rows.Next() {
		addedUser = &models.User{}
		if err := rows.StructScan(addedUser); err != nil {
			return nil, fmt.Errorf("rows.StructScan: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	// Создаем для пользователя запись с его балансом
	queryAddBalance := `
INSERT INTO user_balance (user_id,
                          accruals,
                          drawal,
                          current_balance,
                          last_update)
VALUES ($1,
        DEFAULT,
        DEFAULT,
        DEFAULT,
        DEFAULT);
`

	if _, err := tx.ExecContext(ctx, queryAddBalance, addedUser.UserID); err != nil {
		return nil, fmt.Errorf("add balance: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("tx.Commit: %w", err)
	}

	return addedUser, nil
}
