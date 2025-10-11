package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("YourSecretKey") // change this to a strong secret

func GenerateToken(username, role string) (string, error) {
    claims := jwt.MapClaims{
        "username": username,
        "role":     role,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func ParseToken(tokenStr string) (map[string]interface{}, error) {
    claims := jwt.MapClaims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil || !token.Valid {
        return nil, err
    }
    return claims, nil
}
