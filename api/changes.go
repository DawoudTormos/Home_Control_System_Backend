package api

import (
	"database/sql"
	"net/http"

	"github.com/DawoudTormos/Home_Control_System_Backend/db"
	ginmethods "github.com/DawoudTormos/Home_Control_System_Backend/ginMethods"
	"github.com/gin-gonic/gin"
)

func SetIndexes(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var devices []struct {
			ID    int    `json:"Id"`
			Type  string `json:"Type"`
			Index int    `json:"Index"`
		}

		if err := c.ShouldBindJSON(&devices); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		print(devices)

		ctx := c.Request.Context()
		username := c.GetString("username")

		queries := db.New(dbConn)

		for _, value := range devices {
			if value.Type == "room" {
				var parms db.SetRoomIndexParams
				parms.ID = int32(value.ID)
				parms.Index = int32(value.Index)
				parms.Username = username
				err := queries.SetRoomIndex(ctx, parms)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed change index in a room", "details": err.Error()})
					return
				}
			}

			if value.Type == "switch" {
				var parms db.SetSwitchIndexParams
				parms.ID = int32(value.ID)
				parms.Index = int32(value.Index)
				parms.Username = username
				err := queries.SetSwitchIndex(ctx, parms)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed change index in a room", "details": err.Error()})
					return
				}
			}

			if value.Type == "camera" {
				var parms db.SetCameraIndexParams
				parms.ID = int32(value.ID)
				parms.Index = int32(value.Index)
				parms.Username = username
				err := queries.SetCameraIndex(ctx, parms)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed change index in a room", "details": err.Error()})
					return
				}
			}

			if value.Type == "sensor" {
				var parms db.SetSensorIndexParams
				parms.ID = int32(value.ID)
				parms.Index = int32(value.Index)
				parms.Username = username
				err := queries.SetSensorIndex(ctx, parms)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed change index in a room", "details": err.Error()})
					return
				}
			}
		}

		/*if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rooms", "details": err.Error()})
			return
		}*/

		c.JSON(http.StatusOK, gin.H{"result": "success"})

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

func SetSwitchValue(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var device struct {
			ID    int `json:"Id"`
			Value int `json:"Value"`
		}

		if err := c.ShouldBindJSON(&device); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		ctx := c.Request.Context()
		username := c.GetString("username")

		queries := db.New(dbConn)

		var parms db.SetSwitchValueParams

		parms.ID = int32(device.ID)
		parms.Value = int16(device.Value)
		parms.Username = username
		ginmethods.ReadAndPrintBody(c)
		print(parms.Value, "\n")

		err := queries.SetSwitchValue(ctx, parms)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error in db in updating the switch's value."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": "success"})

	}
}

func AddRoom(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var room struct {
			Name string `json:"name"`
		}

		if err := c.ShouldBindJSON(&room); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		ctx := c.Request.Context()
		username := c.GetString("username")

		queries := db.New(dbConn)

		var parms db.AddRoomParams

		parms.Name = room.Name
		parms.Username = username

		err := queries.AddRoom(ctx, parms)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error in db in add a new room."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": "success"})

	}
}
