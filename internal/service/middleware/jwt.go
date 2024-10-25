package middleware

import (
	"fmt"
	"github.com/gorilla/sessions"
	"time"

	jwt "github.com/golang-jwt/jwt"
)

type JwtToken struct {
	Secret []byte
}

func NewJwtToken(secret string) (*JwtToken, error) {
	return &JwtToken{Secret: []byte(secret)}, nil
}

type JwtCsrfClaims struct {
	SessionID string `json:"sid"`
	UserID    string `json:"uid"`
	jwt.StandardClaims
}

func (tk *JwtToken) Create(s *sessions.Session, tokenExpTime int64) (string, error) {
	data := JwtCsrfClaims{
		SessionID: s.Values["session_id"].(string),
		UserID:    s.Values["id"].(string),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpTime,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	return token.SignedString(tk.Secret)
}

func (tk *JwtToken) Validate(tokenString string) (*JwtCsrfClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCsrfClaims{}, tk.parseSecretGetter)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*JwtCsrfClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (tk *JwtToken) parseSecretGetter(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok || method.Alg() != "HS256" {
		return nil, fmt.Errorf("bad sign method")
	}
	return tk.Secret, nil
}
