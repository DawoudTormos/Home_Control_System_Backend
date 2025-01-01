// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: link_switch.sql

package db

import (
	"context"
	"database/sql"
)

const linkSwitch = `-- name: LinkSwitch :one
UPDATE switches
SET name = $1, color = $2, icon_code = $3, room_id = $4, index = (COALESCE((SELECT MAX(index) FROM switches WHERE switches.room_id = $4), 0) + 1)
WHERE switches.id = $5 AND switches.token = $6
RETURNING switches.id
`

type LinkSwitchParams struct {
	Name     string
	Color    int64
	IconCode int32
	RoomID   int32
	ID       int32
	Token    sql.NullString
}

func (q *Queries) LinkSwitch(ctx context.Context, arg LinkSwitchParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, linkSwitch,
		arg.Name,
		arg.Color,
		arg.IconCode,
		arg.RoomID,
		arg.ID,
		arg.Token,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}
