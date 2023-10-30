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

func TestUsers_AuthUser(t *testing.T) {

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

	wantUser := &models.User{
		UserID:       userID,
		Username:     login,
		PasswordHash: passHash,
		CreatedAt:    updateTime,
		UpdatedAt:    updateTime,
	}

	ctx := context.Background()

	anyErr := fmt.Errorf("anyErr")

	type args struct {
		username string
		password string
		mockFunc func()
	}

	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr error
	}{
		{
			name: "positive",
			args: args{
				username: login,
				password: pass,
				mockFunc: func() {
					mockUsers.EXPECT().GetUserInfo(ctx, login).Return(wantUser, nil)
				},
			},
			want:    wantUser,
			wantErr: nil,
		},
		{
			name: "bad hash",
			args: args{
				username: login,
				password: pass,
				mockFunc: func() {
					u := *wantUser
					u.PasswordHash = "bad"
					mockUsers.EXPECT().GetUserInfo(ctx, login).Return(&u, nil)
				},
			},
			want:    nil,
			wantErr: myerrors.ErrUnauthorized,
		},
		{
			name: "ErrUnauthorized",
			args: args{
				username: login,
				password: pass,
				mockFunc: func() {
					mockUsers.EXPECT().GetUserInfo(ctx, login).Return(nil, sql.ErrNoRows)
				},
			},
			want:    nil,
			wantErr: myerrors.ErrUnauthorized,
		},
		{
			name: "any err",
			args: args{
				username: login,
				password: pass,
				mockFunc: func() {
					mockUsers.EXPECT().GetUserInfo(ctx, login).Return(nil, anyErr)
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
			got, err := uc.AuthUser(ctx, tt.args.username, tt.args.password)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.True(t, errors.Is(err, tt.wantErr))
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
