package storage

import (
  "context"
  "fmt"
  "reflect"
  "strings"

  "github.com/ushakovn-org/loghouse/internal/pkg/models"

  "github.com/ClickHouse/clickhouse-go/v2"
  sq "github.com/Masterminds/squirrel"
)

const (
  // logsTableName таблица с логами
  logsTableName = "logs"
)

var (
  // logFields все поля модели лога для таблицы logsTableName
  logFields = buildFields(models.Log{})
)

// Подавляет ошибки о неиспользуемой переменной
var (
  _ = logFields
)

// builder возвращает builder для sql запросов
func (s *Storage) builder() sq.StatementBuilderType {
  return sq.StatementBuilderType{}.PlaceholderFormat(sq.Question)
}

// buildFields формирует строку, содержащую поля модели с тегами "db:"
func buildFields(model any) string {
  typ := reflect.TypeOf(model)
  count := typ.NumField()
  fields := make([]string, 0, count)

  for index := 0; index < typ.NumField(); index++ {
    tag := typ.Field(index).Tag.Get("db")
    if tag == "" {
      continue
    }
    fields = append(fields, tag)
  }
  return strings.Join(fields, ",")
}

// buildSql формирует строку sql выражения и набор его аргументов
func buildSql(query sq.Sqlizer) (sql string, args []any) {
  sql, args, _ = query.ToSql()
  return sql, args
}

// execx выполняет sql запрос
func execx(ctx context.Context, conn clickhouse.Conn, query sq.Sqlizer) error {
  const wait = false
  sql, args := buildSql(query)

  if err := conn.AsyncInsert(ctx, sql, wait, args...); err != nil {
    return fmt.Errorf("conn.AsyncInsert: %w", err)
  }
  return nil
}
