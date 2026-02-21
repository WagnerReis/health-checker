-- name: CreateUser :exec
INSERT INTO users (
    id, 
    name, 
    email, 
    password, 
    refresh_token, 
    created_at, 
    updated_at
) VALUES (
    COALESCE(sqlc.narg(id), gen_random_uuid()),
    sqlc.arg(name),
    sqlc.arg(email),
    sqlc.arg(password),
    sqlc.narg(refresh_token),
    NOW(),
    NOW()
 );

-- name: Update :exec
UPDATE users SET
    name = sqlc.arg(name),
    email = sqlc.arg(email),
    password = sqlc.arg(password),
    refresh_token = sqlc.narg(refresh_token),
    updated_at = NOW()
WHERE id = sqlc.arg(id);

-- name: FindByID :one
SELECT * FROM users WHERE id = sqlc.arg(id);

-- name: FindByEmail :one
SELECT * FROM users WHERE email = sqlc.arg(email);