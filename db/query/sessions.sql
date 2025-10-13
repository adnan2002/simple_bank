



-- name: CreateSession :one
INSERT INTO sessions (
  id, username, refresh_token, user_agent, client_ip, expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetActiveSessions :many
SELECT 1
FROM sessions
WHERE id = $1
  AND username = $2
  AND refresh_token = $3
  AND is_blocked = FALSE
  AND expires_at > NOW()
  LIMIT 1;


-- name: BlockSession :exec
UPDATE sessions
SET is_blocked = TRUE
WHERE id = $1
RETURNING *;


-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1
LIMIT 1;
