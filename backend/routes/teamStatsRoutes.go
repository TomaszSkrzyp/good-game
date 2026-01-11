package routes

import (
	"net/http"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/handlers"
	"github.com/tomaszSkrzyp/good-game/services"
	"gorm.io/gorm"
)

func RegisterTeamStatsRoutes(mux *http.ServeMux, gormDB *gorm.DB) {
	repo := db.NewTeamStatsRepository(gormDB)
	svc := services.NewTeamStatsService(repo)

	mux.HandleFunc("/teamStats", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateTeamStats(w, r, svc)
		case http.MethodGet:
			handlers.GetTeamStatsByID(w, r, svc)
		case http.MethodPut:
			handlers.UpdateTeamStats(w, r, svc)
		case http.MethodDelete:
			handlers.DeleteTeamStats(w, r, svc)
		default:
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	mux.HandleFunc("/teamStatsFilter", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handlers.FilterTeamStats(w, r, svc)
	})
}
