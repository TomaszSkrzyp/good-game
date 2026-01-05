package routes

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/handlers"
	"github.com/tomaszSkrzyp/good-game/services"
)

func RegisterConferenceRoutes(gormDB *gorm.DB) {
	confRepo := db.NewConferenceRepository(gormDB)
	confService := services.NewConferenceService(confRepo)

	http.HandleFunc("/conferences", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handlers.ListConferences(w, r, confService)
	})

	http.HandleFunc("/conference", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handlers.GetConferenceByID(w, r, confService)
	})
}
