package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
)

var JwtKey string

func LoadJwtKey() {
	JwtKey = os.Getenv("HCS_JWT_SECRET_KEY")

	if JwtKey == "" {
		fmt.Println("Error: HCS_JWT_SECRET_KEY is not set")
		return
	}
	print("------\n")
	print("Secret_KEY", os.Getenv("HCS_JWT_SECRET_KEY"), "\n")
	print("------\n")
}

// ValidateCredentials validates username and password
// Implement this function to connect to PostgreSQL
func validateCredentials(username string, pass string) bool {
	// Placeholder logic; replace with actual database query
	return username == "admin" && pass == "password"
}

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
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
	}
}

// GenerateToken generates a new JWT token for a user
func GenerateToken(username string) (string, error) {

	expirationTime := time.Now().Add(72 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

// TokenAuthMiddleware ensures token is valid as a middleware for most routes that needs authentication
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
			return JwtKey, nil
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

func NewToken() gin.HandlerFunc { // Issue a new token with an extended expiration time
	return func(c *gin.Context) {
		username := c.GetString("username")

		newToken, err := GenerateToken(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":   "Token is valid",
			"new_token": newToken,
			"username":  username,
		})

	}
}
