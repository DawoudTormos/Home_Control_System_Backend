-- name: GetcamerasByRoom :many
SELECT id, name, color, value, index
FROM cameras
WHERE room_id = $1 ;