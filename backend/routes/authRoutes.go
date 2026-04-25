package routes

import (
	"net/http"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/handlers"
	"github.com/tomaszSkrzyp/good-game/middleware"
	"github.com/tomaszSkrzyp/good-game/services"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(mux *http.ServeMux, gormDB *gorm.DB) {
	repo := db.NewUserRepository(gormDB)
	svc := services.NewUserService(repo)

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handlers.LoginHandler(w, r, svc)
	})
	mux.HandleFunc("POST /logout", func(w http.ResponseWriter, r *http.Request) {
		handlers.LogoutHandler(w, r)
	})
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handlers.RegisterUserHandler(w, r, svc)
	})

	mux.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handlers.RefreshHandler(w, r, svc)
	})

	mux.HandleFunc("/profile", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handlers.GetProfileHandler(w, r, svc)
	}))

	mux.HandleFunc("/userSettings", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetUserSettingsHandler(w, r, svc)
		case http.MethodPost:
			handlers.UpdateUserSettingsHandler(w, r, svc)
		default:
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))
}
