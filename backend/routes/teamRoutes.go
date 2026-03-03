package routes

import (
	"net/http"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/handlers"
	"gorm.io/gorm"
)

func RegisterTeamRoutes(mux *http.ServeMux, gormDB *gorm.DB) {
	repo := db.NewTeamRepository(gormDB)

	mux.HandleFunc("/teams", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handlers.FilterTeams(w, r, repo)
	})

	mux.HandleFunc("/team", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handlers.GetTeamByID(w, r, repo)
	})
}
