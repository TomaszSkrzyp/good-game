package routes

import (
	"net/http"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/handlers"
	"github.com/tomaszSkrzyp/good-game/services"
	"gorm.io/gorm"
)

func RegisterConferenceRoutes(mux *http.ServeMux, gormDB *gorm.DB) {
	repo := db.NewConferenceRepository(gormDB)
	svc := services.NewConferenceService(repo)

	mux.HandleFunc("/conferences", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handlers.ListConferences(w, r, svc)
	})

	mux.HandleFunc("/conference", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handlers.GetConferenceByID(w, r, svc)
	})
}
