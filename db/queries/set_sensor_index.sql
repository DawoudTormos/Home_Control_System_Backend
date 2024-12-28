-- name: setSensorIndex :exec
update sensors
set index = $1
where id = $2
;