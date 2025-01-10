package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/miscOS/ddns-bridge/helpers"
)

// UserAuthenticate is a middleware that checks if the user is authenticated
func UserAuthenticate(c *gin.Context) {

	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		c.Abort()
		return
	}

	token = strings.Replace(token, "Bearer ", "", 1)

	claims, err := helpers.VerifyTokenAndMapClaims(token, helpers.GetSecret())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	if username, ok := claims["sub"]; ok {
		c.Set("username", username)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}

	c.Next()
}
