-- name: AddRoom :exec
    INSERT INTO rooms (name, user_id, index)
    VALUES ($1,
     (SELECT ID  FROM users WHERE username = $2),
      COALESCE((SELECT MAX(index)  FROM rooms WHERE rooms.user_id = user_id),0)+1)
    ;
