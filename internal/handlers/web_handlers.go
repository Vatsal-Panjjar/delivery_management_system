package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db     *sqlx.DB
	rdb    *redis.Client
	mu     sync.Mutex
	ctx    context.Context
	active map[int]chan bool
}

func NewHandler(db *sqlx.DB, rdb *redis.Client) *Handler {
	return &Handler{
		db:     db,
		rdb:    rdb,
		ctx:    context.Background(),
		active: make(map[int]chan bool),
	}
}

type User struct {
	ID           int    `db:"id"`
	Username     string `db:"username"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
	Role         string `db:"role"`
}

type Order struct {
	ID              int       `db:"id"`
	UserID          int       `db:"user_id"`
	PickupAddress   string    `db:"pickup_address"`
	DeliveryAddress string    `db:"delivery_address"`
	Status          string    `db:"status"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

// ---------------- AUTH ----------------

func (h *Handler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "📝 Register: POST /register {username, email, password}")
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if username == "" || email == "" || password == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err := h.db.Exec("INSERT INTO users (username, email, password_hash, role) VALUES ($1, $2, $3, 'customer')",
		username, email, string(hash))
	if err != nil {
		http.Error(w, "error registering user", 500)
		return
	}

	fmt.Fprintf(w, "✅ User %s registered successfully!\n", username)
}

func (h *Handler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "🔑 Login: POST /login {email, password}")
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	var user User
	err := h.db.Get(&user, "SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token := fmt.Sprintf("token-%d-%d", user.ID, time.Now().Unix())
	h.rdb.Set(h.ctx, token, user.Role, 24*time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   token,
		Expires: time.Now().Add(24 * time.Hour),
	})

	fmt.Fprintf(w, "✅ Logged in as %s (%s)\n", user.Username, user.Role)
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err == nil {
		h.rdb.Del(h.ctx, c.Value)
		http.SetCookie(w, &http.Cookie{Name: "session_token", Value: "", Expires: time.Now().Add(-1 * time.Hour)})
	}
	fmt.Fprintln(w, "👋 Logged out")
}

// ---------------- MIDDLEWARE ----------------

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		role, err := h.rdb.Get(h.ctx, c.Value).Result()
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		r.Header.Set("X-User-Role", role)
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Header.Get("X-User-Role")
		if role != "admin" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ---------------- USER ----------------

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Queryx("SELECT * FROM orders")
	if err != nil {
		http.Error(w, "error loading dashboard", 500)
		return
	}
	var orders []Order
	for rows.Next() {
		var o Order
		rows.StructScan(&o)
		orders = append(orders, o)
	}
	json.NewEncoder(w).Encode(orders)
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	pickup := r.FormValue("pickup")
	delivery := r.FormValue("delivery")

	if pickup == "" || delivery == "" {
		http.Error(w, "missing fields", 400)
		return
	}

	userID := 1 // temporary (replace with session)
	var id int
	err := h.db.QueryRow("INSERT INTO orders (user_id, pickup_address, delivery_address, status) VALUES ($1,$2,$3,'pending') RETURNING id",
		userID, pickup, delivery).Scan(&id)
	if err != nil {
		http.Error(w, "error creating order", 500)
		return
	}

	fmt.Fprintf(w, "✅ Order %d created\n", id)

	// start tracking
	h.mu.Lock()
	ch := make(chan bool)
	h.active[id] = ch
	h.mu.Unlock()
	go h.trackOrder(id, ch)
}

func (h *Handler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	h.mu.Lock()
	if stop, ok := h.active[id]; ok {
		close(stop)
		delete(h.active, id)
	}
	h.mu.Unlock()

	_, err := h.db.Exec("UPDATE orders SET status='cancelled', updated_at=NOW() WHERE id=$1", id)
	if err != nil {
		http.Error(w, "error cancelling order", 500)
		return
	}
	fmt.Fprintf(w, "🚫 Order %d cancelled\n", id)
}

// ---------------- ADMIN ----------------

func (h *Handler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Queryx("SELECT * FROM orders")
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}
	var orders []Order
	for rows.Next() {
		var o Order
		rows.StructScan(&o)
		orders = append(orders, o)
	}
	json.NewEncoder(w).Encode(orders)
}

func (h *Handler) AdminUpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	status := r.FormValue("status")

	validStatuses := []string{"pending", "in_transit", "delivered"}
	if !contains(validStatuses, status) {
		http.Error(w, "invalid status", 400)
		return
	}

	_, err := h.db.Exec("UPDATE orders SET status=$1, updated_at=NOW() WHERE id=$2", status, id)
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	h.rdb.Set(h.ctx, fmt.Sprintf("order:%d:status", id), status, 0)
	fmt.Fprintf(w, "✅ Order %d status updated to %s\n", id, status)
}

// ---------------- TRACKING ----------------

func (h *Handler) trackOrder(orderID int, stop chan bool) {
	statuses := []string{"pending", "in_transit", "delivered"}
	for _, s := range statuses {
		select {
		case <-stop:
			log.Printf("Order %d tracking stopped", orderID)
			return
		case <-time.After(5 * time.Second):
			h.db.Exec("UPDATE orders SET status=$1, updated_at=NOW() WHERE id=$2", s, orderID)
			h.rdb.Set(h.ctx, fmt.Sprintf("order:%d:status", orderID), s, 0)
			log.Printf("📦 Order %d moved to %s", orderID, s)
			if s == "delivered" {
				h.mu.Lock()
				delete(h.active, orderID)
				h.mu.Unlock()
				return
			}
		}
	}
}

// ---------------- HELPERS ----------------

func contains(arr []string, val string) bool {
	for _, a := range arr {
		if strings.EqualFold(a, val) {
			return true
		}
	}
	return false
}
