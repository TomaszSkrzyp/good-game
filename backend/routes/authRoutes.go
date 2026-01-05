package routes

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/handlers"
	"github.com/tomaszSkrzyp/good-game/services"
)

func RegisterAuthRoutes(gormDB *gorm.DB) {
	userRepo := db.NewUserRepository(gormDB)
	userService := services.NewUserService(userRepo)

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(w, r, userService)
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterUserHandler(w, r, userService)
	})
}
