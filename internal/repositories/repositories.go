package repositories

import (
	"context"
	"github.com/Genry72/gophermart/internal/models"
	_ "github.com/lib/pq"
)

// Userser управление пользователями.
type Userser interface {
	// AddUser Добавление пользователя
	AddUser(ctx context.Context, user *models.User) (*models.User, error)
	// GetUserInfo Получение информации по пользователю
	GetUserInfo(ctx context.Context, username string) (*models.User, error)
}

// Orderer работа с заказами.
type Orderer interface {
	// GetOrderByID получение заказа по ID
	GetOrderByID(ctx context.Context, orderID int64) (*models.Order, error)
	// AddOrder запись заказа в базу
	AddOrder(ctx context.Context, orderID, userID int64) (*models.Order, error)
	// GetOrdersByUserID Получение всех заказов пользователя
	GetOrdersByUserID(ctx context.Context, userID int64) ([]*models.Order, error)
}

// Balancer работа с балансом.
type Balancer interface {
	// GetUserBalance Получение баланса пользователя
	GetUserBalance(ctx context.Context, userID int64) (*models.Balance, error)
	// Withdraw Списание средств
	Withdraw(ctx context.Context, withdraw *models.Withdraw) error
	// Withdrawals Получение информации о выводе средств
	Withdrawals(ctx context.Context, userID int64) ([]*models.Withdraw, error)
}

// Accrualer получение информации по заказам.
type Accrualer interface {
	// GetUnprocessedOrders Получение id заказов по которым нужны проверки статусов
	GetUnprocessedOrders(ctx context.Context) ([]int64, error)
	// GetAccrualInfo получение информации по заказам из accrual
	GetAccrualInfo(ctx context.Context, orderIDs []int64) models.ResponseAccruals
	// WriteStatus Запись информации по заказам в базу
	WriteStatus(ctx context.Context, src models.ResponseAccruals) error
}
