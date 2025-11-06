package service

import (
	"crypto/rand"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (us *UserServ) generateToken(id int, email, role string) (token string, err error) {
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, 
	jwt.MapClaims{
		"id": id,
		"email": email,
		"role": role,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	secretKey := os.Getenv("JWT_SECRET_KEY")
	tokenString, err := jwt.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (us *UserServ) generateValidationCode(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}