package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Vatsal-Panjiar/delivery_management_system/internal/models"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/repo"
)

type AuthHandler struct {
	UserRepo *repo.UserRepo
	JWTKey   []byte
}

func NewAuthHandler(ur *repo.UserRepo, jwtKey []byte) *AuthHandler {
	return &AuthHandler{
		UserRepo: ur,
		JWTKey:   jwtKey,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	req.ID = uuid.NewString()
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.PasswordHash), bcrypt.DefaultCost)
	req.PasswordHash = string(hash)
	req.CreatedAt = time.Now()

	if err := h.UserRepo.Create(&req); err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	user, err := h.UserRepo.GetByUsername(req.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.PasswordHash)) != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, _ := token.SignedString(h.JWTKey)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
