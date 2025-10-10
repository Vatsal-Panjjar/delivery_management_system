package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("super_secret_key")

func GenerateJWT(userID string, role string) (string, error) {
	claims := &jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}
	return nil, err
}
