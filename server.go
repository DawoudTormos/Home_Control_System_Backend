package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	server := gin.Default()

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	server.POST("/submit", func(c *gin.Context) {
		c.String(200, "received submission")
	})

	server.GET("/", func(c *gin.Context) {
		c.String(200, "This the main page")
	})

	server.NoRoute(func(c *gin.Context) {
		c.String(404, "404 Not Found: The page you are looking for does not exist.")
	})
	server.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
