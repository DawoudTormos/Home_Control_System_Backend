-- name: LinkSensor :one
UPDATE sensors
SET name = $1, color = $2, room_id = $3, index = (COALESCE((SELECT MAX(index) FROM sensors WHERE sensors.room_id = $3), 0) + 1)
WHERE sensors.id = $4 AND sensors.token = $5
RETURNING sensors.id;