-- name: SetCameraIndex :exec
update cameras
set index = $1
where id = $2
;