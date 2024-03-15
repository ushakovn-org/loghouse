package models

import "time"

// CreateLogParams параметры создания лога
type CreateLogParams struct {
  // Транспорт сервиса
  Transport Transport
  // Эндпоинт запроса к сервису
  Endpoint string
  // Идентификатор трассировки
  TraceID *string
  // Заголовки запроса
  RequestHeaders map[string]string
  // Заголовки ответа
  ResponseHeaders map[string]string
  // Тело запроса
  RequestBody string
  // Тело ответа
  ResponseBody *string
  // Текст ошибки
  ErrorMessage *string
  // Код ответа
  ResponseStatusCode *int32
  // Текстовый код ответа
  ResponseStatusCodeString *string
  // Время отправки запроса
  RequestSendTime *time.Time
  // Время отправки ответа
  ResponseSendTime *time.Time
  // Время обработки запроса
  Latency *string
  // Запрос к сервису завершился успешно
  IsSuccess bool
}
