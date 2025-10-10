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
	Repo      *repo.UserRepo
	JWTSecret []byte
}

func NewAuthHandler(userRepo *repo.UserRepo, secret []byte) *AuthHandler {
	return &AuthHandler{
		Repo:      userRepo,
		JWTSecret: secret,
	}
}

// Signup creates a new user
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"` // "customer" or "admin"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerEr
