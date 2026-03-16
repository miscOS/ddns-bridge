package auth

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/miscOS/ddns-bridge/models"
	"golang.org/x/crypto/bcrypt"
)

type LocalAuthenticator struct {
	id       string
	secret   string
	register string
}

func NewLocalAuthenticator(secret string, register string) (*LocalAuthenticator, error) {

	if len(secret) < 32 {
		log.Default().Println("Warning: Local JWT secret is less than 32 characters long, generating a random secret")
		secret = generateKey()
	}

	if register == "" {
		log.Default().Println("Warning: Local JWT register key is not set, registration is open to anyone")
	}

	return &LocalAuthenticator{
		id:       "local",
		secret:   secret,
		register: register,
	}, nil
}

func generateKey() string {

	random := randomString(32)
	return string(hmac.New(sha512.New, []byte(random)).Sum(nil))
}

func randomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (a *LocalAuthenticator) Authenticate(c *gin.Context) (*models.User, error) {
	h := c.GetHeader("Authorization")

	if !strings.HasPrefix(h, "Bearer ") {
		return nil, nil
	}

	raw := strings.TrimPrefix(h, "Bearer ")

	// Skip if instance is not responsible
	if ok, err := a.isResponsible(raw); err != nil || !ok {
		return nil, err
	}

	claims, err := a.verifyToken(raw)
	if err != nil {
		return nil, err
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid token: missing subject")
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return nil, errors.New("invalid token: subject is not a valid UUID")
	}

	return &models.User{
		ID: userID,
	}, nil
}

func (a *LocalAuthenticator) GetID() string {
	return a.id
}

func (a *LocalAuthenticator) AllowRegistration(register string) bool {
	return a.register == register
}

func (a *LocalAuthenticator) GenerateToken(userID string) (string, error) {

	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": a.id,                                 // Issuer
		"sub": userID,                               // User identifier
		"exp": time.Now().Add(4 * time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                    // Issued at
	})

	// Sign the token with the secret key
	if token, err := claims.SignedString([]byte(a.secret)); err != nil {
		return "", err
	} else {
		return token, nil
	}
}

func (a *LocalAuthenticator) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (a *LocalAuthenticator) VerifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (a *LocalAuthenticator) isResponsible(token string) (bool, error) {

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false, errors.New("provided token is not a valid JWT")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	var jwt struct {
		Iss string `json:"iss"`
	}
	if err := json.Unmarshal(payload, &jwt); err != nil {
		return false, err
	}

	return a.id == jwt.Iss, nil
}

func (a *LocalAuthenticator) parseToken(tokenString string) (*jwt.Token, error) {

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.secret), nil
	})
}

func (a *LocalAuthenticator) verifyToken(tokenString string) (jwt.MapClaims, error) {

	token, err := a.parseToken(tokenString)

	if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
		return mapClaims, err
	} else {
		return nil, fmt.Errorf("unexpected claims type: %T", token.Claims)
	}
}
