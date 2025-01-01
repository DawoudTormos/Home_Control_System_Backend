-- name: LinkSwitch :one
UPDATE switches
SET name = $1, color = $2, icon_code = $3, room_id = $4, index = (COALESCE((SELECT MAX(index) FROM switches WHERE switches.room_id = $4), 0) + 1)
WHERE switches.id = $5 AND switches.token = $6
RETURNING switches.id;
