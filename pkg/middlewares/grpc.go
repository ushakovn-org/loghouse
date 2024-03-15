package middlewares

import (
  "context"
  "fmt"
  "strings"
  "time"

  "github.com/ushakovn-org/loghouse/internal/pb/loghouse"
  "github.com/ushakovn-org/loghouse/pkg/ctxdetach"

  "github.com/samber/lo"
  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/pkg/tracing/tracer"
  "google.golang.org/grpc"
  "google.golang.org/grpc/metadata"
  "google.golang.org/grpc/status"
  "google.golang.org/protobuf/encoding/protojson"
  "google.golang.org/protobuf/proto"
  "google.golang.org/protobuf/types/known/timestamppb"
)

// UnaryServerInterceptor interceptor для унарных gRPC запросов,
// отправляет запрос на сохранение логов в loghouse
func UnaryServerInterceptor(config Config) grpc.UnaryServerInterceptor {
  if err := config.Validate(); err != nil {
    log.Fatalf("loghouse: config invalid: %v", err)
  }
  client := loghouse.NewLogHouseClient(config.Conn)

  return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
    reqSendTime := time.Now().UTC()

    resp, respErr := handler(ctx, req)

    respSendTime := time.Now().UTC()

    timeoutCtx, cancel := context.WithTimeout(ctxdetach.Do(ctx), config.Timeout)

    go func(ctx context.Context) {
      defer cancel()

      latency := lo.ToPtr(fmt.Sprintf("%d ms", respSendTime.Sub(reqSendTime).Milliseconds()))

      var traceID *string

      if trace := tracer.SpanFromContext(ctx).SpanContext().TraceID(); trace.IsValid() {
        traceID = lo.ToPtr(trace.String())
      }
      endpoint := info.FullMethod

      if endpoint == "" {
        log.Warnf("loghouse: grpc method name not specified")
        return
      }

      msg, ok := req.(proto.Message)
      if !ok {
        log.Warnf("loghouse: failed request assertion to proto message: %T", req)
        return
      }
      buf, err := protojson.Marshal(msg)
      if err != nil {
        log.Warnf("loghouse: failed request proto marshalling: %v", err)
        return
      }
      reqBody := string(buf)

      respStatus := status.Code(respErr)

      statusCode := lo.ToPtr(int32(respStatus))
      statusString := lo.ToPtr(respStatus.String())

      var respBody *string

      if resp != nil {
        if msg, ok = resp.(proto.Message); !ok {
          log.Warnf("loghouse: failed response assertion to proto message: %T", req)
          return
        }
        if buf, err = protojson.Marshal(msg); err != nil {
          log.Warnf("loghouse: failed response proto marshalling: %v", err)
          return
        }
        respBody = lo.ToPtr(string(buf))
      }

      isSuccess := true
      var errMsg *string

      if respErr != nil {
        isSuccess = false
        errMsg = lo.ToPtr(respErr.Error())
      }

      mdHeaders := func(md metadata.MD, ok bool) map[string]string {
        if md, ok = metadata.FromIncomingContext(ctx); ok {
          headers := make(map[string]string, len(md))

          for key, values := range md {
            headers[key] = strings.Join(values, ",")
          }
          return headers
        }
        return nil
      }

      reqHeaders := mdHeaders(metadata.FromIncomingContext(ctx))
      respHeaders := mdHeaders(metadata.FromOutgoingContext(ctx))

      createLogReq := &loghouse.CreateLogAsyncRequest{
        Transport:                loghouse.Transport_TRANSPORT_GRPC,
        Endpoint:                 endpoint,
        TraceId:                  traceID,
        RequestHeaders:           reqHeaders,
        ResponseHeaders:          respHeaders,
        RequestBody:              reqBody,
        ResponseBody:             respBody,
        ErrorMessage:             errMsg,
        ResponseStatusCode:       statusCode,
        ResponseStatusCodeString: statusString,
        IsSuccess:                isSuccess,
        Latency:                  latency,
        RequestSendTime:          timestamppb.New(reqSendTime),
        ResponseSendTime:         timestamppb.New(respSendTime),
      }
      if _, err = client.CreateLogAsync(ctx, createLogReq); err != nil {
        log.Warnf("loghouse: create log async call failed: %v", err)
        return
      }
    }(timeoutCtx)

    return resp, respErr
  }
}
