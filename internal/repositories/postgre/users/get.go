package users

import (
	"context"
	"fmt"
	"github.com/Genry72/gophermart/internal/models"
)

func (u *UserStorage) GetUserInfo(ctx context.Context, username string) (*models.User, error) {
	query := `
select user_id,
       username,
       password_hash,
       email,
       first_name,
       last_name,
       phone,
       created_at,
       updated_at
from users where username = $1
`
	row := u.conn.QueryRowxContext(ctx, query, username)

	var user models.User

	if err := row.StructScan(&user); err != nil {
		return nil, fmt.Errorf("rows.StructScan: %w", err)
	}

	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("row.Err: %w", err)
	}

	return &user, nil

}
