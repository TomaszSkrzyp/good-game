package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tomaszSkrzyp/good-game/middleware"
	"github.com/tomaszSkrzyp/good-game/services"
)

type LoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request, us *services.UserService) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.UserName) == "" || len(req.Password) < 6 {
		ErrorResponse(w, http.StatusBadRequest, "Username required and password must be at least 6 characters")
		return
	}

	newUser, err := us.Register(r.Context(), req.UserName, req.Password, req.Email, 2)
	if err != nil {
		ErrorResponse(w, http.StatusConflict, err.Error())
		return
	}

	JSONResponse(w, http.StatusCreated, newUser)
}

func LoginHandler(w http.ResponseWriter, r *http.Request, us *services.UserService) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid login request")
		return
	}

	user, err := us.Authenticate(r.Context(), req.UserName, req.Password)
	if err != nil {
		ErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	roleName := user.Role.Name

	// pass the roleName to GenerateToken
	accessToken, err := middleware.GenerateToken(user.ID, user.UserName, roleName, 15*time.Minute)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to create access token")
		return
	}

	refreshToken, err := middleware.GenerateToken(user.ID, user.UserName, roleName, 7*24*time.Hour)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to create refresh token")
		return
	}
	isProduction := os.Getenv("APP_ENV") == "production"

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   isProduction,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"token":    accessToken,
		"userName": user.UserName,
		"email":    user.Email,
		"role":     roleName,
	})
}
func RefreshHandler(w http.ResponseWriter, r *http.Request, us *services.UserService) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		ErrorResponse(w, http.StatusUnauthorized, "Missing refresh token")
		return
	}

	claims, err := middleware.ValidateToken(cookie.Value)
	if err != nil {
		ErrorResponse(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	userID := uint(claims["id"].(float64))
	user, err := us.GetUserByID(r.Context(), userID)
	if err != nil {
		ErrorResponse(w, http.StatusUnauthorized, "User session no longer valid")
		return
	}

	newToken, _ := middleware.GenerateToken(user.ID, user.UserName, user.Role.Name, 15*time.Minute)

	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"token": newToken,
		"role":  user.Role.Name,
	})
}

func GetProfileHandler(w http.ResponseWriter, r *http.Request, us *services.UserService) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == 0 {
		ErrorResponse(w, http.StatusUnauthorized, "Auth context missing")
		return
	}

	user, err := us.GetUserByID(r.Context(), userID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "User record not found")
		return
	}

	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"id":       user.ID,
		"userName": user.UserName,
		"email":    user.Email,
		"role":     user.Role.Name,
	})
}

func GetUserSettingsHandler(w http.ResponseWriter, r *http.Request, us *services.UserService) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == 0 {
		ErrorResponse(w, http.StatusUnauthorized, "Auth context missing")
		return
	}

	user, err := us.GetUserByID(r.Context(), userID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	JSONResponse(w, http.StatusOK, map[string]bool{"hideScores": user.HideScores})
}

func UpdateUserSettingsHandler(w http.ResponseWriter, r *http.Request, us *services.UserService) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == 0 {
		ErrorResponse(w, http.StatusUnauthorized, "Auth context missing")
		return
	}

	var body struct {
		HideScores bool `json:"hideScores"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := us.SetHideScores(r.Context(), userID, body.HideScores); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to update settings")
		return
	}

	JSONResponse(w, http.StatusOK, map[string]string{"message": "Settings updated"})
}
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	isProduction := os.Getenv("APP_ENV") == "production"

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteLaxMode,
	})

	JSONResponse(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}
