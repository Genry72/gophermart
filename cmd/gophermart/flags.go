package main

import (
	"flag"
	"os"
	"strconv"
)

func parseFlags() {
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&flagPgDsn, "d", "postgres://postgres:pass@localhost:5432/gophermart?sslmode=disable", "строка подключения к базе данных")
	flag.StringVar(&flagAuthKey, "authkey", "1111", "token auth key")
	flag.Int64Var(&flagTokenLifeTime, "lt", 3, "token life time token in hour")
	flag.StringVar(&flagAccural, "r", "http://localhost:8080", "адрес системы расчёта начислений: переменная окружения ОС")
	flag.Parse()

	if runAddr := os.Getenv(envRunAddr); runAddr != "" {
		flagRunAddr = runAddr
	}

	if value := os.Getenv(envPgDSN); value != "" {
		flagPgDsn = value
	}

	if value := os.Getenv(envAuthKey); value != "" {
		flagAuthKey = value
	}

	if value := os.Getenv(envTokenLifeTime); value != "" {
		if lt, err := strconv.ParseInt(value, 10, 64); err == nil {
			flagTokenLifeTime = lt
		}
	}

	if value := os.Getenv(envAccuralSystemAddress); value != "" {
		flagAccural = value
	}

}
