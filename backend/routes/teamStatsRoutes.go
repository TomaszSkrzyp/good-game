package routes

import (
	"net/http"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/handlers"
	"github.com/tomaszSkrzyp/good-game/services"
	"gorm.io/gorm"
)

func RegisterTeamStatsRoutes(gormDB *gorm.DB) {
	teamStatsRepo := db.NewTeamStatsRepository(gormDB)
	teamStatsService := services.NewTeamStatsService(teamStatsRepo)

	http.HandleFunc("/teamStats", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateTeamStats(w, r, teamStatsService)
		case http.MethodGet:
			handlers.GetTeamStatsByID(w, r, teamStatsService)
		case http.MethodPut:
			handlers.UpdateTeamStats(w, r, teamStatsService)
		case http.MethodDelete:
			handlers.DeleteTeamStats(w, r, teamStatsService)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/teamStatsFilter", func(w http.ResponseWriter, r *http.Request) {
		handlers.FilterTeamStats(w, r, teamStatsService)
	})
}
