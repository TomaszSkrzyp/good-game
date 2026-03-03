package routes

import (
	"net/http"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/handlers"
	"gorm.io/gorm"
)

func RegisterConferenceRoutes(mux *http.ServeMux, gormDB *gorm.DB) {
	repo := db.NewConferenceRepository(gormDB)

	mux.HandleFunc("/conferences", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handlers.ListConferences(w, r, repo)
	})

	mux.HandleFunc("/conference", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handlers.GetConferenceByID(w, r, repo)
	})
}
