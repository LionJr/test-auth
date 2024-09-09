package token_manager

import "github.com/dgrijalva/jwt-go/v4"

type AccessClaims struct {
	jwt.StandardClaims

	UserID string `json:"user_id"`
	IP     string `json:"ip"`
}

type RefreshClaims struct {
	jwt.StandardClaims

	IP string `json:"ip"`
}
