-- name: GetsensorsByRoom :many
SELECT id, name, color, type, value ,index
FROM sensors
WHERE room_id = $1;