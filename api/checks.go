package api

import (
	"database/sql"
	"net/http"

	"github.com/DawoudTormos/Home_Control_System_Backend/db"
	"github.com/gin-gonic/gin"
)

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
		}

		c.JSON(http.StatusOK, result)

	}
}
