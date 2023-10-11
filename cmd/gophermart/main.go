package main

import (
	"context"
	"errors"
	"github.com/Genry72/gophermart/internal/handlers"
	"github.com/Genry72/gophermart/internal/logger"
	"github.com/Genry72/gophermart/internal/repositories/postgre"
	"github.com/Genry72/gophermart/internal/usecases"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	flagRunAddr       string
	flagPgDsn         string
	flagAuthKey       string
	flagTokenLifeTime int64
	flagAccural       string
)

const (
	envRunAddr = "RUN_ADDRESS"
	// Строка с адресом подключения к БД
	envPgDSN = "DATABASE_URI"
	// Ключ для генерации токена
	envAuthKey = "AUTH_KEY"
	// Время жизни токена в часах
	envTokenLifeTime = "TOKEN_LIFE_TIME"
	// Адрес системы расчёта начислений: переменная окружения ОС
	envAccuralSystemAddress = "ACCRUAL_SYSTEM_ADDRESS"
)

func main() {
	zapLogger := logger.NewZapLogger("info")

	defer func() {
		_ = zapLogger.Sync()
	}()

	// обрабатываем аргументы командной строки
	parseFlags()

	repo, err := postgre.NewPGStorage(flagPgDsn, zapLogger)
	if err != nil {
		zapLogger.Fatal("postgre.NewPGStorage", zap.Error(err))
	}

	zapLogger.Info("Connect to db success")

	defer repo.Stop()

	usecase := usecases.NewUsecase(repo, zapLogger)

	zapLogger.Info("Starting server", zap.String("port", flagRunAddr))

	server := handlers.NewHandler(usecase, flagRunAddr, flagAuthKey, time.Duration(flagTokenLifeTime)*time.Hour, zapLogger)

	go func() {
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zapLogger.Fatal("run server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		zapLogger.Error("Server Shutdown:", zap.Error(err))
	}

	zapLogger.Info("Server exiting")
}
