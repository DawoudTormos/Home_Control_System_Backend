-- name: GetRooms :many
SELECT rooms.id, rooms.name 
FROM rooms
JOIN users
ON users.id = rooms.user_id
WHERE users.username = $1;