package api

import (
	"database/sql"
	"net/http"
	"sync"
	"time"

	"github.com/DawoudTormos/Home_Control_System_Backend/db"
	"github.com/gin-gonic/gin"
)

type linkRequest struct {
	ReqID     int32
	Username  string
	RoomID    int32
	DeviceID  int32
	Name      string
	Color     int64
	IconCode  int32
	Status    bool
	TimeStart int64 //Unix time in go size. Good to work forever with in32 limit
}

var (
	linkingRequests []linkRequest = []linkRequest{}
	requestsMutex   sync.Mutex
)

func StartPeriodicCheck() { //stops request after 10 minutes
	go func() {
		for {
			requestsMutex.Lock()
			for i, req := range linkingRequests {
				if time.Now().Unix()-req.TimeStart > 10*60 {
					linkingRequests = append(linkingRequests[:i], linkingRequests[i+1:]...)
				}
			}
			requestsMutex.Unlock()

			time.Sleep(1 * time.Second)
		}
	}()
}

func CheckDeviceExists(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var device struct {
			ID int32 `json:"ID"`
		}

		if err := c.ShouldBindJSON(&device); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		ctx := c.Request.Context()
		//username := c.GetString("username")

		queries := db.New(dbConn)

		result, err := queries.CheckDeviceExists(ctx, device.ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error in db.", "details": err.Error()})
			return
		} else if result.ID == 0 {
			c.JSON(http.StatusOK, gin.H{"error": "Device does not exist."})
			return
		}

		c.JSON(http.StatusOK, result)

	}
}

func CheckDeviceExistsAndStartLinking(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var device struct {
			ID         int32  `json:"ID"`
			RoomID     int32  `json:"roomId"`
			DeviceName string `json:"deviceName"`
			Color      int64  `json:"color"`
			IconCode   int32  `json:"icon_code"`
			IconFamily string `json:"icon_family"`
			SType      string `json:"SType"`
		}

		if err := c.ShouldBindJSON(&device); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		ctx := c.Request.Context()
		username := c.GetString("username")

		queries := db.New(dbConn)

		result, err := queries.CheckDeviceExists(ctx, device.ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error in db.", "details": err.Error()})
			return
		} else if result.ID == 0 {
			c.JSON(http.StatusOK, gin.H{"error": "Device does not exist."})
		}

		req := linkRequest{
			ReqID:     int32(len(linkingRequests)),
			Username:  username,
			RoomID:    int32(device.RoomID),
			DeviceID:  int32(result.ID),
			Name:      device.DeviceName,
			Color:     device.Color,
			IconCode:  device.IconCode,
			Status:    false,
			TimeStart: time.Now().Unix(),
		}

		if req.IconCode == 0 {
			req.IconCode = 57800
		}

		requestsMutex.Lock()
		linkingRequests = append(linkingRequests, req)
		requestsMutex.Unlock()

		response := map[string]interface{}{
			"device": result,
			"ReqID":  req.ReqID,
		}

		c.JSON(http.StatusOK, response)

	}
}

func CheckDeviceRequestState(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var check struct {
			ReqID int32 `json:"reqID"`
		}

		if err := c.ShouldBindJSON(&check); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		//ctx := c.Request.Context()
		username := c.GetString("username")

		var isFound bool = false
		var status bool = false

		requestsMutex.Lock()
		for i, req := range linkingRequests {
			print(i, "\n")
			//i++
			if check.ReqID == req.ReqID && req.Username == username {
				print(req.DeviceID, " device id: found\n")
				status = req.Status
				isFound = true
				break
			}

		}
		requestsMutex.Unlock()

		if !isFound {
			c.JSON(http.StatusOK, gin.H{"error": "Request not found."})
			return
		} else if !status {
			c.JSON(http.StatusOK, gin.H{"error": "none", "status": "pending"})
			return
		} else if status {
			c.JSON(http.StatusOK, gin.H{"error": "none", "status": "success"})
			return
		}

	}
}

func AcceptDeviceLinking(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var device struct {
			ID         int32  `json:"ID"`
			DeviceType string `json:"type"`
			Token      string `json:"token"`
		}

		if err := c.ShouldBindJSON(&device); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		ctx := c.Request.Context()
		//username := c.GetString("username")

		queries := db.New(dbConn)

		var request linkRequest

		var isFound bool = false
		var indexFound int = 0

		requestsMutex.Lock()
		for i, req := range linkingRequests {
			print(i, "\n")
			//i++
			if device.ID == req.DeviceID && !req.Status {
				print(req.DeviceID, " device id: found\n")
				request = req
				indexFound = i
				isFound = true
				break
			}

		}
		requestsMutex.Unlock()

		if !isFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "device has no link request"})
			return
		} else {
			print(request.DeviceID, " device id: found\n")
		}

		if device.DeviceType == "switch" {

			var linkedDevice = db.LinkSwitchParams{
				Name:     request.Name,
				Color:    request.Color,
				IconCode: request.IconCode,
				RoomID:   request.RoomID,
				ID:       request.DeviceID,
				Token:    sql.NullString{String: device.Token, Valid: true},
			}
			print(linkedDevice.Color, "\n")
			print(linkedDevice.Name, "\n")
			//print(linkedDevice.Token, "\n")
			result, err := queries.LinkSwitch(ctx, linkedDevice)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error in db.", "details": err.Error()})
				return
			} else if result != request.DeviceID {
				c.JSON(http.StatusBadRequest, gin.H{"error": "No row updated"})
				return

			}

		} else if device.DeviceType == "sensor" {

			var linkedDevice = db.LinkSensorParams{
				Name:   request.Name,
				Color:  request.Color,
				RoomID: request.RoomID,
				ID:     request.DeviceID,
				Token:  sql.NullString{String: device.Token, Valid: true},
			}
			print(linkedDevice.Color, "\n")
			print(linkedDevice.Name, "\n")
			//print(linkedDevice.Token, "\n")
			result, err := queries.LinkSensor(ctx, linkedDevice)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error in db.", "details": err.Error()})
				return
			} else if result != request.DeviceID {
				c.JSON(http.StatusBadRequest, gin.H{"error": "No row updated"})
				return

			}

		} else if device.DeviceType == "camera" {

			var linkedDevice = db.LinkCameraParams{
				Name:   request.Name,
				Color:  request.Color,
				RoomID: request.RoomID,
				ID:     request.DeviceID,
				Token:  sql.NullString{String: device.Token, Valid: true},
			}
			print(linkedDevice.Color, "\n")
			print(linkedDevice.Name, "\n")
			//print(linkedDevice.Token, "\n")
			result, err := queries.LinkCamera(ctx, linkedDevice)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error in db.", "details": err.Error()})
				return
			} else if result != request.DeviceID {
				c.JSON(http.StatusBadRequest, gin.H{"error": "No row updated"})
				return

			}

		} else {
			c.JSON(http.StatusOK, gin.H{"error": "invalid device or bad request."})
			return

		}

		requestsMutex.Lock()
		linkingRequests[indexFound].Status = true
		requestsMutex.Unlock()

		c.JSON(http.StatusOK, gin.H{"result": "success"})
	}
}
