package auth

import (
	"bytes"
	//"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/DawoudTormos/Home_Control_System_Backend/db"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
)

var JwtKey []byte

func LoadJwtKey() {
	JwtKey = []byte(os.Getenv("HCS_JWT_SECRET_KEY"))

	print("\n------\n")

	if bytes.Equal(JwtKey, []byte{}) {
		fmt.Println("Error: HCS_JWT_SECRET_KEY is not set")
		print("------\n")
		os.Exit(1) // Exit with a non-zero status (indicates an error)
	}
	print("Secret_KEY: ", os.Getenv("HCS_JWT_SECRET_KEY"), "\n")
	print("------\n\n")
}

// ValidateCredentials validates username and password
// Implement this function to connect to PostgreSQL
func validateCredentials(c *gin.Context, dbConn *sql.DB, username, password string) bool {

	ctx := c.Request.Context()

	queries := db.New(dbConn)

	// Fetch the salt and hashed password for the user
	creds, err := queries.GetUserCredentials(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Username not found")
		} else {
			log.Println("Database query error:", err)
		}
		return false
	}

	// Hash the salted password

	/*hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	log.Println(string(hash)) */

	err = bcrypt.CompareHashAndPassword([]byte(creds.HashedPassword), []byte(password))
	//saltedHash := hex.EncodeToString(hash[:])
	//println(saltedHash)

	if err != nil {
		log.Println("Invalid password")
		log.Println(" message: ", err.Error())
		return false
	}

	return true
}

/*
func validateCredentials(username string, pass string) bool {
	// Placeholder logic; replace with actual database query
	return username == "admin" && pass == "password"
}*/

func CheckLogin(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if !validateCredentials(c, dbConn, credentials.Username, credentials.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, err := GenerateToken(credentials.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token", "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func SignUp(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Email    string `json:"email"`
		}

		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if len(credentials.Email) >= 150 || len(credentials.Username) >= 100 || len(credentials.Password) >= 72 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error":     "One of the entered inputs is iggere than allowed value",
				"errorCode": "01",
			})
			return
		}

		//Hashing the password with bycrypt
		hashVal, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error":     "error Creating a new account",
				"errorCode": "0",
			})
			log.Println("eror while hashing.\n message: ", err.Error())
			return
		}

		ctx := c.Request.Context()

		queries := db.New(dbConn)

		usernameUsed, err := queries.CheckUsernameExists(ctx, credentials.Username)
		if usernameUsed {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error":     "Username used already!",
				"errorCode": "1",
			})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error":     "error Creating a new account",
				"errorCode": "11",
			})
			return

		}

		emailUsed, err := queries.CheckEmailExists(ctx, credentials.Email)
		if emailUsed {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error":     "Email used already!",
				"errorCode": "2",
			})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error":     "error Creating a new account",
				"errorCode": "21",
			})
			return

		}

		_, err = queries.AddUser(ctx, db.AddUserParams{
			Username:       credentials.Username,
			Email:          credentials.Email,
			HashedPassword: string(hashVal),
		},
		)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error":     "error creating an account",
				"errorCode": "3",
			})
			return
		}

		token, err := GenerateToken(credentials.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error":     "error returning a token",
				"errorCode": "4",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "account created",
			"token":    token,
			"username": credentials.Username,
		})

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
