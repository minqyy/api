package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/minqyy/api/internal/config"
	"time"
)

type Manager struct {
	config config.Token
}

type DefaultClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

type Pair struct {
	AccessToken  string
	RefreshToken string
}

func New(config config.Token) *Manager {
	return &Manager{config: config}
}

func (m *Manager) GenerateTokenPair(userID string) (*Pair, error) {
	accessToken, err := m.generateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken := m.generateRefreshToken()

	return &Pair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (m *Manager) generateAccessToken(ID string) (string, error) {
	claims := DefaultClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.config.Access.TTL)),
		},
		ID: ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.config.Access.Secret))
}

func (m *Manager) ParseAccessToken(accessToken string) (*DefaultClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(accessToken, &DefaultClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("token.ParseAccessToken: invalid signing method")
		}
		return []byte(m.config.Access.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*DefaultClaims)
	if !ok {
		return nil, errors.New("token.ParseAccessToken: token claims are not of type *DefaultClaims")
	}

	return claims, nil
}

func (m *Manager) generateRefreshToken() string {
	return uuid.NewString()
}
