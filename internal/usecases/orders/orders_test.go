package orders

import (
	"github.com/Genry72/gophermart/internal/logger"
	"github.com/Genry72/gophermart/internal/repositories/postgre/orders"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"go.uber.org/zap"
	"testing"
)

func TestNewOrders(t *testing.T) {
	zapLogger := logger.NewZapLogger("info")

	db, _, err := sqlxmock.Newx()
	assert.NoError(t, err)

	repo := orders.NewOrderStorage(db, zapLogger)

	type args struct {
		conn *sqlx.DB
		log  *zap.Logger
	}

	tests := []struct {
		name string
		args args
		want *Orders
	}{
		{
			name: "positive",
			args: args{
				conn: db,
				log:  zapLogger,
			},
			want: &Orders{
				log:  zapLogger,
				repo: repo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewOrders(tt.args.conn, tt.args.log), "NewOrders(%v, %v)", tt.args.conn, tt.args.log)
		})
	}
}
