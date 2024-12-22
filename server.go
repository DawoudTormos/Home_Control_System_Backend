package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey string

// Claims struct for JWT
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// ValidateCredentials validates username and password
// Implement this function to connect to PostgreSQL
func validateCredentials(username string, pass string) bool {
	// Placeholder logic; replace with actual database query
	return username == "admin" && pass == "password"
}

// GenerateToken generates a new JWT token for a user
func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// TokenAuthMiddleware ensures token is valid (not wrong or expired)
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is missing"})
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		print()

		c.Set("username", claims.Username)
		c.Next()
	}
}

func main() {

	jwtKey = os.Getenv("HCS_JWT_SECRET_KEY")

	if jwtKey == "" {
		fmt.Println("Error: HCS_JWT_SECRET_KEY is not set")
		return
	}
	print("------\n")
	print("Secret_KEY", os.Getenv("HCS_JWT_SECRET_KEY"), "\n")
	print("------\n")

	server := gin.Default()

	// Login route for generating tokens
	server.POST("/login", func(c *gin.Context) {
		var credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if !validateCredentials(credentials.Username, credentials.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, err := GenerateToken(credentials.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	// Check token validity
	server.GET("/check-token", TokenAuthMiddleware(), func(c *gin.Context) {
		username := c.GetString("username")
		c.JSON(http.StatusOK, gin.H{"message": "Token is valid", "username": username})
	})

	// Protected route
	protected := server.Group("/secure", TokenAuthMiddleware())
	{
		protected.GET("/data", func(c *gin.Context) {
			username := c.GetString("username")
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Welcome %s, here is your secure data.", username)})
		})
	}

	server.Run(":8080")
}
