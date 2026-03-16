package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/miscOS/ddns-bridge/models"
)

type OIDCAuthenticator struct {
	id       string
	verifier *oidc.IDTokenVerifier
}

func NewOIDCAuthenticator(issuer, clientID string) (*OIDCAuthenticator, error) {
	provider, err := oidc.NewProvider(context.Background(), issuer)
	if err != nil {
		return nil, err
	}

	return &OIDCAuthenticator{
		id: issuer,
		verifier: provider.Verifier(&oidc.Config{
			ClientID: clientID,
		}),
	}, nil
}

func (a *OIDCAuthenticator) Authenticate(c *gin.Context) (*models.User, error) {

	h := c.GetHeader("Authorization")
	if !strings.HasPrefix(h, "Bearer ") {
		return nil, nil
	}

	raw := strings.TrimPrefix(h, "Bearer ")

	// Skip if instance is not responsible
	if ok, err := a.isResponsible(raw); err != nil || !ok {
		return nil, err
	}

	token, err := a.verifier.Verify(c.Request.Context(), raw)
	if err != nil {
		return nil, err
	}

	var claims struct {
		Sub string `json:"sub"`
	}
	_ = token.Claims(&claims)

	userID, err := uuid.Parse(claims.Sub)
	if err != nil {
		return nil, errors.New("invalid token: subject is not a valid UUID")
	}

	return &models.User{
		ID: userID,
	}, nil
}

func (a *OIDCAuthenticator) GetID() string {
	return a.id
}

func (a *OIDCAuthenticator) isResponsible(token string) (bool, error) {

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
