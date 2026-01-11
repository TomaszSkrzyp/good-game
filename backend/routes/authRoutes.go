package routes

import (
	"net/http"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/handlers"
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

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handlers.RegisterUserHandler(w, r, svc)
	})

	mux.HandleFunc("/profile", handlers.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetProfileHandler(w, r, svc)
	}))
}
