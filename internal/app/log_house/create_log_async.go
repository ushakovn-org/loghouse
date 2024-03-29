// Code generated by Boiler; YOU MUST CHANGE THIS.

package log_house

import (
  "context"
  "time"

  "github.com/samber/lo"
  log "github.com/sirupsen/logrus"
  "github.com/ushakovn-org/loghouse/internal/config"
  desc "github.com/ushakovn-org/loghouse/internal/pb/loghouse"
  "github.com/ushakovn-org/loghouse/internal/pkg/models"
  "github.com/ushakovn-org/loghouse/pkg/ctxdetach"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
)

// CreateLogAsync асинхронно создает лог обработки запроса
func (s *LogHouse) CreateLogAsync(ctx context.Context, req *desc.CreateLogAsyncRequest) (*desc.CreateLogAsyncResponse, error) {
  if err := req.ValidateAll(); err != nil {
  	return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
  }

  ctx, cancel := context.WithTimeout(
  	ctxdetach.Do(ctx),
  	config.Get(ctx, config.LoghouseWorkersTimeout).Duration(),
  )

  s.pool.Submit(func() {
  	defer cancel()

  	var (
  		reqSendTime  *time.Time
  		respSendTime *time.Time
  	)
  	if req.RequestSendTime != nil {
  		reqSendTime = lo.ToPtr(req.RequestSendTime.AsTime())
  	}
  	if req.ResponseSendTime != nil {
  		respSendTime = lo.ToPtr(req.ResponseSendTime.AsTime())
  	}
  	err := s.storage.CreateLogAsync(ctx, models.CreateLogParams{
  		Transport:                models.Transport(req.Transport),
  		Endpoint:                 req.Endpoint,
  		TraceID:                  req.TraceId,
  		RequestHeaders:           req.RequestHeaders,
  		ResponseHeaders:          req.ResponseHeaders,
  		RequestBody:              req.RequestBody,
  		ResponseBody:             req.ResponseBody,
  		ErrorMessage:             req.ErrorMessage,
  		ResponseStatusCode:       req.ResponseStatusCode,
  		ResponseStatusCodeString: req.ResponseStatusCodeString,
  		RequestSendTime:          reqSendTime,
  		ResponseSendTime:         respSendTime,
  		Latency:                  req.Latency,
  		IsSuccess:                req.IsSuccess,
  	})
  	if err != nil {
  		log.Errorf("LogHouse.CreateLogAsync: storage.CreateLogAsync: %v", err)
  	}
  })

  return &desc.CreateLogAsyncResponse{}, nil
}
