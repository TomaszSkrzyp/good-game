package routes

import (
	"net/http"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/handlers"
	"github.com/tomaszSkrzyp/good-game/middleware"
	"gorm.io/gorm"
)

func RegisterConfigRoutes(mux *http.ServeMux, gormDB *gorm.DB) {
	repo := db.NewConfigRepository(gormDB)

	// GET /config - Publicly accessible
	mux.HandleFunc("GET /config", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAlgorithmConfig(w, r)
	})

	// POST /config - Restricted to Admins
	updateHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateConfig(w, r, repo)
	})

	mux.Handle("POST /config",
		middleware.AuthMiddleware(
			middleware.AdminOnly(
				updateHandler,
			),
		),
	)
}
