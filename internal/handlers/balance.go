package handlers

import (
	"github.com/Genry72/gophermart/internal/models"
	"github.com/Genry72/gophermart/internal/models/myerrors"
	"github.com/gin-gonic/gin"
	"github.com/theplant/luhn"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (h *Handler) getUserBalance(c *gin.Context) {
	userID, ok := c.Request.Context().Value(models.CtxKeyUserID{}).(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, myerrors.ErrBadCtxUserID.Error())
		return
	}

	ctx := c.Request.Context()

	balance, err := h.useCases.Balances.GetUserBalance(ctx, userID)
	if err != nil {
		h.log.Error(" h.useCases.Users.GetUserBalanc", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, balance)
}

func (h *Handler) withdraw(c *gin.Context) {
	userID, ok := c.Request.Context().Value(models.CtxKeyUserID{}).(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, myerrors.ErrBadCtxUserID.Error())
		return
	}

	draw := models.Withdraw{}
	if err := c.ShouldBindJSON(&draw); err != nil {
		h.log.Error("withdraw.ShouldBindJSON", zap.Error(err))
		c.String(http.StatusBadRequest, ErrBadBody.Error())

		return
	}

	orderID, err := strconv.ParseInt(draw.Order, 10, 64)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, myerrors.ErrBadFormatOrder.Error())
		return
	}

	if !luhn.Valid(int(orderID)) {
		c.JSON(http.StatusUnprocessableEntity, myerrors.ErrBadFormatOrder.Error())
		return
	}

	draw.UserID = userID

	ctx := c.Request.Context()

	if err := h.useCases.Balances.Withdraw(ctx, &draw); err != nil {
		h.log.Error(" h.useCases.Users.Withdraw", zap.Error(err))
		status := checkError(err)
		c.JSON(status, err.Error())

		return
	}

	c.JSON(http.StatusOK, "ok")
}

func (h *Handler) withdrawals(c *gin.Context) {
	userID, ok := c.Request.Context().Value(models.CtxKeyUserID{}).(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, myerrors.ErrBadCtxUserID.Error())
		return
	}

	ctx := c.Request.Context()

	drawals, err := h.useCases.Balances.Withdrawals(ctx, userID)
	if err != nil {
		h.log.Error("h.useCases.Users.Withdrawals", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, drawals)
}
