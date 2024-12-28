-- name: SetRoomIndex :exec
UPDATE rooms
SET index = $1
WHERE rooms.id = $2 
  AND user_id IN (
    SELECT users.id 
    FROM users 
    WHERE users.username = $3
  );