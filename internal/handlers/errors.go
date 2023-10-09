package handlers

import (
	"errors"
	"github.com/Genry72/gophermart/internal/models/myerrors"
	"net/http"
)

var ErrBadBody = errors.New("неверный формат запроса")

func checkError(err error) int {
	switch {
	case errors.Is(err, myerrors.ErrUserAlreadyExist):
		return http.StatusConflict
	case errors.Is(err, myerrors.ErrUnauthorized):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
