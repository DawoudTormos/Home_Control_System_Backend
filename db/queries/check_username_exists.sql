-- name: CheckUsernameExists :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE username = $1
) AS exists;
