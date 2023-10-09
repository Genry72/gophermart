package handlers

import (
	"context"
	"github.com/Genry72/gophermart/internal/handlers/jwtAuth"
	"github.com/Genry72/gophermart/internal/handlers/jwtAuth/jwtToken"
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
	authToken jwtAuth.Auther
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
		Addr:    hostPort,
		Handler: g,
	}

	h := &Handler{
		useCases:  useCases,
		log:       log,
		ginEngine: g,
		server:    srv,
		authToken: jwtToken.NewJwtToken(tokenKey, jwtLifeTime),
	}

	h.initRoutes()

	return h
}

func (s *Handler) Start() error {
	return s.server.ListenAndServe()
}

func (s *Handler) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (h *Handler) initRoutes() {
	api := h.ginEngine.Group("api")
	{
		user := api.Group("user")
		{
			user.POST("/register", h.addUser)
			user.POST("/login", h.authUser)
		}
	}
}
