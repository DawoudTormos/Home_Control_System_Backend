-- name: ValidateCredentials :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND hashed_password = $2);
