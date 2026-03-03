package routes

import (
	"net/http"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/handlers"
	"gorm.io/gorm"
)

func RegisterTeamStatsRoutes(mux *http.ServeMux, gormDB *gorm.DB) {
	repo := db.NewTeamStatsRepository(gormDB)

	mux.HandleFunc("/teamStats", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateTeamStats(w, r, repo)
		case http.MethodGet:
			handlers.GetTeamStatsByID(w, r, repo)
		case http.MethodPut:
			handlers.UpdateTeamStats(w, r, repo)
		case http.MethodDelete:
			handlers.DeleteTeamStats(w, r, repo)
		default:
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	mux.HandleFunc("/teamStatsFilter", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handlers.FilterTeamStats(w, r, repo)
	})
}
