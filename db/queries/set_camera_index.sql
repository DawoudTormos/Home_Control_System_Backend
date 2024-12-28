-- name: SetCameraIndex :exec
UPDATE cameras
SET index = $1
WHERE cameras.id = $2 
  AND room_id IN  (
    SELECT rooms.id 
    FROM rooms 
    WHERE rooms.user_id  IN  (
            SELECT users.id 
            FROM users 
            WHERE users.username = $3
        )
  );