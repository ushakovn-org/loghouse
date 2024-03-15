package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

// Config конфиг для репозитория логов Storage
type Config struct {
	// Адреса кластера ClickHouse
	Addr []string
}

// Storage репозиторий логов
type Storage struct {
	// Подключение к ClickHouse
	conn clickhouse.Conn
}

// NewStorage создает новый репозиторий логов Storage
func NewStorage(ctx context.Context, config Config) (*Storage, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: config.Addr,

		MaxIdleConns:    5,
		MaxOpenConns:    10,
		ConnMaxLifetime: time.Hour,

		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: time.Second * 30,
	})
	if err = conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("clickhouse conn ping failed: %w", err)
	}
	return &Storage{
		conn: conn,
	}, nil
}
