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
	case errors.Is(err, myerrors.ErrOrderAlreadyUploadByUser):
		return http.StatusOK
	case errors.Is(err, myerrors.ErrOrderUploadByAnotherUser):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
