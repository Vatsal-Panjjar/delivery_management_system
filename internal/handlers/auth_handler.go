package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/auth"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// TODO: validate user from DB instead of hardcoding
	if req.Username == "admin" && req.Password == "admin123" {
		token, _ := auth.GenerateJWT("1", "admin")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}
