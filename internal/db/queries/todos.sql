-- name: ListTodos :many
SELECT id, title, done, created_at
FROM todos
ORDER BY id DESC;

-- name: CreateTodo :one
INSERT INTO todos (title)
VALUES ($1)
RETURNING id, title, done, created_at;
