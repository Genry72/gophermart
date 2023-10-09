package repositories

import (
	"context"
	"github.com/Genry72/gophermart/internal/models"
	_ "github.com/lib/pq"
)

type Userser interface {
	AddUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserInfo(ctx context.Context, username string) (*models.User, error)
}
