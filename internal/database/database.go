package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ConfigDB struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type DB struct {
	Pool *pgxpool.Pool
	Ctx  context.Context
}

func NewDB(ctx context.Context, cfg ConfigDB) (*DB, error) {
	conn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)

	config, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфига: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("База данных не отвечает: %w", err)
	}

	log.Println("База данных успешно подключена!")
	return &DB{
		Pool: pool,
		Ctx:  ctx,
	}, nil
}
