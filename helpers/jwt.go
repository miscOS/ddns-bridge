package helpers

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Add a new global variable for the secret key
var secretKey = []byte("your-secret-key")

// Function to create JWT tokens with claims
func CreateToken(uid uint) (string, error) {

	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "ddns-bridge",                        // Issuer
		"uid": uid,                                  // User identifier
		"exp": time.Now().Add(4 * time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                    // Issued at
	})

	// Sign the token with the secret key
	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	// Print information about the created token
	fmt.Printf("Token claims added: %+v\n", claims)
	return tokenString, nil
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key
		return secretKey, nil
	})

	// Return claims and error
	if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
		return mapClaims, err
	} else {
		return nil, fmt.Errorf("unexpected claims type: %T", token.Claims)
	}
}
