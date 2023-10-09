package jwtToken

import (
	"github.com/Genry72/gophermart/internal/models"
	"testing"
	"time"
)

func TestGetToken(t *testing.T) {
	type args struct {
		tokenKey string
		lifetime time.Duration
		userName string
		id       int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "positive",
			args: args{
				tokenKey: "1111",
				lifetime: 10 * time.Second,
				userName: "userName",
				id:       5,
			},

			wantErr: false,
		},
		{
			name: "empty pass",
			args: args{
				tokenKey: "",
				lifetime: 10 * time.Second,
				userName: "userName",
				id:       5,
			},

			wantErr: false,
		},
		{
			name: "expired token",
			args: args{
				tokenKey: "",
				//lifetime: 10 * time.Second,
				userName: "userName",
				id:       5,
			},

			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := NewJwtToken(tt.args.tokenKey, tt.args.lifetime)
			u := &models.User{
				UserID:   tt.args.id,
				Username: tt.args.userName,
			}

			token, err := j.GetToken(u)
			if err != nil {
				t.Errorf("GetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotUserID, gotUserName, err := j.ValidateAndParseToken(token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAndParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotUserID != tt.args.id {
				if err == nil {
					t.Errorf("ValidateAndParseToken.UserID got = %v, want %v", gotUserID, tt.args.id)
				}
			}

			if gotUserName != tt.args.userName {
				if err == nil {
					t.Errorf("ValidateAndParseToken.UserID got = %v, want %v", gotUserName, tt.args.userName)
				}
			}
		})
	}
}
