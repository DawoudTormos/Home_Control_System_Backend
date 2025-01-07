package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/DawoudTormos/Home_Control_System_Backend/db"
	ginmethods "github.com/DawoudTormos/Home_Control_System_Backend/ginMethods"
	"github.com/gin-gonic/gin"
)

func SendSensorData(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		bodyBytes, _ := ginmethods.ReadAndPrintBody(c)

		bodyString := strings.TrimSpace(string(bodyBytes))

		parts := strings.Split(bodyString, ",")
		if len(parts) != 3 {
			c.String(http.StatusBadRequest, "Invalid format. Expected: ID,token")
			println("Invalid format. Expected: ID,token")
			return
		}

		id := strings.TrimSpace(parts[0])
		token := strings.TrimSpace(parts[1])
		value := strings.TrimSpace(parts[2])

		deviceID, err := strconv.Atoi(id)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid ID format")
			println("Invalid ID format")
			return
		}

		valueS, err := strconv.ParseFloat(value, 32)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid value format")
			println("Invalid value format")
			return
		}
		valueS = valueS * 100

		ctx := c.Request.Context()

		queries := db.New(dbConn)

		err = queries.SetSensorValue(ctx, db.SetSensorValueParams{
			ID:    int32(deviceID),
			Token: sql.NullString{String: token, Valid: token != ""},
			Value: int32(valueS),
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error in db.", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, "success")

	}
}

func GetSwitchStatus(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		bodyBytes, _ := ginmethods.ReadAndPrintBody(c)

		parts := strings.Split(string(bodyBytes), ",")
		if len(parts) != 2 {
			c.String(http.StatusBadRequest, "Invalid format. Expected: ID,token")
			return
		}

		id := strings.TrimSpace(parts[0])
		token := strings.TrimSpace(parts[1])

		deviceID, err := strconv.Atoi(id)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid ID format")
			return
		}

		ctx := c.Request.Context()

		queries := db.New(dbConn)

		result, err := queries.GetSwitchState(ctx, db.GetSwitchStateParams{
			ID:    int32(deviceID),
			Token: sql.NullString{String: token, Valid: token != ""},
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error in db.", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)

	}
}
