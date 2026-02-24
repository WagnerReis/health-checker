-- name: CreateUser :exec
INSERT INTO users (
    id, 
    name, 
    email, 
    password, 
    created_at, 
    updated_at
) VALUES (
    COALESCE(sqlc.narg(id), gen_random_uuid()),
    sqlc.arg(name),
    sqlc.arg(email),
    sqlc.arg(password),
    NOW(),
    NOW()
 );

-- name: Update :exec
UPDATE users SET
    name = sqlc.arg(name),
    email = sqlc.arg(email),
    password = sqlc.arg(password),
    updated_at = NOW()
WHERE id = sqlc.arg(id);

-- name: FindByID :one
SELECT * FROM users WHERE id = sqlc.arg(id);

-- name: FindByEmail :one
SELECT * FROM users WHERE email = sqlc.arg(email);

-- Refresh Token Queries

-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (
    id, 
    user_id, 
    token_hash,
    expires_at, 
    revoked, 
    created_at, 
    updated_at
) VALUES (
    COALESCE(sqlc.narg(id), gen_random_uuid()),
    sqlc.arg(user_id),
    sqlc.arg(token_hash),
    sqlc.arg(expires_at),
    false,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);

-- name: FindByTokenHash :one
SELECT * FROM refresh_tokens WHERE token_hash = sqlc.arg(tokenHash);

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens SET 
    revoked = true,
    updated_at = CURRENT_TIMESTAMP
WHERE token_hash = sqlc.arg(token_hash)
RETURNING 1;

-- name: RevokeAllByUser :exec
UPDATE refresh_tokens SET 
    revoked = true,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = sqlc.arg(userID);

-- name: DeleteAllExpired :exec
DELETE FROM refresh_tokens WHERE expires_at < CURRENT_TIMESTAMP;

-- Monitor Queries

-- name: CreateMonitor :exec
INSERT INTO monitors (
    id, 
    user_id, 
    name, 
    url, 
    method, 
    headers, 
    body, 
    interval, 
    expected_status_code,
    timeout,
    created_at, 
    updated_at
) VALUES (
    COALESCE(sqlc.narg(id), gen_random_uuid()),
    sqlc.arg(user_id),
    sqlc.arg(name),
    sqlc.arg(url),
    sqlc.arg(method),
    sqlc.narg(headers),
    sqlc.narg(body),
    sqlc.arg(interval),
    sqlc.narg(expected_status_code),
    sqlc.arg(timeout),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);
