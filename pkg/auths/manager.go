package auths

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/mikalai2006/go-template-api/internal/domain"
)

type TokenManager interface {
	NewJWT(claims domain.DataForClaims, ttl time.Duration) (string, error)
	Parse(accessToken string) (*AuthClaims, error)
	NewRefreshToken() (string, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}

type AuthClaims struct {
	Roles []string `json:"roles"`
	jwt.StandardClaims
}

func (m *Manager) NewJWT(claims domain.DataForClaims, ttl time.Duration) (string, error) {
	claimsData := AuthClaims{
		Roles: claims.Roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			Subject:   claims.UserID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsData)

	signedToken, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (m *Manager) Parse(accessToken string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&AuthClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexcepted signing method: %v", token.Header["alg"])
			}
			return []byte(m.signingKey), nil
		},
	)

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (m *Manager) NewRefreshToken() (string, error) {
	r := uuid.New()

	return r.String(), nil
}
