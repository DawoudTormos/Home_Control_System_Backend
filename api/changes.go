package api

import (
	"database/sql"
	"net/http"

	"github.com/DawoudTormos/Home_Control_System_Backend/db"
	"github.com/gin-gonic/gin"
)

func SetRoomIndex(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()
		username := c.GetString("username")

		queries := db.New(dbConn)

		rooms, err := queries.SetRoomIndex(ctx, username)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rooms", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"rooms": rooms})

	}
}

func GetDevicess(dbConn *sql.DB) gin.HandlerFunc {
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
			//fmt.Printf("Index: %d, Value: %d\n", index, value)
			switches, _ := queries.GetswitchesByRoom(ctx, value.ID)
			sensors, _ := queries.GetsensorsByRoom(ctx, value.ID)
			cameras, _ := queries.GetcamerasByRoom(ctx, value.ID)

			if switches != nil {
				for _, value := range switches {

					devices = append(devices, value)
				}
			}

			if sensors != nil {
				for _, value := range sensors {

					devices = append(devices, value)
				}
			}

			if cameras != nil {
				for _, value := range cameras {

					devices = append(devices, value)
				}
			}

			devicesByRoom[value.Name] = devices
			//print(switches[0].Color)
			//print(devices[0].Color)
		}

		c.JSON(http.StatusOK, devicesByRoom)

	}
}
