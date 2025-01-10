package helpers

import (
	"crypto/hmac"
	"crypto/sha512"

	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("unsecure-default-key")

func RandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenerateSecret() {

	randomString := RandomString(32)
	hmacKey := hmac.New(sha512.New, []byte(randomString))

	secretKey = hmacKey.Sum(nil)
}

func GetSecret() []byte {
	return secretKey
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func VerifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
