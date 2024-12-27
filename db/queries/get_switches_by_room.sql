-- name: GetswitchesByRoom :many
SELECT id, name, color, icon_code, icon_family, type, value ,index
FROM switches
WHERE room_id = $1;