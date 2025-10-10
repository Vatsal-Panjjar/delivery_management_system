package auth


import (
"context"
"errors"
"net/http"
"strings"


"github.com/golang-jwt/jwt/v5"
)


var ErrUnauthorized = errors.New("unauthorized")


type Claims struct {
UserID int `json:"user_id"`
Role string `json:"role"`
jwt.RegisteredClaims
}


func GenerateToken(secret string, userID int, role string) (string, error) {
claims := Claims{UserID: userID, Role: role}
t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
return t.SignedString([]byte(secret))
}


// Middleware extracts JWT and puts claims into context
func Middleware(secret string) func(next http.Handler) http.Handler {
return func(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
auth := r.Header.Get("Authorization")
if auth == "" { http.Error(w, "missing auth", http.StatusUnauthorized); return }
parts := strings.SplitN(auth, " ", 2)
if len(parts) != 2 || parts[0] != "Bearer" { http.Error(w, "bad auth", http.StatusUnauthorized); return }
tok, err := jwt.ParseWithClaims(parts[1], &Claims{}, func(t *jwt.Token) (interface{}, error) { return []byte(secret), nil })
if err != nil || !tok.Valid { http.Error(w, "invalid token", http.StatusUnauthorized); return }
claims := tok.Claims.(*Claims)
ctx := context.WithValue(r.Context(), "claims", claims)
next.ServeHTTP(w, r.WithContext(ctx))
})
}
}


func FromContext(r *http.Request) (*Claims, error) {
c := r.Context().Value("claims")
if c == nil { return nil, ErrUnauthorized }
claims, ok := c.(*Claims)
if !ok { return nil, ErrUnauthorized }
return claims, nil
}
