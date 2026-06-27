package auth

import (
	"context"
	"database/sql"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"

	"navisha-english/backend/internal/middleware"
)

type GoogleHandler struct {
	db          *sql.DB
	oauthConfig *oauth2.Config
}

func NewGoogleHandler(db *sql.DB) *GoogleHandler {
	return &GoogleHandler{
		db: db,
		oauthConfig: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

// Redirect redirects the user to Google OAuth consent screen
func (h *GoogleHandler) Redirect(c *gin.Context) {
	state := uuid.New().String()

	// Store state in a short-lived cookie for CSRF protection
	// Use secure=true in production so the cookie is only sent over HTTPS
	secure := os.Getenv("GIN_MODE") == "release"
	c.SetCookie("oauth_state", state, 600, "/", "", secure, true)

	url := h.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Callback handles the OAuth2 callback from Google
func (h *GoogleHandler) Callback(c *gin.Context) {
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3010"
	}

	// Validate state to prevent CSRF
	stateCookie, err := c.Cookie("oauth_state")
	if err != nil || stateCookie != c.Query("state") {
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/login?error=invalid_state")
		return
	}
	// Clear the state cookie using the same secure flag as when it was set
	secure := os.Getenv("GIN_MODE") == "release"
	c.SetCookie("oauth_state", "", -1, "/", "", secure, true)

	code := c.Query("code")
	if code == "" {
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/login?error=missing_code")
		return
	}

	// Exchange code for token
	token, err := h.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/login?error=token_exchange_failed")
		return
	}

	// Validate ID token and extract user info
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/login?error=no_id_token")
		return
	}

	payload, err := idtoken.Validate(context.Background(), idToken, h.oauthConfig.ClientID)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/login?error=invalid_id_token")
		return
	}

	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	googleID, _ := payload.Claims["sub"].(string)
	emailVerified, _ := payload.Claims["email_verified"].(bool)

	if email == "" || !emailVerified {
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/login?error=no_email")
		return
	}

	// Upsert user — create if not exists, update name if changed
	user, err := h.upsertUser(googleID, email, name)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/login?error=db_error")
		return
	}

	// Generate JWT tokens
	accessToken, refreshToken, err := generateTokens(user.ID, user.Email)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/login?error=token_generation_failed")
		return
	}

	if err := saveRefreshToken(h.db, user.ID, refreshToken); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/login?error=session_error")
		return
	}

	// Redirect to frontend with tokens in query params.
	// Use url.QueryEscape to safely encode token values.
	redirectURL := frontendURL + "/auth/callback" +
		"?access_token=" + url.QueryEscape(accessToken) +
		"&refresh_token=" + url.QueryEscape(refreshToken)

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// VerifyIDToken verifies a Google ID token sent directly from the frontend
// This supports the Google One Tap / Sign-In button flow
func (h *GoogleHandler) VerifyIDToken(c *gin.Context) {
	var req struct {
		IDToken string `json:"id_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload, err := idtoken.Validate(context.Background(), req.IDToken, h.oauthConfig.ClientID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Google ID token"})
		return
	}

	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	googleID, _ := payload.Claims["sub"].(string)
	emailVerified, _ := payload.Claims["email_verified"].(bool)

	if email == "" || !emailVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found or not verified in token"})
		return
	}

	user, err := h.upsertUser(googleID, email, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user"})
		return
	}

	accessToken, refreshToken, err := generateTokens(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	if err := saveRefreshToken(h.db, user.ID, refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	})
}

func (h *GoogleHandler) upsertUser(googleID, email, name string) (*User, error) {
	var user User
	err := h.db.QueryRow(`
		INSERT INTO navisha_english_users (google_id, email, name, password_hash)
		VALUES ($1, $2, $3, '')
		ON CONFLICT (email) DO UPDATE SET
			google_id = EXCLUDED.google_id,
			name = EXCLUDED.name,
			updated_at = NOW()
		RETURNING id, name, email, level, created_at`,
		googleID, email, name,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Level, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func saveRefreshToken(db *sql.DB, userID, token string) error {
	_, err := db.Exec(
		`INSERT INTO navisha_english_refresh_tokens (id, user_id, token, expires_at) VALUES ($1, $2, $3, $4)`,
		uuid.New().String(), userID, token, time.Now().Add(30*24*time.Hour),
	)
	return err
}

func generateTokens(userID, email string) (string, string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "changeme-use-env-variable"
	}

	claims := &middleware.Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	refreshToken := uuid.New().String()
	return accessToken, refreshToken, nil
}