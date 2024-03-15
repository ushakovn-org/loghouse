// Code generated by Boiler; DO NOT EDIT.

package config

import (
	"context"
	"github.com/ushakovn/boiler/pkg/config"
	"github.com/ushakovn/boiler/pkg/config/types"
)

const (
	// Размер пула воркеров
	LoghouseWorkersCount configKey = "loghouse_workers_count"
	// Таймаут для воркеров
	LoghouseWorkersTimeout configKey = "loghouse_workers_timeout"
)

const (
	// Адрес кластера ClickHouse
	ClickhouseAddress configKey = "clickhouse_address"
)

// configKey strict type for config key
type configKey string

// Get value of the specified key
func Get(ctx context.Context, key configKey) types.Value {
	return config.ContextClient(ctx).GetValue(ctx, string(key))
}

// Watch watches changes to the value of the specified key
func Watch(ctx context.Context, key configKey, action func(value types.Value)) {
	config.ContextClient(ctx).WatchValue(ctx, string(key), action)
}
