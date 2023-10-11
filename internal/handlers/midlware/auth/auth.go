package auth

import (
	"context"
	"github.com/Genry72/gophermart/internal/handlers/jwtauth"
	"github.com/Genry72/gophermart/internal/models"
	"github.com/Genry72/gophermart/internal/models/myerrors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Auth Аутентификация пользователя
func Auth(a jwtauth.Auther) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		tokenSplited := strings.Split(authHeader, " ")

		if len(tokenSplited) != 2 || strings.ToLower(tokenSplited[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, myerrors.ErrBadAuthHeader.Error())
			return
		}

		userID, _, err := a.ValidateAndParseToken(tokenSplited[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, myerrors.ErrBadAuthHeader.Error())
			return
		}

		ctx := context.WithValue(c.Request.Context(), models.CtxKeyUserID{}, userID)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}

}
