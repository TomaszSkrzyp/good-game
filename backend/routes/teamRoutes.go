package routes

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/handlers"
	"github.com/tomaszSkrzyp/good-game/services"
)

func RegisterTeamRoutes(gormDB *gorm.DB) {
	teamRepo := db.NewTeamRepository(gormDB)
	teamService := services.NewTeamService(teamRepo)

	http.HandleFunc("/teams", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handlers.FilterTeams(w, r, teamService)
	})

	http.HandleFunc("/team", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handlers.GetTeamByID(w, r, teamService)
	})
}
