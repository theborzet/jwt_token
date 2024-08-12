package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	config "github.com/theborzet/jwt_token/configs"
)

type Claims struct {
	UserID    string `json:"sub"`
	IPAddress string `json:"ip"`
	jwt.StandardClaims
}

type RefreshToken struct {
	Token     string
	ExpiresAt time.Time
}

type TokenManager interface {
	NewJWT(userId, ipAddress string) (string, error)
	ParseJWT(accessToken string) (string, string, error)
	NewRefreshToken() (*RefreshToken, error)
}

type Manager struct {
	SigningKey string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{SigningKey: cfg.Token.SigningKey,
		AccessTTL:  cfg.Token.AccessTokenLifetime,
		RefreshTTL: cfg.Token.RefreshTokenLifetime}
}

func (m *Manager) NewJWT(userId, ipAddress string) (string, error) {
	claims := Claims{
		UserID:    userId,
		IPAddress: ipAddress,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.AccessTTL).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.SigningKey))
}

func (m *Manager) ParseJWT(accessToken string) (string, string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexcepted method: %v", token.Header["alg"])
		}

		return []byte(m.SigningKey), nil
	})
	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(Claims)
	if !ok || !token.Valid {
		return "", "", fmt.Errorf("token receiving error")
	}

	return claims.UserID, claims.IPAddress, nil
}

func (m *Manager) NewRefreshToken() (*RefreshToken, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	refreshToken := RefreshToken{
		Token:     base64.StdEncoding.EncodeToString(b),
		ExpiresAt: time.Now().Add(m.RefreshTTL),
	}

	return &refreshToken, nil
}
