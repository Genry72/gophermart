package orders

import (
	"context"
	"errors"
	"fmt"
	"github.com/Genry72/gophermart/internal/logger"
	"github.com/Genry72/gophermart/internal/models"
	mockRepo "github.com/Genry72/gophermart/internal/repositories/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrders_GetOrdersByUserID(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockOrders := mockRepo.NewMockOrderer(mockCtl)

	zapLogger := logger.NewZapLogger("info")

	uc := &Orders{
		log:  zapLogger,
		repo: mockOrders,
	}

	ctx := context.Background()
	userID := int64(321)
	orderID := int64(123)
	myerr := fmt.Errorf("myerr")

	orders := []*models.Order{{
		OrderID: fmt.Sprint(orderID),
		UserID:  userID,
	}}

	type args struct {
		mockFunc func()
	}

	tests := []struct {
		name    string
		args    args
		want    []*models.Order
		wantErr error
	}{
		{
			name: "positive",
			args: args{func() {
				mockOrders.EXPECT().GetOrdersByUserID(ctx, userID).Return(orders, nil)
			}},
			want:    orders,
			wantErr: nil,
		},
		{
			name: "negative",
			args: args{func() {
				mockOrders.EXPECT().GetOrdersByUserID(ctx, userID).Return(nil, myerr)
			}},
			want:    nil,
			wantErr: myerr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.mockFunc != nil {
				tt.args.mockFunc()
			}
			got, err := uc.GetOrdersByUserID(ctx, userID)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.True(t, errors.Is(err, tt.wantErr))
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
