create table logs
(
    transport                   UInt32   not null,
    endpoint                    String   not null,
    trace_id                    String,
    request_headers             Map(String, String),
    response_headers            Map(String, String),
    request_body                String   not null,
    response_body               String,
    error_message               String,
    response_status_code        Int32,
    response_status_code_string String,
    request_send_time           DateTime,
    response_send_time          DateTime,
    latency                     String,
    is_success                  Bool     not null,
    created_at                  DateTime default now()
) ENGINE = Log();