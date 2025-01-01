-- name: LinkCamera :one
UPDATE cameras
SET name = $1, color = $2, room_id = $3, index = (COALESCE((SELECT MAX(index) FROM cameras WHERE cameras.room_id = $3), 0) + 1)
WHERE cameras.id = $4 AND cameras.token = $5
RETURNING cameras.id;
