package handlers

import (
	"github.com/Genry72/gophermart/internal/models"
	"github.com/Genry72/gophermart/internal/models/myerrors"
	"github.com/gin-gonic/gin"
	"github.com/theplant/luhn"
	"io"
	"net/http"
	"strconv"
)

func (h *Handler) uploadOrder(c *gin.Context) {
	userID, ok := c.Request.Context().Value(models.CtxKeyUserID{}).(int64)
	if !ok {
		c.String(http.StatusInternalServerError, myerrors.ErrBadCtxUserID.Error())
		return
	}

	// Проверка на валидность номера заказа из тела запроса
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	orderID, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		c.String(http.StatusUnprocessableEntity, myerrors.ErrBadFormatOrder.Error())
		return
	}

	if !luhn.Valid(int(orderID)) {
		c.String(http.StatusUnprocessableEntity, myerrors.ErrBadFormatOrder.Error())
		return
	}

	ctx := c.Request.Context()

	order, err := h.useCases.Orders.AddOrder(ctx, orderID, userID)
	if err != nil {
		status := checkError(err)
		c.String(status, err.Error())

		return
	}

	c.JSON(http.StatusAccepted, order)
}

func (h *Handler) getOrders(c *gin.Context) {
	userID, ok := c.Request.Context().Value(models.CtxKeyUserID{}).(int64)
	if !ok {
		c.String(http.StatusInternalServerError, myerrors.ErrBadCtxUserID.Error())
		return
	}

	ctx := c.Request.Context()

	orders, err := h.useCases.Orders.GetOrdersByUserID(ctx, userID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if len(orders) == 0 {
		c.String(http.StatusNoContent, "") // При статусе 204 тело не возвращается
		return
	}

	c.JSON(http.StatusOK, orders)
}
