-- name: Find :one
SELECT * FROM users u WHERE u.id = $1;

-- name: FindByEmail :one
SELECT * FROM users u WHERE u.email = $1;

-- name: Store :exec
INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4);
