package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims {
		"sub": username,
		"epx": time.Now().Add(20 * time.Minute),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenSing, err :=  token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("erro ao criar token: %w", err)
	}

	return tokenSing, nil
}