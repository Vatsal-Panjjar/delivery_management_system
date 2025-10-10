package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/models"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserRepo *repo.UserRepo
	JWTKey   []byte
}

func NewAuthHandler(userRepo *repo.UserRepo, jwtKey []byte) *AuthHandler {
	return &AuthHandler{UserRepo: userRepo, JWTKey: jwtKey}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var u models.User
	json.NewDecoder(r.Body).Decode(&u)
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	u.PasswordHash = string(hash)
	h.UserRepo.Create(&u)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&creds)
	u, _ := h.UserRepo.GetByUsername(creds.Username)
	bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(creds.Password))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": u.ID, "role": u.Role, "exp": time.Now().Add(24 * time.Hour).Unix()})
	tokenString, _ := token.SignedString(h.JWTKey)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
