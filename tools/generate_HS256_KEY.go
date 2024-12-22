package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

// Function to generate a random key for JWT signing
func generateJWTKey() []byte {
	key := make([]byte, 32) // 256 bits
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal("Error generating key:", err)
	}
	return key
}

func main() {
	// Generate the secret key
	jwtKey := generateJWTKey()

	// Print the base64 encoded version of the key
	encodedKey := base64.StdEncoding.EncodeToString(jwtKey)
	log.Println("Generated JWT secret key:", encodedKey)

	// Use the jwtKey in your application
}
