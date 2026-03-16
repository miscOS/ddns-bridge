package auth

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/miscOS/ddns-bridge/models"
)

type Authenticator interface {
	GetID() string
	Authenticate(c *gin.Context) (*models.User, error)
}

type AuthenticatorRegistry struct {
	items map[string]Authenticator
	mu    sync.RWMutex
}

var registry *AuthenticatorRegistry = &AuthenticatorRegistry{
	items: make(map[string]Authenticator),
	mu:    sync.RWMutex{},
}

func Register(item Authenticator) error {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	registry.items[item.GetID()] = item
	return nil
}

func GetAuthenticators() []Authenticator {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	items := make([]Authenticator, 0, len(registry.items))

	for _, item := range registry.items {
		items = append(items, item)
	}

	return items
}

func GetLocalAuthenticator() (*LocalAuthenticator, error) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	item, ok := registry.items["local"]
	if !ok {
		return nil, fmt.Errorf("local authenticator not found")
	}
	return item.(*LocalAuthenticator), nil
}

func GetMiddleware() gin.HandlerFunc {

	authenticators := GetAuthenticators()

	return func(c *gin.Context) {
		for _, authenticator := range authenticators {
			u, err := authenticator.Authenticate(c)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
			if u != nil {
				c.Set("user", u)
				c.Next()
				return
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to authenticate user"})
		c.Abort()
	}
}
