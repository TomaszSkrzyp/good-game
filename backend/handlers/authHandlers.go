package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/tomaszSkrzyp/good-game/services"
)

// Global key for JWT signing loaded from env
var jwtKey = []byte(os.Getenv("JWT_KEY"))

type LoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// Helper for generating JWT tokens with specific expiration time
func generateToken(userID uint, username string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"id":   userID,
		"user": username,
		"exp":  time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func sendError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request, us *services.UserService) {
	w.Header().Set("Content-Type", "application/json")

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Basic validation before reaching service layer
	if strings.TrimSpace(req.UserName) == "" || len(req.Password) < 6 {
		sendError(w, http.StatusBadRequest, "Username required and password must be at least 6 characters")
		return
	}

	newUser, err := us.Register(req.UserName, req.Password, req.Email, 2)
	if err != nil {
		sendError(w, http.StatusConflict, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func LoginHandler(w http.ResponseWriter, r *http.Request, us *services.UserService) {
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid login request")
		return
	}

	user, err := us.Authenticate(req.UserName, req.Password)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Access Token: short-lived (15 min) for API calls
	accessToken, err := generateToken(user.ID, user.UserName, 15*time.Minute)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create access token")
		return
	}

	// Refresh Token: long-lived (7 days) for session persistence
	refreshToken, err := generateToken(user.ID, user.UserName, 7*24*time.Hour)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create refresh token")
		return
	}

	// Set refresh token in HttpOnly cookie to prevent XSS attacks
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   false, // Set to true if using HTTPS
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	// Send non-sensitive data and access token to client
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":    accessToken,
		"userName": user.UserName,
		"email":    user.Email,
	})
}

// Handler for issuing new access tokens using valid refresh cookie
func RefreshHandler(w http.ResponseWriter, r *http.Request, us *services.UserService) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Missing refresh token")
		return
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		sendError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	userID := uint(claims["id"].(float64))
	user, err := us.GetUserByID(userID)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "User session no longer valid")
		return
	}

	// Generate fresh access token
	newToken, _ := generateToken(user.ID, user.UserName, 15*time.Minute)

	json.NewEncoder(w).Encode(map[string]string{
		"token": newToken,
	})
}

func GetProfileHandler(w http.ResponseWriter, r *http.Request, us *services.UserService) {
	w.Header().Set("Content-Type", "application/json")

	// Context value injected by AuthMiddleware
	val := r.Context().Value(UserIDKey)
	if val == nil {
		sendError(w, http.StatusUnauthorized, "Auth context missing")
		return
	}

	var userID uint
	switch v := val.(type) {
	case uint:
		userID = v
	case float64:
		userID = uint(v)
	default:
		sendError(w, http.StatusInternalServerError, "Internal context type error")
		return
	}

	user, err := us.GetUserByID(userID)
	if err != nil {
		sendError(w, http.StatusNotFound, "User record not found")
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       user.ID,
		"userName": user.UserName,
		"email":    user.Email,
	})
}
