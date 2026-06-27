package auth

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Level     string `json:"level"`
	CreatedAt string `json:"created_at"`
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM navisha_english_users WHERE email=$1)", req.Email).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	var user User
	err = h.db.QueryRow(
		`INSERT INTO navisha_english_users (name, email, password_hash) VALUES ($1, $2, $3)
		RETURNING id, name, email, level, created_at`,
		req.Name, req.Email, string(hash),
	).Scan(&user.ID, &user.Name, &user.Email, &user.Level, &user.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	accessToken, refreshToken, err := generateTokens(user.ID, user.Email)
	if err != nil {
		// Roll back user creation to avoid orphaned accounts
		h.db.Exec("DELETE FROM navisha_english_users WHERE id=$1", user.ID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	if err := h.saveRefreshToken(user.ID, refreshToken); err != nil {
		// Roll back user creation to avoid orphaned accounts
		h.db.Exec("DELETE FROM navisha_english_users WHERE id=$1", user.ID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	var passwordHash string
	err := h.db.QueryRow(
		`SELECT id, name, email, level, created_at, password_hash FROM navisha_english_users WHERE email=$1`,
		req.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Level, &user.CreatedAt, &passwordHash)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	accessToken, refreshToken, err := generateTokens(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	if err := h.saveRefreshToken(user.ID, refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userID string
	var expiresAt time.Time
	err := h.db.QueryRow(
		`SELECT user_id, expires_at FROM navisha_english_refresh_tokens WHERE token=$1`,
		req.RefreshToken,
	).Scan(&userID, &expiresAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if time.Now().After(expiresAt) {
		h.db.Exec("DELETE FROM navisha_english_refresh_tokens WHERE token=$1", req.RefreshToken)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		return
	}

	var user User
	err = h.db.QueryRow(
		`SELECT id, name, email, level, created_at FROM navisha_english_users WHERE id=$1`, userID,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Level, &user.CreatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	accessToken, newRefreshToken, err := generateTokens(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	if _, err := h.db.Exec("DELETE FROM navisha_english_refresh_tokens WHERE token=$1", req.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invalidate token"})
		return
	}
	if err := h.saveRefreshToken(user.ID, newRefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	// Ignore bind errors — logout should always succeed from client perspective.
	// We still try to read and invalidate the token if provided.
	_ = c.ShouldBindJSON(&req)
	if req.RefreshToken != "" {
		h.db.Exec("DELETE FROM navisha_english_refresh_tokens WHERE token=$1", req.RefreshToken)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *Handler) Me(c *gin.Context) {
	userID := c.GetString("user_id")
	var user User
	err := h.db.QueryRow(
		`SELECT id, name, email, level, created_at FROM navisha_english_users WHERE id=$1`, userID,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Level, &user.CreatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) saveRefreshToken(userID, token string) error {
	return saveRefreshToken(h.db, userID, token)
}
