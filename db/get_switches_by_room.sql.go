// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: get_switches_by_room.sql

package db

import (
	"context"
)

const getswitchesByRoom = `-- name: GetswitchesByRoom :many
SELECT id, name, color, icon_code, icon_family, type, value ,index
FROM switches
WHERE room_id = $1
`

type GetswitchesByRoomRow struct {
	ID         int32
	Name       string
	Color      int64
	IconCode   int32
	IconFamily string
	Type       int32
	Value      int16
	Index      int32
}

func (q *Queries) GetswitchesByRoom(ctx context.Context, roomID int32) ([]GetswitchesByRoomRow, error) {
	rows, err := q.db.QueryContext(ctx, getswitchesByRoom, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetswitchesByRoomRow
	for rows.Next() {
		var i GetswitchesByRoomRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Color,
			&i.IconCode,
			&i.IconFamily,
			&i.Type,
			&i.Value,
			&i.Index,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
