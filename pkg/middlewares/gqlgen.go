package middlewares

import (
  "context"
  "fmt"
  "strings"
  "time"

  "github.com/99designs/gqlgen/graphql"
  "github.com/samber/lo"
  log "github.com/sirupsen/logrus"
  "github.com/ushakovn-org/loghouse/internal/pb/loghouse"
  "github.com/ushakovn-org/loghouse/pkg/ctxdetach"
  "github.com/ushakovn/boiler/pkg/tracing/tracer"
  "google.golang.org/protobuf/types/known/timestamppb"
)

// GqlgenResponseMiddleware middleware для GraphQL операций,
// отправляет запрос на сохранение логов в loghouse
func GqlgenResponseMiddleware(config Config) graphql.ResponseMiddleware {
  if err := config.Validate(); err != nil {
    log.Fatalf("loghouse: config invalid: %v", err)
  }
  client := loghouse.NewLogHouseClient(config.Conn)

  return func(ctx context.Context, handler graphql.ResponseHandler) *graphql.Response {
    opResp := handler(ctx)

    respSendTime := time.Now().UTC()

    timeoutCtx, cancel := context.WithTimeout(ctxdetach.Do(ctx), config.Timeout)

    go func(ctx context.Context) {
      defer cancel()

      if !graphql.HasOperationContext(ctx) {
        return
      }
      opCtx := graphql.GetOperationContext(ctx)

      if opCtx == nil || opCtx.Operation == nil {
        return
      }
      endpoint := fmt.Sprintf("%s %s", opCtx.Operation.Operation, opCtx.Operation.Name)

      reqSendTime := opCtx.Stats.OperationStart.UTC()

      latency := lo.ToPtr(fmt.Sprintf("%d ms", respSendTime.Sub(reqSendTime).Milliseconds()))

      var traceID *string

      if trace := tracer.SpanFromContext(ctx).SpanContext().TraceID(); trace.IsValid() {
        traceID = lo.ToPtr(trace.String())
      }

      reqBody := opCtx.RawQuery
      respBody := lo.ToPtr(string(opResp.Data))

      isSuccess := true
      var errMsg *string

      if opErrs := graphql.GetErrors(ctx); len(opErrs) > 0 {
        isSuccess = false
        errMsg = lo.ToPtr(opErrs.Error())
      }

      reqHeaders := make(map[string]string)

      for key, values := range opCtx.Headers {
        reqHeaders[key] = strings.Join(values, ",")
      }

      respHeaders := make(map[string]string)

      for key, value := range opResp.Extensions {
        reqHeaders[key] = fmt.Sprintf("%v", value)
      }

      createLogReq := &loghouse.CreateLogAsyncRequest{
        Transport:        loghouse.Transport_TRANSPORT_GRAPHQL,
        Endpoint:         endpoint,
        TraceId:          traceID,
        RequestHeaders:   reqHeaders,
        ResponseHeaders:  respHeaders,
        RequestBody:      reqBody,
        ResponseBody:     respBody,
        ErrorMessage:     errMsg,
        IsSuccess:        isSuccess,
        Latency:          latency,
        RequestSendTime:  timestamppb.New(reqSendTime),
        ResponseSendTime: timestamppb.New(respSendTime),
      }
      if _, err := client.CreateLogAsync(ctx, createLogReq); err != nil {
        log.Warnf("loghouse: create log async call failed: %v", err)
        return
      }
    }(timeoutCtx)

    return opResp
  }
}
