package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"gin-sqlc-demo/internal/db/sqlc"
)

type Handlers struct {
	q      *sqlc.Queries
	logger *slog.Logger
}

func New(q *sqlc.Queries, logger *slog.Logger) *Handlers {
	return &Handlers{q: q, logger: logger}
}

func (h *Handlers) ListTodos(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	items, err := h.q.ListTodos(ctx)
	if err != nil {
		h.logger.Error("list todos failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	c.JSON(http.StatusOK, items)
}

type createTodoReq struct {
	Title string `json:"title"`
}

func (h *Handlers) CreateTodo(c *gin.Context) {
	var req createTodoReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	todo, err := h.q.CreateTodo(ctx, sqlc.CreateTodoParams{Title: req.Title})
	if err != nil {
		h.logger.Error("create todo failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	c.JSON(http.StatusCreated, todo)
}
