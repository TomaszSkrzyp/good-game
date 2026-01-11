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

func sendError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request, us *services.UserService) {
	w.Header().Set("Content-Type", "application/json")

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if strings.TrimSpace(req.UserName) == "" || len(req.Password) < 6 {
		sendError(w, http.StatusBadRequest, "Username required and password min. 6 chars")
		return
	}
	if !strings.Contains(req.Email, "@") {
		sendError(w, http.StatusBadRequest, "Invalid email format")
		return
	}
	newUser, err := us.Register(req.UserName, req.Password, req.Email, 2)
	if err != nil {

		sendError(w, http.StatusConflict, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		ID        uint   `json:"id"`
		UserName  string `json:"userName"`
		Email     string `json:"email"`
		RoleID    uint   `json:"roleId"`
		CreatedAt string `json:"createdAt"`
	}{
		ID:        newUser.ID,
		UserName:  newUser.UserName,
		Email:     newUser.Email,
		RoleID:    newUser.RoleID,
		CreatedAt: newUser.CreatedAt.Format(time.RFC3339),
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request, us *services.UserService) {
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.UserName == "" || req.Password == "" {
		sendError(w, http.StatusBadRequest, "Missing credentials")
		return
	}

	user, err := us.Authenticate(req.UserName, req.Password)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"id":    user.ID,
		"user":  user.UserName,
		"email": user.Email,
		"exp":   expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Could not generate token")
		return
	}

	json.NewEncoder(w).Encode(struct {
		Token    string `json:"token"`
		UserName string `json:"userName"`
		Email    string `json:"email"`
	}{
		Token:    tokenString,
		UserName: user.UserName,
		Email:    user.Email,
	})
}

func GetProfileHandler(w http.ResponseWriter, r *http.Request, us *services.UserService) {
	w.Header().Set("Content-Type", "application/json")

	val := r.Context().Value(UserIDKey)
	if val == nil {
		sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var userID uint
	switch v := val.(type) {
	case uint:
		userID = v
	case float64:
		userID = uint(v)
	default:
		sendError(w, http.StatusInternalServerError, "Invalid User ID format")
		return
	}

	user, err := us.GetUserByID(userID)
	if err != nil {
		sendError(w, http.StatusNotFound, "User not found")
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       user.ID,
		"userName": user.UserName,
		"email":    user.Email,
	})
}
