package api

import (
	"database/sql"
	"net/http"

	"github.com/DawoudTormos/Home_Control_System_Backend/db"
	"github.com/gin-gonic/gin"
)

type Device interface {
}

func GetRooms(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()
		username := c.GetString("username")

		queries := db.New(dbConn)

		rooms, err := queries.GetRooms(ctx, username)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rooms", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, rooms)

	}
}

func GetDevices(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()
		username := c.GetString("username")

		queries := db.New(dbConn)

		rooms, err := queries.GetRooms(ctx, username)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "details": err.Error()})
			return
		}

		devices := []Device{}
		devicesByRoom := make(map[string][]Device)

		for _, value := range rooms {
			devices = []Device{}
			//fmt.Printf("Index: %d, Value: %d\n", index, value)
			switches, _ := queries.GetswitchesByRoom(ctx, value.ID)
			sensors, _ := queries.GetsensorsByRoom(ctx, value.ID)
			cameras, _ := queries.GetcamerasByRoom(ctx, value.ID)

			if switches != nil {
				for _, value := range switches {
					value2 := toMapSwitch(value)
					devices = append(devices, value2)
				}
			}

			if sensors != nil {
				for _, value := range sensors {
					value2 := toMapSensor(value)
					devices = append(devices, value2)
				}
			}

			if cameras != nil {
				for _, value := range cameras {
					value2 := toMapCamera(value)
					devices = append(devices, value2)
				}
			}

			devicesByRoom[value.Name] = devices
			//print(switches[0].Color)
			//print(devices[0].Color)
		}

		c.JSON(http.StatusOK, devicesByRoom)

	}
}

func toMapSwitch(row db.GetswitchesByRoomRow) map[string]interface{} {
	return map[string]interface{}{
		"ID":         row.ID,
		"Name":       row.Name,
		"Color":      row.Color,
		"IconCode":   row.IconCode,
		"IconFamily": row.IconFamily,
		"Type":       row.Type,
		"Value":      row.Value,
		"Index":      row.Index,
		"SType":      "switch",
	}
}

func toMapCamera(row db.GetcamerasByRoomRow) map[string]interface{} {
	return map[string]interface{}{
		"ID":    row.ID,
		"Name":  row.Name,
		"Color": row.Color,
		"Value": row.Value,
		"Index": row.Index,
		"SType": "camera",
	}
}

func toMapSensor(row db.GetsensorsByRoomRow) map[string]interface{} {
	return map[string]interface{}{
		"ID":    row.ID,
		"Name":  row.Name,
		"Color": row.Color,
		"Type":  row.Type,
		"Value": row.Value,
		"Index": row.Index,
		"SType": "sensor",
	}
}
