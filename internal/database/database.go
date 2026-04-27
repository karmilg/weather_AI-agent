package database

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/karmilg/weather_AI-agent/config"
)

type DB struct {
	Pool *pgxpool.Pool
	Ctx  context.Context
}

func NewDB(ctx context.Context, cfg *config.Config) (*DB, error) {
	log.Printf("DEBUG: DatabaseURL = '%s'", cfg.DatabaseURL)
	log.Printf("DEBUG: DBHost = '%s'", cfg.DBHost)
	var connString string

	if cfg.DatabaseURL != "" {
		connString = cfg.DatabaseURL

		if !strings.Contains(connString, "sslmode=") {
			if strings.Contains(connString, "?") {
				connString += "&sslmode=require"
			} else {
				connString += "?sslmode=require"
			}
		}
	} else {
		connString = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
		)
	}

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга DATABASE_URL: %w", err)
	}

	poolConfig.MaxConns = 25
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = 30 * 60
	poolConfig.MaxConnIdleTime = 5 * 60

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания пула: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("База данных не отвечает: %w", err)
	}

	log.Println("✅ База данных успешно подключена! (pgxpool)")
	return &DB{Pool: pool, Ctx: ctx}, nil
}
