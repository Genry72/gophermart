package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Genry72/gophermart/internal/logger"
	"github.com/Genry72/gophermart/internal/models"
	"github.com/Genry72/gophermart/internal/models/myerrors"
	mockRepo "github.com/Genry72/gophermart/internal/repositories/mocks"
	"github.com/Genry72/gophermart/pkg/cryptor"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUsers_CreateUser(t *testing.T) {

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockUsers := mockRepo.NewMockUserser(mockCtl)

	zapLogger := logger.NewZapLogger("info")

	uc := &Users{
		log:  zapLogger,
		repo: mockUsers,
	}

	login := "login"
	pass := "pass"
	userID := int64(1)
	passHash, err := cryptor.Sha256(pass)
	assert.NoError(t, err)
	updateTime := time.Now()
	ctx := context.Background()
	anyAerr := fmt.Errorf("anyErr")
	wantUser := &models.User{
		UserID:       userID,
		Username:     login,
		PasswordHash: passHash,
		CreatedAt:    updateTime,
		UpdatedAt:    updateTime,
	}

	registerUser := &models.UserRegister{
		Username: login,
		Password: pass,
	}

	type args struct {
		user     *models.UserRegister
		mockFunc func()
	}

	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr error
	}{
		// CreateUser
		{
			name: "CreateUser positive",
			args: args{
				user: registerUser,
				mockFunc: func() {
					mockUsers.EXPECT().GetUserInfo(ctx, registerUser.Username).Return(nil, sql.ErrNoRows)
					u := *wantUser
					u.UserID = 0
					u.CreatedAt = time.Time{}
					u.UpdatedAt = time.Time{}
					mockUsers.EXPECT().AddUser(ctx, &u).Return(wantUser, nil)
				},
			},
			want:    wantUser,
			wantErr: nil,
		},
		{
			name: "CreateUser ErrUserAlreadyExist",
			args: args{
				user: registerUser,
				mockFunc: func() {
					mockUsers.EXPECT().GetUserInfo(ctx, registerUser.Username).Return(nil, nil)
				},
			},
			want:    nil,
			wantErr: myerrors.ErrUserAlreadyExist,
		},
		{
			name: "CreateUser any err",
			args: args{
				user: registerUser,
				mockFunc: func() {
					mockUsers.EXPECT().GetUserInfo(ctx, registerUser.Username).Return(nil, anyAerr)
				},
			},
			want:    nil,
			wantErr: anyAerr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.mockFunc != nil {
				tt.args.mockFunc()
			}

			got, err := uc.CreateUser(ctx, tt.args.user)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.True(t, errors.Is(err, tt.wantErr))
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
