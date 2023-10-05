-- name: Find :one
SELECT * FROM users u WHERE u.id = $1;

-- name: FindByEmail :one
SELECT * FROM users u WHERE u.email = $1;

-- name: Store :exec
INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4);

-- name: Update :exec
UPDATE users u SET name = $2, email = $3, password = $4 WHERE u.id = $1;
