package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	mockAuth "github.com/Genry72/gophermart/internal/handlers/jwtauth/mocks"
	"github.com/Genry72/gophermart/internal/logger"
	"github.com/Genry72/gophermart/internal/models"
	"github.com/Genry72/gophermart/internal/models/myerrors"
	"github.com/Genry72/gophermart/internal/usecases"
	mockUc "github.com/Genry72/gophermart/internal/usecases/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

//func TestRegister(t *testing.T) {
//	login := "login1"
//
//	password := "pass"
//
//	hashPass, err := cryptor.Sha256(password)
//
//	assert.NoError(t, err)
//
//	zapLogger := logger.NewZapLogger("info")
//
//	db, mock, err := sqlxmock.Newx()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//
//	defer func() {
//		_ = db.Close()
//	}()
//
//	testCases := []struct {
//		// input
//		name string
//		body *models.UserRegister // Отправляемое боди
//
//		// mockFunc func
//		mockDB func()
//
//		// expected
//		expectedStatusCode     int // Ожидаемый код ответа
//		expectedSuccessMessage *models.User
//		expectedErr            error // Ожидаемая ошибка
//	}{
//		{
//			name: "#1",
//			body: &models.UserRegister{
//				Username: login,
//				Password: password,
//			},
//			mockDB: func() {
//				mock.ExpectQuery("select").WillReturnError(sql.ErrNoRows)
//
//				rows := sqlxmock.NewRows([]string{"user_id", "username", "password_hash"}).AddRow(1, login, hashPass)
//
//				mock.ExpectQuery("INSERT INTO users").WillReturnRows(rows)
//			},
//			expectedStatusCode: http.StatusOK,
//			expectedSuccessMessage: &models.User{
//				Username:     login,
//				PasswordHash: hashPass,
//			},
//		},
//	}
//
//	g := gin.New()
//
//	useMidlWare(g, zapLogger)
//
//	tokenKey := "abc"
//
//	jwtLifeTime := 24 * time.Hour
//
//	for _, testCase := range testCases {
//		testCase.mockDB()
//
//		handlers := Handler{
//			useCases:  usecases.NewUsecase(db, zapLogger),
//			log:       zapLogger,
//			ginEngine: g,
//			authToken: jwttoken.NewJwtToken(tokenKey, jwtLifeTime),
//		}
//
//		handlers.initRoutes()
//
//		recorder := httptest.NewRecorder()
//
//		var body io.Reader
//		if testCase.body != nil {
//			b, err := json.Marshal(testCase.body)
//			assert.NoError(t, err)
//			body = bytes.NewReader(b)
//		}
//
//		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/user/register", body)
//		require.NoError(t, err)
//
//		g.ServeHTTP(recorder, req)
//		if status := recorder.Code; status != testCase.expectedStatusCode {
//			t.Errorf("status codes differ: expected %d, got %d err %v", testCase.expectedStatusCode, status, recorder.Body.String())
//		}
//
//		respBody := &models.User{}
//
//		err = json.Unmarshal(recorder.Body.Bytes(), respBody)
//		assert.NoError(t, err)
//
//		assert.Equal(t, testCase.body.Username, respBody.Username)
//	}
//}

var t1 = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

