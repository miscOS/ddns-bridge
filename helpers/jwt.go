package helpers

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Function to create JWT tokens with claims
func CreateToken(username string, secretKey []byte) (string, error) {

	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "ddns-bridge",                        // Issuer
		"sub": username,                             // User identifier
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

func ParseToken(tokenString string, secretKey []byte) (*jwt.Token, error) {

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key
		return secretKey, nil
	})
}

func VerifyTokenAndMapClaims(tokenString string, secretKey []byte) (jwt.MapClaims, error) {

	// Parse the token
	token, err := ParseToken(tokenString, secretKey)

	// Return claims and error
	if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
		return mapClaims, err
	} else {
		return nil, fmt.Errorf("unexpected claims type: %T", token.Claims)
	}
}
