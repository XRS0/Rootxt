package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	Secret []byte
	TTL    time.Duration
}

type Claims struct {
	UserID int64 `json:"uid"`
	jwt.RegisteredClaims
}

func (m JWTManager) Generate(userID int64) (string, error) {
	expiresAt := time.Now().Add(m.TTL)
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.Secret)
}

func (m JWTManager) Verify(tokenString string) (int64, error) {
	parsed, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.Secret, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.UserID, nil
}
