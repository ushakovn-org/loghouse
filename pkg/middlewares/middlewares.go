package middlewares

import (
  "fmt"
  "time"

  "google.golang.org/grpc"
)

// Config конфигурация для UnaryServerInterceptor
type Config struct {
  // Conn grpc подключение к сервису loghouse
  Conn grpc.ClientConnInterface
  // Timeout таймаут на обращение к сервису loghouse
  Timeout time.Duration
}

// Validate валидация для Config
func (c Config) Validate() error {
  if c.Conn == nil {
    return fmt.Errorf("grpc connection must be specified")
  }
  if c.Timeout <= 0 {
    return fmt.Errorf("non negative timeout must be specified")
  }
  return nil
}

// WithDefault применяет значения по умолчанию для Config
func (c Config) WithDefault() Config {
  c.Timeout = 500 * time.Millisecond
  return c
}

// WithConn проставляет Config.Conn и возвращает Config
func (c Config) WithConn(conn grpc.ClientConnInterface) Config {
  c.Conn = conn
  return c
}

// WithTimeout проставляет Config.Timeout и возвращает Config
func (c Config) WithTimeout(timeout time.Duration) Config {
  c.Timeout = timeout
  return c
}

// NewConfig возвращает новый Config
// со значениями по умолчанию
func NewConfig() Config {
  return new(Config).WithDefault()
}
