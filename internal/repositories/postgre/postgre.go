package postgre

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"time"
)

type PGStorage struct {
	Conn *sqlx.DB
	log  *zap.Logger
}

func NewPGStorage(dsn string, log *zap.Logger) (*PGStorage, error) {

	if err := migration(dsn); err != nil {
		return nil, fmt.Errorf("migration: %w", err)
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Connect: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(10 * time.Second)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(1 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	pg := &PGStorage{
		Conn: db,
		log:  log,
	}

	return pg, nil
}

func migration(dsn string) error {
	m, err := migrate.New(
		"file://internal/repositories/postgre/migrations",
		dsn)
	if err != nil {
		return fmt.Errorf("migrate.New: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("m.Up: %w", err)
	}

	return nil
}

func (pg *PGStorage) Stop() {
	if err := pg.Conn.Close(); err != nil {
		pg.log.Error("pg.Conn.Close", zap.Error(err))
		return
	}

	pg.log.Info("Database success closed")
}
