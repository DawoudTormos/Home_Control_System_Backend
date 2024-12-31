-- name: CheckDeviceExists :one

SELECT sensors.id, 'sensor' as device_type
FROM sensors 
WHERE sensors.id = $1
UNION ALL
SELECT cameras.id, 'camera' as device_type
FROM cameras 
WHERE cameras.id = $1
UNION ALL
SELECT switches.id, 'switch' as device_type
FROM switches 
WHERE switches.id = $1
LIMIT 1;