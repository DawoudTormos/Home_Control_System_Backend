package main

import (
	//"auth"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/DawoudTormos/Home_Control_System_Backend/api"
	"github.com/DawoudTormos/Home_Control_System_Backend/auth"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
)

func main() {

	connStr := "postgres://postgres:admin@localhost:5432/HCS?sslmode=disable"
	dbConn, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer dbConn.Close()

	auth.LoadJwtKey()

	server := gin.Default()

	// Login route for generating tokens
	server.POST("/login", auth.CheckLogin(dbConn))

	// Sign Up route
	server.POST("/signup", auth.SignUp(dbConn))

	// Check token validity
	server.POST("/check-token", auth.TokenAuthMiddleware(), auth.NewToken())

	// Protected route
	protected := server.Group("/secure", auth.TokenAuthMiddleware())
	{
		protected.GET("/data", func(c *gin.Context) {
			username := c.GetString("username")
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Welcome %s, here is your secure data.", username)})
		})

		protected.GET("/getRooms", api.GetRooms(dbConn))
		protected.GET("/getDevices", api.GetDevices(dbConn))
		protected.GET("/setIndexes", api.SetIndexes(dbConn))

	}

	server.Run(":8080")
}
