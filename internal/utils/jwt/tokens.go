package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTUtil interface {
	Generate(typeToken string, accountID string) (string, error)
	Validate(tokenStr string) (string, error)
}

type impl struct {
	Secret []byte
	TTL    map[string]time.Duration
}

func NewJWTUtil(ttl map[string]time.Duration, secret []byte) JWTUtil {
	return &impl{
		TTL:    ttl,
		Secret: secret,
	}
}

func (u *impl) Generate(typeToken string, accountID string) (string, error) {
	if accountID == "" {
		return "", fmt.Errorf("JWTUtil : Generate: sub is empty")
	}

	ttl, ok := u.TTL[typeToken]
	if !ok {
		ttl = u.TTL["access"]
	}

	expired := time.Now().Add(ttl).Unix()
	claims := jwt.MapClaims{
		"sub": accountID,
		"exp": expired,
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	sigStr, err := token.SignedString(u.Secret)

	if err != nil {
		return "", err
	}
	return sigStr, nil
}

// 30 * 24 * time.Hour
func (u *impl) Validate(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return u.Secret, nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("JWTUtils: Validate: %w", err)
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims["sub"].(string), nil
}
