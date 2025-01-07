-- name: GetSwitchState :one
SELECT value as state
FROM switches
WHERE id = $1 and token = $2;