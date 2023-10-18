package handlers

import (
	"github.com/Genry72/gophermart/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const successAuthUser = "пользователь успешно аутентифицирован"

func (h *Handler) addUser(c *gin.Context) {
	ctx := c.Request.Context()

	user := &models.UserRegister{}

	if err := c.ShouldBindJSON(user); err != nil {
		h.log.Error("addUser.ShouldBindJSON", zap.Error(err))
		c.String(http.StatusBadRequest, ErrBadBody.Error())

		return
	}

	newUser, err := h.useCases.Users.CreateUser(ctx, user)
	if err != nil {
		h.log.Error("h.useCases.Users.AddUser", zap.Error(err))
		status := checkError(err)
		c.String(status, err.Error())

		return
	}

	token, err := h.authToken.GetToken(newUser)
	if err != nil {
		h.log.Error("h.authToken.GetToken", zap.Error(err))
		status := checkError(err)
		c.String(status, err.Error())
		return
	}

	c.Header("Authorization", "Bearer "+token)

	c.JSON(http.StatusOK, newUser)

}

func (h *Handler) authUser(c *gin.Context) {
	ctx := c.Request.Context()

	user := &models.UserRegister{}

	if err := c.ShouldBindJSON(user); err != nil {
		h.log.Error("authUser.ShouldBindJSON", zap.Error(err))
		c.String(http.StatusBadRequest, ErrBadBody.Error())

		return
	}

	userInfo, err := h.useCases.Users.AuthUser(ctx, user.Username, user.Password)
	if err != nil {
		h.log.Error("h.useCases.Users.AuthUser", zap.Error(err))
		status := checkError(err)
		c.String(status, err.Error())

		return
	}

	token, err := h.authToken.GetToken(userInfo)
	if err != nil {
		h.log.Error("h.authToken.GetToken", zap.Error(err))
		status := checkError(err)
		c.String(status, err.Error())
		return
	}

	c.Header("Authorization", "Bearer "+token)

	c.String(http.StatusOK, successAuthUser)

}
