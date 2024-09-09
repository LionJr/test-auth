package token_manager

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"time"
)

type TokenManager struct {
	accessSigningKey  string
	refreshSigningKey string
}

type Time struct {
	time.Time
}

func NewManager(accessSigningKey, refreshSigningKey string) (*TokenManager, error) {
	if accessSigningKey == "" || refreshSigningKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &TokenManager{
		accessSigningKey:  accessSigningKey,
		refreshSigningKey: refreshSigningKey,
	}, nil
}

func (m *TokenManager) NewAccessToken(userId, jti, ip string, ttl time.Duration) (string, error) {
	accessClaims := &AccessClaims{
		StandardClaims: jwt.StandardClaims{
			ID:        jti,
			ExpiresAt: (*jwt.Time)(&Time{time.Now().Add(ttl)}),
		},
		UserID: userId,
		IP:     ip,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims)

	return token.SignedString([]byte(m.accessSigningKey))
}

func (m *TokenManager) ParseAccessToken(accessToken string) (*AccessClaims, error) {
	accessClaims := &AccessClaims{}
	token, err := jwt.ParseWithClaims(accessToken, accessClaims, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.accessSigningKey), nil
	})
	if err != nil {
		return nil, err
	}

	accessClaims, ok := token.Claims.(*AccessClaims)
	if !ok {
		return nil, fmt.Errorf("error get user claims from access token")
	}

	return accessClaims, nil
}

func (m *TokenManager) NewRefreshToken(jti, ip string, ttl time.Duration) (string, error) {
	refreshClaims := &RefreshClaims{
		StandardClaims: jwt.StandardClaims{
			ID:        jti,
			ExpiresAt: (*jwt.Time)(&Time{time.Now().Add(ttl)}),
		},
		IP: ip,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshClaims)

	return token.SignedString([]byte(m.refreshSigningKey))
}

func (m *TokenManager) ParseRefreshToken(refreshToken string) (*RefreshClaims, error) {
	refreshClaims := &RefreshClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, refreshClaims, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.refreshSigningKey), nil
	})
	if err != nil {
		return nil, err
	}

	refreshClaims, ok := token.Claims.(*RefreshClaims)
	if !ok {
		return nil, fmt.Errorf("error get user claims from refresh token")
	}

	return refreshClaims, nil
}
