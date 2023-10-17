package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/Genry72/gophermart/internal/handlers/jwtauth/jwttoken"
	"github.com/Genry72/gophermart/internal/logger"
	"github.com/Genry72/gophermart/internal/models"
	"github.com/Genry72/gophermart/internal/usecases"
	"github.com/Genry72/gophermart/pkg/cryptor"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	login := "login1"

	password := "pass"

	hashPass, err := cryptor.Sha256(password)

	assert.NoError(t, err)

	zapLogger := logger.NewZapLogger("info")

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer func() {
		_ = db.Close()
	}()

	testCases := []struct {
		// input
		name string
		body *models.UserRegister // Отправляемое боди

		// mock func
		mockDB func()

		// expected
		expectedStatusCode     int // Ожидаемый код ответа
		expectedSuccessMessage *models.User
		expectedErr            error // Ожидаемая ошибка
	}{
		{
			name: "#1",
			body: &models.UserRegister{
				Username: login,
				Password: password,
			},
			mockDB: func() {
				mock.ExpectQuery("select").WillReturnError(sql.ErrNoRows)

				rows := sqlxmock.NewRows([]string{"user_id", "username", "password_hash"}).AddRow(1, login, hashPass)

				mock.ExpectQuery("INSERT INTO users").WillReturnRows(rows)
			},
			expectedStatusCode: http.StatusOK,
			expectedSuccessMessage: &models.User{
				Username:     login,
				PasswordHash: hashPass,
			},
		},
	}

	g := gin.New()

	use(g, zapLogger)

	tokenKey := "abc"

	jwtLifeTime := 24 * time.Hour

	for _, testCase := range testCases {
		testCase.mockDB()

		handlers := Handler{
			useCases:  usecases.NewUsecase(db, zapLogger),
			log:       zapLogger,
			ginEngine: g,
			authToken: jwttoken.NewJwtToken(tokenKey, jwtLifeTime),
		}

		handlers.initRoutes()

		recorder := httptest.NewRecorder()

		var body io.Reader
		if testCase.body != nil {
			b, err := json.Marshal(testCase.body)
			assert.NoError(t, err)
			body = bytes.NewReader(b)
		}

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/user/register", body)
		require.NoError(t, err)

		g.ServeHTTP(recorder, req)
		if status := recorder.Code; status != testCase.expectedStatusCode {
			t.Errorf("status codes differ: expected %d, got %d err %v", testCase.expectedStatusCode, status, recorder.Body.String())
		}

		respBody := &models.User{}

		err = json.Unmarshal(recorder.Body.Bytes(), respBody)
		assert.NoError(t, err)

		assert.Equal(t, testCase.body.Username, respBody.Username)
		// todo проверка заголовка авторизации
	}
}
