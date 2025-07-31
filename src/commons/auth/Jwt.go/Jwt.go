package auth

import (
	"time"

	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(username string) (string, error) {
	return generateJWT(username, 30 * time.Minute)
}

func GenerateRefreshJWT(username string) (string, error) {
	return generateJWT(username, 7 * 24 * time.Hour)
}

func generateJWT(username string, duration time.Duration) (string, error) {
	secret := configuration.Instance().Secret()

	expirationTime := time.Now().Add(duration)
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
		return claims, err
	}
	
	return claims, nil
}