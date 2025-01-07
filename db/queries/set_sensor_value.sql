-- name: SetSensorValue :exec
UPDATE sensors
SET value = $1
WHERE id = $2 and token = $3;
