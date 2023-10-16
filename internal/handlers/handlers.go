package handlers

import (
	"context"
	"github.com/Genry72/gophermart/internal/handlers/jwtauth"
	"github.com/Genry72/gophermart/internal/handlers/jwtauth/jwttoken"
	"github.com/Genry72/gophermart/internal/handlers/midlware/auth"
	"github.com/Genry72/gophermart/internal/handlers/midlware/gzip"
	midlwareLog "github.com/Genry72/gophermart/internal/handlers/midlware/log"
	"github.com/Genry72/gophermart/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Handler struct {
	useCases  *usecases.Usecase
	log       *zap.Logger
	ginEngine *gin.Engine
	server    *http.Server
	authToken jwtauth.Auther
}

func NewHandler(useCases *usecases.Usecase,
	hostPort string, tokenKey string, jwtLifeTime time.Duration, log *zap.Logger) *Handler {
	gin.SetMode(gin.ReleaseMode)
	g := gin.New()

	// Подключаем логирование и работу со сжатием запросов
	g.Use(midlwareLog.ResponseLogger(log))
	g.Use(midlwareLog.RequestLogger(log))
	g.Use(gzip.Gzip(log))

	srv := &http.Server{
		Addr:              hostPort,
		Handler:           g,
		ReadHeaderTimeout: time.Second,
	}

	h := &Handler{
		useCases:  useCases,
		log:       log,
		ginEngine: g,
		server:    srv,
		authToken: jwttoken.NewJwtToken(tokenKey, jwtLifeTime),
	}

	h.initRoutes()

	return h
}

func (h *Handler) Start() error {
	return h.server.ListenAndServe()
}

func (h *Handler) Stop(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}

func (h *Handler) initRoutes() {
	api := h.ginEngine.Group("api")
	{
		user := api.Group("user")
		{
			user.POST("/register", h.addUser)
			user.POST("/login", h.authUser)
			user.POST("/orders", auth.Auth(h.authToken), h.uploadOrder)
			user.GET("/orders", auth.Auth(h.authToken), h.getOrders)
			user.GET("/balance", auth.Auth(h.authToken), h.getUserBalance)
			user.GET("/withdrawals", auth.Auth(h.authToken), h.withdrawals)
		}

		balance := user.Group("/balance")
		{
			balance.POST("/withdraw", auth.Auth(h.authToken), h.withdraw)
		}
	}
}
