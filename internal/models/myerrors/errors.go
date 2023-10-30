package myerrors

import "errors"

var (
	ErrUserAlreadyExist         = errors.New("логин уже занят")
	ErrUnauthorized             = errors.New("неверная пара логин/пароль")
	ErrBadAuthHeader            = errors.New("пользователь не аутентифицирован")
	ErrBadCtxUserID             = errors.New("userID из контекста имеет некорректный формат")
	ErrBadFormatOrder           = errors.New("неверный формат номера заказа")
	ErrOrderUploadByAnotherUser = errors.New("номер заказа уже был загружен другим пользователем")
	ErrOrderAlreadyUploadByUser = errors.New("номер заказа уже был загружен этим пользователем")
	ErrStatusCodeNotCorrect     = errors.New("некорректный код ответа")
	ErrNoMoney                  = errors.New("на счету недостаточно средств")
)
