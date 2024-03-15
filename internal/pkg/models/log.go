package models

import (
	"time"
)

type Transport int

const (
	TransportUnknown Transport = 0
	TransportGRPC    Transport = 1
	TransportGraphQL Transport = 2
	TransportHTTP    Transport = 3
)

type RequestHeaders string

type Log struct {
	Transport        Transport         `db:"transport"`
	TraceID          *string           `db:"trace_id"`
	RequestHeaders   map[string]string `db:"request_headers"`
	RequestBody      string            `db:"request_body"`
	ResponseBody     string            `db:"response_body"`
	RequestSendTime  *time.Time        `db:"request_send_time"`
	ResponseSendTime *time.Time        `db:"response_send_time"`
	Latency          *string           `db:"latency"`
}
