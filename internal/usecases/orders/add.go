package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Genry72/gophermart/internal/models"
	"github.com/Genry72/gophermart/internal/models/myerrors"
)

func (o *Orders) AddOrder(ctx context.Context, orderID int64, userID int64) (*models.Order, error) {
	// Проверка существования ордера с данным id
	existOrder, err := o.repo.GetOrderByID(ctx, orderID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("o.repo.GetOrderByID: %w", err)
	}

	if existOrder != nil && existOrder.UserID != userID {
		return nil, myerrors.ErrOrderUploadByAnotherUser
	}

	if existOrder != nil && existOrder.UserID == userID {
		return nil, myerrors.ErrOrderAlreadyUploadByUser
	}

	return o.repo.AddOrder(ctx, orderID, userID)
}
