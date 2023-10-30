package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Genry72/gophermart/internal/models"
	"github.com/Genry72/gophermart/internal/models/myerrors"
	"github.com/Genry72/gophermart/pkg/cryptor"
)

func (u *Users) AuthUser(ctx context.Context, username, password string) (*models.User, error) {
	// Проверка существования логина
	user, err := u.repo.GetUserInfo(ctx, username)
	switch {
	case err == nil:
		break
	case errors.Is(err, sql.ErrNoRows):
		return nil, myerrors.ErrUnauthorized
	default:
		return nil, fmt.Errorf("u.repo.GetUserInfo: %w", err)
	}

	pass, err := cryptor.Sha256(password)
	if err != nil {
		return nil, fmt.Errorf("cryptor.Sha256: %w", err)
	}

	if pass != user.PasswordHash {
		return nil, myerrors.ErrUnauthorized
	}

	return user, nil
}
