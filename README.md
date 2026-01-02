# gin-sqlc-demo

A minimal Gin REST service with:
- Postgres via pgxpool (pgx v5)
- sqlc type-safe queries
- OpenTelemetry tracing (OTLP -> OTel Collector -> Jaeger)
- Structured logging (slog) + request_id + trace_id/span_id
- Graceful shutdown

## Endpoints
- `GET /healthz` -> `{"status":"ok"}`
- `GET /v1/hello` -> `{"message":"hello"}`
- `GET /v1/todos` -> list todos from Postgres
- `POST /v1/todos` -> create todo with body `{"title":"..."}`

## Local run (Docker)
1. Start everything:
   ```bash
   make up
# go-service-template-postgres
