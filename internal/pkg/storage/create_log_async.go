package storage

import (
  "context"
  "fmt"

  "github.com/ushakovn-org/loghouse/internal/pkg/models"

  "github.com/ushakovn/boiler/pkg/tracing/tracer"
  "go.opentelemetry.io/otel/attribute"
)

// CreateLogAsync асинхронно создает лог обработки запроса
func (s *Storage) CreateLogAsync(ctx context.Context, params models.CreateLogParams) error {
  ctx, span := tracer.StartContextWithSpan(ctx, "storage.CreateLog")
  defer span.End()

  fields := map[string]any{
    "transport":    params.Transport,
    "endpoint":     params.Endpoint,
    "request_body": params.RequestBody,
    "is_success":   params.IsSuccess,
  }
  if params.TraceID != nil {
    fields["trace_id"] = *params.TraceID
  }
  if len(params.RequestHeaders) > 0 {
    fields["request_headers"] = params.RequestHeaders
  }
  if len(params.ResponseHeaders) > 0 {
    fields["response_headers"] = params.RequestHeaders
  }
  if params.ResponseBody != nil {
    fields["response_body"] = *params.ResponseBody
  }
  if params.ErrorMessage != nil {
    fields["error_message"] = *params.ErrorMessage
  }
  if params.ResponseStatusCode != nil {
    fields["response_status_code"] = *params.ResponseStatusCode
  }
  if params.ResponseStatusCodeString != nil {
    fields["response_status_code_string"] = *params.ResponseStatusCodeString
  }
  if params.RequestSendTime != nil {
    fields["request_send_time"] = *params.ResponseSendTime
  }
  if params.RequestSendTime != nil {
    fields["response_send_time"] = *params.ResponseSendTime
  }
  if params.Latency != nil {
    fields["latency"] = *params.Latency
  }

  for key, value := range fields {
    strValue := fmt.Sprintf("%v", value)

    span.SetAttributes(attribute.KeyValue{
      Key:   attribute.Key(key),
      Value: attribute.StringValue(strValue),
    })
  }

  query := s.builder().
    Insert(logsTableName).
    SetMap(fields)

  return execx(ctx, s.conn, query)
}
