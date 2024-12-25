-- name: GetUserCredentials :one
SELECT username, hashed_password
FROM users
WHERE username = $1;