package handler_helper

import (
	"easystore/models"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateUUID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Failed to generate UUID: %v", err)
	}
	return id.String()
}

func GenerateEmployeeLoginJwt(e *models.Employee) (string, error) {

	claims := jwt.MapClaims{
		"empID":    e.ID,
		"empName":  e.Name,
		"empEmail": e.Email,
		"exp":      time.Now().Add(time.Hour * 1).Unix(), // Expires in 1 hour
		"iat":      time.Now().Unix(),
	}

	secretKey := os.Getenv("JSON_SECRET_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil

}