func TestHandlers(t *testing.T) {
	// Мокаем юзеркейсы и работу с токеном
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockUsers := mockUc.NewMockUserser(mockCtl)
	mockOrders := mockUc.NewMockOrderser(mockCtl)
	mockBalance := mockUc.NewMockBalancer(mockCtl)
	mockAuthToken := mockAuth.NewMockAuther(mockCtl)
	uc := &usecases.Usecase{
		Users:    mockUsers,
		Orders:   mockOrders,
		Balances: mockBalance,
	}

	login := "login1"

	password := "pass"

	expectedUser := &models.User{
		Username:  login,
		UserID:    1,
		CreatedAt: t1,
	}

	expectedToken := "token"

	luhnorderID := "79927398713"

	expectedOrder := &models.Order{
		OrderID: luhnorderID,
		//UserID:  1, // userID не возвращаем в ответе
	}

	anyErr := fmt.Errorf("anyErr")

	zapLogger := logger.NewZapLogger("info")

	// Создаем роутер и устанавливаем мидлварю
	g := gin.New()

	useMidlWare(g, zapLogger)

	handlers := Handler{
		useCases:  uc,
		log:       zapLogger,
		ginEngine: g,
		authToken: mockAuthToken,
	}

	handlers.initRoutes()

	type args struct {
		url         string // url запроса
		method      string
		requestBody any               // Отправляемое боди (структура, либо строка)
		headers     map[string]string // Добавляемые заголовки в запрос
		mockFunc    func()
	}

	tests := []struct {
		name string
		args args
		// expected
		expectedStatusCode int               // Ожидаемый код ответа
		expectedBody       any               // Ожидаемое тело ответа
		expectedHeaders    map[string]string // Ожидаемые заголовки в ответе
		expectedErr        error             // Ожидаемая ошибка в теле ответа
		// Функция для парсинга тела ответа
		parseResponseBody func(b []byte) (interface{}, error)
	}{
		// users
		{
			name: "Register positive",
			args: args{
				url:    "/api/user/register",
				method: http.MethodPost,
				requestBody: models.UserRegister{
					Username: login,
					Password: password,
				},
				mockFunc: func() {
					mockUsers.EXPECT().CreateUser(context.Background(), &models.UserRegister{
						Username: login,
						Password: password,
					}).Return(expectedUser, nil)
					mockAuthToken.EXPECT().GetToken(expectedUser).Return(expectedToken, nil)
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       expectedUser,
			expectedHeaders:    map[string]string{"Authorization": "Bearer " + expectedToken},
			parseResponseBody: func(b []byte) (interface{}, error) {
				respBody := &models.User{}
				err := json.Unmarshal(b, respBody)
				if err != nil {
					return nil, fmt.Errorf("json.Unmarshal: %w %s", err, string(b))
				}
				return respBody, nil
			},
		},
		{
			name: "Register bad body format",
			args: args{
				url:         "/api/user/register",
				method:      http.MethodPost,
				requestBody: "{user",
			},

			expectedStatusCode: http.StatusBadRequest,

			expectedErr: ErrBadBody,
		},
		{
			name: "Register err GetToken",
			args: args{
				url:    "/api/user/register",
				method: http.MethodPost,
				requestBody: models.UserRegister{
					Username: login,
					Password: password,
				},
				mockFunc: func() {
					mockUsers.EXPECT().CreateUser(context.Background(), &models.UserRegister{
						Username: login,
						Password: password,
					}).Return(expectedUser, nil)
					mockAuthToken.EXPECT().GetToken(expectedUser).Return("", fmt.Errorf("myErr"))
				},
			},

			expectedStatusCode: http.StatusInternalServerError,

			expectedErr: fmt.Errorf("myErr"),
		},

		{
			name: "Register err CreateUser",
			args: args{
				url:    "/api/user/register",
				method: http.MethodPost,
				requestBody: models.UserRegister{
					Username: login,
					Password: password,
				},
				mockFunc: func() {
					mockUsers.EXPECT().CreateUser(context.Background(), &models.UserRegister{
						Username: login,
						Password: password,
					}).Return(nil, myerrors.ErrUserAlreadyExist)
				},
			},

			expectedStatusCode: http.StatusConflict,

			expectedErr: myerrors.ErrUserAlreadyExist,
		},

		{
			name: "Auth user positive",
			args: args{
				url:    "/api/user/login",
				method: http.MethodPost,
				requestBody: models.UserRegister{
					Username: login,
					Password: password,
				},
				mockFunc: func() {
					mockUsers.EXPECT().AuthUser(context.Background(), login, password).Return(expectedUser, nil)
					mockAuthToken.EXPECT().GetToken(expectedUser).Return(expectedToken, nil)
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       successAuthUser,
			expectedHeaders:    map[string]string{"Authorization": "Bearer " + expectedToken},
			parseResponseBody: func(b []byte) (interface{}, error) {
				return string(b), nil
			},
		},
		{
			name: "Auth user bad body format",
			args: args{
				url:         "/api/user/login",
				method:      http.MethodPost,
				requestBody: "{user",
			},

			expectedStatusCode: http.StatusBadRequest,

			expectedErr: ErrBadBody,
		},
		{
			name: "Auth err from uc ErrUnauthorized",
			args: args{
				url:    "/api/user/login",
				method: http.MethodPost,
				requestBody: models.UserRegister{
					Username: login,
					Password: password,
				},
				mockFunc: func() {
					mockUsers.EXPECT().AuthUser(context.Background(), login, password).Return(nil, myerrors.ErrUnauthorized)
				},
			},

			expectedStatusCode: http.StatusUnauthorized,

			expectedErr: myerrors.ErrUnauthorized,
		},
		{
			name: "Auth err GetToken",
			args: args{
				url:    "/api/user/login",
				method: http.MethodPost,
				requestBody: models.UserRegister{
					Username: login,
					Password: password,
				},
				mockFunc: func() {
					mockUsers.EXPECT().AuthUser(context.Background(), login, password).Return(expectedUser, nil)
					mockAuthToken.EXPECT().GetToken(expectedUser).Return("", fmt.Errorf("myErr"))
				},
			},

			expectedStatusCode: http.StatusInternalServerError,

			expectedErr: fmt.Errorf("myErr"),
		},

		// orders
		{
			name: "uploadOrder positive",
			args: args{
				url:         "/api/user/orders",
				method:      http.MethodPost,
				requestBody: luhnorderID,
				headers:     map[string]string{"Authorization": "Bearer " + expectedToken},
				mockFunc: func() {
					mockAuthToken.EXPECT().ValidateAndParseToken(expectedToken).
						Return(expectedUser.UserID, expectedUser.Username, nil)

					mockOrders.EXPECT().AddOrder(gomock.Any(), int64(79927398713), expectedUser.UserID).
						Return(expectedOrder, nil)
				},
			},
			expectedStatusCode: http.StatusAccepted,
			expectedBody:       expectedOrder,
			parseResponseBody: func(b []byte) (interface{}, error) {
				respBody := &models.Order{}
				err := json.Unmarshal(b, respBody)
				if err != nil {
					return nil, fmt.Errorf("json.Unmarshal: %w %s", err, string(b))
				}
				return respBody, nil
			},
		},
		{
			name: "uploadOrder err AlreadyUploadByUser",
			args: args{
				url:         "/api/user/orders",
				method:      http.MethodPost,
				requestBody: luhnorderID,
				headers:     map[string]string{"Authorization": "Bearer " + expectedToken},
				mockFunc: func() {
					mockAuthToken.EXPECT().ValidateAndParseToken(expectedToken).
						Return(expectedUser.UserID, expectedUser.Username, nil)
					mockOrders.EXPECT().AddOrder(gomock.Any(), int64(79927398713), expectedUser.UserID).
						Return(nil, myerrors.ErrOrderUploadByAnotherUser)
				},
			},
			expectedStatusCode: http.StatusConflict,
			expectedErr:        myerrors.ErrOrderUploadByAnotherUser,
		},
		{
			name: "uploadOrder OrderAlreadyUploadByUser",
			args: args{
				url:         "/api/user/orders",
				method:      http.MethodPost,
				requestBody: luhnorderID,
				headers:     map[string]string{"Authorization": "Bearer " + expectedToken},
				mockFunc: func() {
					mockAuthToken.EXPECT().ValidateAndParseToken(expectedToken).
						Return(expectedUser.UserID, expectedUser.Username, nil)
					mockOrders.EXPECT().AddOrder(gomock.Any(), int64(79927398713), expectedUser.UserID).
						Return(nil, myerrors.ErrOrderAlreadyUploadByUser)
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedErr:        myerrors.ErrOrderAlreadyUploadByUser,
		},
		{
			name: "uploadOrder err no body",
			args: args{
				url:         "/api/user/orders",
				method:      http.MethodPost,
				requestBody: nil,
				headers:     map[string]string{"Authorization": "Bearer " + expectedToken},
				mockFunc: func() {
					mockAuthToken.EXPECT().ValidateAndParseToken(expectedToken).
						Return(expectedUser.UserID, expectedUser.Username, nil)
				},
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedErr:        myerrors.ErrBadFormatOrder,
		},
		{
			name: "uploadOrder orderID not valid",
			args: args{
				url:         "/api/user/orders",
				method:      http.MethodPost,
				requestBody: "123",
				headers:     map[string]string{"Authorization": "Bearer " + expectedToken},
				mockFunc: func() {
					mockAuthToken.EXPECT().ValidateAndParseToken(expectedToken).
						Return(expectedUser.UserID, expectedUser.Username, nil)
				},
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedErr:        myerrors.ErrBadFormatOrder,
		},
		{
			name: "uploadOrder add order error",
			args: args{
				url:         "/api/user/orders",
				method:      http.MethodPost,
				requestBody: luhnorderID,
				headers:     map[string]string{"Authorization": "Bearer " + expectedToken},
				mockFunc: func() {
					mockAuthToken.EXPECT().ValidateAndParseToken(expectedToken).
						Return(expectedUser.UserID, expectedUser.Username, nil)
					mockOrders.EXPECT().AddOrder(gomock.Any(), int64(79927398713), expectedUser.UserID).
						Return(nil, sql.ErrNoRows)
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        sql.ErrNoRows,
		},
		{
			name: "get orders positive",
			args: args{
				url:     "/api/user/orders",
				method:  http.MethodGet,
				headers: map[string]string{"Authorization": "Bearer " + expectedToken},
				mockFunc: func() {
					mockAuthToken.EXPECT().ValidateAndParseToken(expectedToken).
						Return(expectedUser.UserID, expectedUser.Username, nil)
					mockOrders.EXPECT().GetOrdersByUserID(gomock.Any(), expectedUser.UserID).
						Return([]*models.Order{expectedOrder}, nil)
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       []*models.Order{expectedOrder},
			parseResponseBody: func(b []byte) (interface{}, error) {
				respBody := make([]*models.Order, 0)
				err := json.Unmarshal(b, &respBody)
				if err != nil {
					return nil, fmt.Errorf("json.Unmarshal: %w %s", err, string(b))
				}
				return respBody, nil
			},
		},
		{
			name: "get orders err bd",
			args: args{
				url:     "/api/user/orders",
				method:  http.MethodGet,
				headers: map[string]string{"Authorization": "Bearer " + expectedToken},
				mockFunc: func() {
					mockAuthToken.EXPECT().ValidateAndParseToken(expectedToken).
						Return(expectedUser.UserID, expectedUser.Username, nil)
					mockOrders.EXPECT().GetOrdersByUserID(gomock.Any(), expectedUser.UserID).
						Return(nil, anyErr)
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        anyErr,
		},
		{
			name: "get orders no conntent",
			args: args{
				url:     "/api/user/orders",
				method:  http.MethodGet,
				headers: map[string]string{"Authorization": "Bearer " + expectedToken},
				mockFunc: func() {
					mockAuthToken.EXPECT().ValidateAndParseToken(expectedToken).
						Return(expectedUser.UserID, expectedUser.Username, nil)
					mockOrders.EXPECT().GetOrdersByUserID(gomock.Any(), expectedUser.UserID).
						Return(nil, nil)
				},
			},
			expectedStatusCode: http.StatusNoContent,
			expectedErr:        fmt.Errorf(""), // При статусе 204 тело не возвращается
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.mockFunc != nil {
				tt.args.mockFunc()
			}

			recorder := httptest.NewRecorder()
			// Переводим в байты тело запроса
			var body io.Reader

			if tt.args.requestBody != nil {
				switch tb := tt.args.requestBody.(type) {
				case string: // Передаем тело запроса как есть
					body = bytes.NewReader([]byte(tb))

				default: // Маршалим в структуру
					b, err := json.Marshal(tb)
					assert.NoError(t, err)

					body = bytes.NewReader(b)
				}
			} else {
				body = http.NoBody
			}

			req, err := http.NewRequestWithContext(context.Background(), tt.args.method, tt.args.url, body)
			require.NoError(t, err)

			// добавляем заголовки в запрос
			for k, v := range tt.args.headers {
				req.Header.Add(k, v)
			}

			g.ServeHTTP(recorder, req)

			// код ответа
			assert.Equalf(t, tt.expectedStatusCode, recorder.Code, "body: %s", recorder.Body.String())

			// тело ответа, нет ошибки
			if tt.expectedErr == nil {
				respBody, err := tt.parseResponseBody(recorder.Body.Bytes())
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, respBody)

				// Проверка заголовков
				for k, v := range tt.expectedHeaders {
					assert.Equal(t, v, recorder.Header().Get(k))
				}
			}

			// тело ответа, есть ошибка
			if tt.expectedErr != nil {
				assert.Equal(t, tt.expectedErr.Error(), recorder.Body.String())
			}

		})

	}
}
