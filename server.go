package main

import (
	//"auth"
	"fmt"
	"net/http"

	"github.com/DawoudTormos/Home_Control_System_Backend/auth"
	"github.com/gin-gonic/gin"
)

func main() {

	auth.LoadJwtKey()

	server := gin.Default()

	// Login route for generating tokens
	server.POST("/login", auth.CheckLogin())

	// Check token validity
	server.GET("/check-token", auth.TokenAuthMiddleware(), auth.NewToken())

	// Protected route
	protected := server.Group("/secure", auth.TokenAuthMiddleware())
	{
		protected.GET("/data", func(c *gin.Context) {
			username := c.GetString("username")
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Welcome %s, here is your secure data.", username)})
		})
	}

	server.Run(":8080")
}
