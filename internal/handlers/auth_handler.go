package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/cache"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/models"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("supersecretkey") // replace with environment variable in prod

// AuthHandler handles user authentication
type AuthHandler struct {
	UserRepo *repo.UserRepo
	Cache    *cache.RedisCache
}

// NewAuthHandler returns a new AuthHandler
func NewAuthHandler(r *repo.UserRepo, c *cache.RedisCache) *AuthHandler {
	return &AuthHandler{UserRepo: r, Cache: c}
}

// SignupHandle
