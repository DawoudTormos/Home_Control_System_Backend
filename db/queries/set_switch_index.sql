-- name: setSwitchIndex :exec
update switches
set index = $1
where id = $2
;