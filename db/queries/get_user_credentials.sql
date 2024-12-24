-- name: GetUserCredentials :one
SELECT salt, hashed_password
FROM users
WHERE username = $1;