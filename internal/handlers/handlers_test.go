package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	mockAuth "github.com/Genry72/gophermart/internal/handlers/jwtauth/mocks"
	"github.com/Genry72/gophermart/internal/logger"
	"github.com/Genry72/gophermart/internal/models"
	"github.com/Genry72/gophermart/internal/models/myerrors"
	"github.com/Genry72/gophermart/internal/usecases"
	mockUc "github.com/Genry72/gophermart/internal/usecases/mocks"
	"github.com/Genry72/gophermart/pkg/slices"
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

	expectedDraw := &models.Withdraw{
		//UserID: 1, // userID не возвращаем в ответе
		Order:  luhnorderID,
		Points: 50,
	}

	expectedBalance := &models.Balance{
		Current:   100,
		Withdrawn: 50,
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
		// balance
		{
			name: "getUserBalance positive",
			args: args{
				url:     "/api/user/balance",
				method:  http.MethodGet,
				headers: map[string]string{"Authorization": "Bearer " + expectedToken},
				mockFunc: func() {
					mockAuthToken.EXPECT().ValidateAndParseToken(expectedToken).
						Return(expectedUser.UserID, expectedUser.Username, nil)
					mockBalance.EXPECT().GetUserBalance(gomock.Any(), expectedUser.UserID).Return(expectedBalance, nil)
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       expectedBalance,
			parseResponseBody: func(b []byte) (interface{}, error) {
				respBody := &models.Balance{}
				err := json.Unmarshal(b, respBody)
				if err != nil {
					return nil, fmt.Errorf("json.Unmarshal: %w %s", err, string(b))
				}
				return respBody, nil
			},
		},
		{
			name: "getUserBalance err bd",
			args: args{
				url:     "/api/user/balance",
				method:  http.MethodGet,
				headers: map[string]string{"Authorization": "Bearer " + expectedToken},
				mockFunc: func() {
					mockAuthToken.EXPECT().ValidateAndParseToken(expectedToken).
						Return(expectedUser.UserID, expectedUser.Username, nil)
					mockBalance.EXPECT().GetUserBalance(gomock.Any(), expectedUser.UserID).Return(nil, anyErr)
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        anyErr,
		},
		{
			name: "withdraw positive",
			args: args{
				url:         "/api/user/balance/withdraw",
				method:      http.MethodPost,
				requestBody: expectedDraw,
				headers:     map[string]string{"Authorization": "Bearer " + expectedToken},
				mockFunc: func() {
					mockAuthToken.EXPECT().ValidateAndParseToken(expectedToken).
						Return(expectedUser.UserID, expectedUser.Username, nil)
					exDarw := *expectedDraw
					exDarw.UserID = 1 // UserID подменится из контекта
					mockBalance.EXPECT().Withdraw(gomock.Any(), &exDarw).Return(nil)
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       expectedDraw,
			parseResponseBody: func(b []byte) (interface{}, error) {
				respBody := &models.Withdraw{}
				err := json.Unmarshal(b, respBody)
				if err != nil {
					return nil, fmt.Errorf("json.Unmarshal: %w %s", err, string(b))
				}
				return respBody, nil
			},
		},
		{
			name: "withdraw empry body",
			args: args{
				url:         "/api/user/balance/withdraw",
				method:      http.MethodPost,
				requestBody: "",
				headers:     map[string]string{"Authorization": "Bearer " + expectedToken},
				mockFunc: func() {
					mockAuthToken.EXPECT().ValidateAndParseToken(expectedToken).
						Return(expectedUser.UserID, expectedUser.Username, nil)
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErr:        ErrBadBody,
		},
		{
			name: "withdraw err format order",
			args: args{
				url:         "/api/user/balance/withdraw",
				method:      http.MethodPost,
				requestBody: "{}",
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
			name: "withdraw not valid orderID",
			args: args{
				url:    "/api/user/balance/withdraw",
				method: http.MethodPost,
				requestBody: &models.Withdraw{
					UserID: 1,
					Order:  "123",
				},
				headers: map[string]string{"Authorization": "Bearer " + expectedToken},
				mockFunc: func() {
					mockAuthToken.EXPECT().ValidateAndParseToken(expectedToken).
						Return(expectedUser.UserID, expectedUser.Username, nil)
				},
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedErr:        myerrors.ErrBadFormatOrder,
		},
		{
			name: "withdraw no money",
			args: args{
				url:         "/api/user/balance/withdraw",
				method:      http.MethodPost,
				requestBody: expectedDraw,
				headers:     map[string]string{"Authorization": "Bearer " + expectedToken},
				mockFunc: func() {
					mockAuthToken.EXPECT().ValidateAndParseToken(expectedToken).
						Return(expectedUser.UserID, expectedUser.Username, nil)
					exDarw := *expectedDraw
					exDarw.UserID = 1 // UserID подменится из контекта
					mockBalance.EXPECT().Withdraw(gomock.Any(), &exDarw).Return(myerrors.ErrNoMoney)
				},
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedErr:        myerrors.ErrNoMoney,
		},
		{
			name: "withdraw err bd",
			args: args{
				url:         "/api/user/balance/withdraw",
				method:      http.MethodPost,
				requestBody: expectedDraw,
				headers:     map[string]string{"Authorization": "Bearer " + expectedToken},
				mockFunc: func() {
					mockAuthToken.EXPECT().ValidateAndParseToken(expectedToken).
						Return(expectedUser.UserID, expectedUser.Username, nil)
					exDarw := *expectedDraw
					exDarw.UserID = 1 // UserID подменится из контекта
					mockBalance.EXPECT().Withdraw(gomock.Any(), &exDarw).Return(anyErr)
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        anyErr,
		},
		{
			name: "withdrawals positive",
			args: args{
				url:     "/api/user/withdrawals",
				method:  http.MethodGet,
				headers: map[string]string{"Authorization": "Bearer " + expectedToken},

				mockFunc: func() {
					mockAuthToken.EXPECT().ValidateAndParseToken(expectedToken).
						Return(expectedUser.UserID, expectedUser.Username, nil)
					mockBalance.EXPECT().Withdrawals(gomock.Any(), expectedUser.UserID).Return([]*models.Withdraw{expectedDraw}, nil)
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       []*models.Withdraw{expectedDraw},
			parseResponseBody: func(b []byte) (interface{}, error) {
				respBody := []*models.Withdraw{expectedDraw}
				err := json.Unmarshal(b, &respBody)
				if err != nil {
					return nil, fmt.Errorf("json.Unmarshal: %w %s", err, string(b))
				}
				return respBody, nil
			},
		},
		{
			name: "withdrawals err bd",
			args: args{
				url:     "/api/user/withdrawals",
				method:  http.MethodGet,
				headers: map[string]string{"Authorization": "Bearer " + expectedToken},

				mockFunc: func() {
					mockAuthToken.EXPECT().ValidateAndParseToken(expectedToken).
						Return(expectedUser.UserID, expectedUser.Username, nil)
					mockBalance.EXPECT().Withdrawals(gomock.Any(), expectedUser.UserID).Return(nil, anyErr)
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        anyErr,
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

func TestNewHandlerAndAuth(t *testing.T) {
	useCases := new(usecases.Usecase)
	hostPort := "localhost:8080"
	tokenKey := "token"
	jwtLifeTime := time.Hour
	zapLogger := logger.NewZapLogger("info")

	handler := NewHandler(useCases, hostPort, tokenKey, jwtLifeTime, zapLogger)

	// роуты по которым нет роверки токена
	noAuthRoute := []string{"/api/user/register", "/api/user/login"}

	for _, v := range handler.ginEngine.Routes() {
		req, err := http.NewRequestWithContext(context.Background(), v.Method, v.Path, nil)
		assert.NoError(t, err)

		resp := httptest.NewRecorder()

		handler.ginEngine.ServeHTTP(resp, req)

		if !slices.Inslice(noAuthRoute, v.Path) {
			assert.Equal(t, http.StatusUnauthorized, resp.Code)
		} else {
			assert.Equal(t, http.StatusBadRequest, resp.Code)
		}
	}

}

func TestStartHandler(t *testing.T) {
	useCases := new(usecases.Usecase)
	hostPort := ":0"
	tokenKey := ""
	jwtLifeTime := time.Hour
	zapLogger := logger.NewZapLogger("info")

	handler := NewHandler(useCases, hostPort, tokenKey, jwtLifeTime, zapLogger)

	server := httptest.NewServer(handler.ginEngine)
	defer server.Close()

	go func() {
		if err := handler.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			assert.NoError(t, err)
		}
	}()

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := handler.Stop(ctx); err != nil {
			assert.NoError(t, err)
		}
	}()

	time.Sleep(time.Second)

	req, err := http.NewRequestWithContext(context.Background(),
		"POST", "http://"+server.Listener.Addr().String()+"/api/user/login", nil)
	assert.NoError(t, err)

	client := &http.Client{}

	res, err := client.Do(req)
	assert.NoError(t, err)

	defer func() {
		err := res.Body.Close()
		assert.NoError(t, err)
	}()

}
