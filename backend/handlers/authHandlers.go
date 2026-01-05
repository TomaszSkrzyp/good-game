package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tomaszSkrzyp/good-game/services"
)

type LoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}
type RegisterRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request, us *services.UserService) {
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid request"}`))
		return
	}
	newUser, err := us.Register(req.UserName, req.Password, 2) //2 - userrole
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
		return
	}
	resp, err := json.Marshal(struct {
		ID       uint   `json:"id"`
		UserName string `json:"userName"`
		RoleID   uint   `json:"roleId"`
		Created  string `json:"createdAt"`
	}{
		ID:       newUser.ID,
		UserName: newUser.UserName,
		RoleID:   newUser.RoleID,
		Created:  newUser.CreatedAt.Format(time.RFC3339),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)

}

func LoginHandler(w http.ResponseWriter, r *http.Request, us *services.UserService) {

	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid request"}`))
		return
	}
	user, err := us.Authenticate(req.UserName, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"invalid credentials"}`))
		return
	}
	resp, err := json.Marshal(struct {
		ID        uint   `json:"id"`
		UserName  string `json:"userName"`
		LastLogin string `json:"lastLogin"`
	}{
		ID:        user.ID,
		UserName:  user.UserName,
		LastLogin: user.LastLoginAt.Format(time.RFC3339),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
