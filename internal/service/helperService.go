package service

import (
	"crypto/rand"
	"fmt"
	"log"
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

func (us *PaymentServ) totalDay(endDate, startDate string) (int, error) {
	endDateDay := endDate
	startDateDay := startDate
	templateDate := "2006-01-02"

	parseEndDate, err := time.Parse(templateDate, endDateDay)
	if err != nil {
		log.Print(err.Error())
		return 0, err
	}

	parseStartDate, err := time.Parse(templateDate, startDateDay)
	if err != nil {
		log.Print(err.Error())
		return 0, err
	}

	getDayFromEndDate := parseEndDate.Day()
	getDayFromStartDate := parseStartDate.Day()

	totalDay := getDayFromEndDate - getDayFromStartDate

	return totalDay, nil
}