-- name: GetsensorsByRoom :many
SELECT id, name, color, type_id as type, value ,index
FROM sensors
WHERE room_id = $1;