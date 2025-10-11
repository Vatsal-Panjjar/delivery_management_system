package auth

import (
    "errors"
    "net/http"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// Secret key for JWT signing (use env variable in production)
var jwtSecret = []byte("supersecretkey123") // change this in production

// Claims represents JWT claims
type Claims struct {
    Username string `json:"username"`
    Role     string `json:"role"` // "user" or "admin"
    jwt.RegisteredClaims
}

// GenerateToken generates a JWT for a given username and role
func GenerateToken(username, role string) (string, error) {
    claims := &Claims{
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// VerifyToken extracts and validates JWT from request header
func VerifyToken(r *http.Request) (*Claims, error) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        return nil, errors.New("authorization header missing")
    }

    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return nil, errors.New("invalid authorization header format")
    }

    tokenStr := parts[1]

    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}

// RequireAdmin middleware: require admin role
func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        claims, err := VerifyToken(r)
        if err != nil {
            http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
            return
        }
        if claims.Role != "admin" {
            http.Error(w, "Forbidden: admin access required", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    }
}
