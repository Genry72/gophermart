package myerrors

import "errors"

var ErrUserAlreadyExist = errors.New("логин уже занят")
var ErrUnauthorized = errors.New("неверная пара логин/пароль")
