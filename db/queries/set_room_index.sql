-- name: setRoomIndex :exec
update rooms
set index = $1
where id = $2
;