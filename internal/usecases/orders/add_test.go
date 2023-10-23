package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Genry72/gophermart/internal/logger"
	"github.com/Genry72/gophermart/internal/models"
	"github.com/Genry72/gophermart/internal/models/myerrors"
	mockRepo "github.com/Genry72/gophermart/internal/repositories/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrders_AddOrder(t *testing.T) {

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockOrders := mockRepo.NewMockOrderer(mockCtl)

	zapLogger := logger.NewZapLogger("info")

	uc := &Orders{
		log:  zapLogger,
		repo: mockOrders,
	}

	ctx := context.Background()
	orderID := int64(123)
	userID := int64(321)
	anyErr := fmt.Errorf("anyErr")
	order := &models.Order{
		OrderID: fmt.Sprint(orderID),
		UserID:  userID,
		Status:  "status",
		Accrual: 1,
	}

	type args struct {
		orderID  int64
		userID   int64
		mockFunc func()
	}

	tests := []struct {
		name    string
		args    args
		want    *models.Order
		wantErr error
	}{
		{
			name: "positive",
			args: args{
				orderID: orderID,
				userID:  userID,
				mockFunc: func() {
					mockOrders.EXPECT().GetOrderByID(ctx, orderID).Return(nil, sql.ErrNoRows)
					mockOrders.EXPECT().AddOrder(ctx, orderID, userID).Return(order, nil)
				},
			},
			want:    order,
			wantErr: nil,
		},
		{
			name: "ErrOrderUploadByAnotherUser",
			args: args{
				orderID: orderID,
				userID:  userID,
				mockFunc: func() {
					o := *order
					o.UserID = 456
					mockOrders.EXPECT().GetOrderByID(ctx, orderID).Return(&o, nil)
				},
			},
			want:    nil,
			wantErr: myerrors.ErrOrderUploadByAnotherUser,
		},
		{
			name: "ErrOrderAlreadyUploadByUse",
			args: args{
				orderID: orderID,
				userID:  userID,
				mockFunc: func() {
					mockOrders.EXPECT().GetOrderByID(ctx, orderID).Return(order, nil)
				},
			},
			want:    nil,
			wantErr: myerrors.ErrOrderAlreadyUploadByUser,
		},
		{
			name: "anyErr",
			args: args{
				orderID: orderID,
				userID:  userID,
				mockFunc: func() {
					mockOrders.EXPECT().GetOrderByID(ctx, orderID).Return(nil, anyErr)
				},
			},
			want:    nil,
			wantErr: anyErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.mockFunc != nil {
				tt.args.mockFunc()
			}

			got, err := uc.AddOrder(ctx, tt.args.orderID, tt.args.userID)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.True(t, errors.Is(err, tt.wantErr))
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
