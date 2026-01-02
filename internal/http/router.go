package http

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"gin-sqlc-demo/internal/db/sqlc"
	"gin-sqlc-demo/internal/http/handlers"
	"gin-sqlc-demo/internal/http/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Deps struct {
	Logger *slog.Logger
	DB     *pgxpool.Pool
}

func NewRouter(d Deps) *gin.Engine {
	r := gin.New()

	q := sqlc.New(d.DB)

	// OTel middleware should be early so spans are created before our logger runs.
	r.Use(otelgin.Middleware(getServiceName()))
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(d.Logger))

	// Basic safety middleware
	r.Use(gin.Recovery())

	h := handlers.New(q, d.Logger)

	r.GET("/healthz", h.Health)
	r.GET("/v1/hello", h.Hello)
	r.GET("/v1/todos", h.ListTodos)
	r.POST("/v1/todos", h.CreateTodo)

	// Default timeouts are on the server, but per-handler DB timeouts are used too.
	_ = time.Second

	return r
}

func getServiceName() string {
	// Keep consistent with cmd/server/main.go default
	// (otelgin uses this name for span attributes)
	return "gin-sqlc-demo"
}
