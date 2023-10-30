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

func (u *Users) CreateUser(ctx context.Context, user *models.UserRegister) (*models.User, error) {
	// Проверка существования логина
	_, err := u.repo.GetUserInfo(ctx, user.Username)
	switch {
	case err == nil:
		return nil, myerrors.ErrUserAlreadyExist
	case errors.Is(err, sql.ErrNoRows):
		break
	default:
		return nil, fmt.Errorf("u.repo.GetUserInfo: %w", err)
	}

	pass, err := cryptor.Sha256(user.Password)
	if err != nil {
		return nil, fmt.Errorf("cryptor.Sha256: %w", err)
	}
	newUser := &models.User{
		Username:     user.Username,
		PasswordHash: pass,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Phone:        user.Phone,
	}

	return u.repo.AddUser(ctx, newUser)
}
