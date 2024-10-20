package utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(userId primitive.ObjectID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userId,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil

}
