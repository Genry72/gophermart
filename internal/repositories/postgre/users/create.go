package users

import (
	"context"
	"fmt"
	"github.com/Genry72/gophermart/internal/models"
	"go.uber.org/zap"
)

func (u *UserStorage) AddUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
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

	rows, err := u.conn.NamedQueryContext(ctx, query, user)
	if err != nil {
		return nil, fmt.Errorf("u.conn.NamedQueryContext: %w", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			u.log.Error("rows.Close()", zap.Error(err))
			return
		}
	}()

	var result *models.User

	for rows.Next() {
		result = &models.User{}
		if err := rows.StructScan(result); err != nil {
			return nil, fmt.Errorf("rows.StructScan: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return result, nil
}
