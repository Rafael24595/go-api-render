package auth

import (
	"time"

	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(username string) (string, error) {
	secret := configuration.Instance().Secret()

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(secret)
}

func ValidateJWT(tokenString string) (*Claims, error) {
	secret := configuration.Instance().Secret()

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return secret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}
	
	return claims, nil
}