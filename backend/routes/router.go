package routes

import (
	"net/http"

	"github.com/tomaszSkrzyp/good-game/handlers"
	"github.com/tomaszSkrzyp/good-game/middleware"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) http.Handler {
	mux := http.NewServeMux()

	RegisterAuthRoutes(mux, db)
	RegisterGameRoutes(mux, db)
	RegisterTeamRoutes(mux, db)
	RegisterConferenceRoutes(mux, db)
	RegisterTeamStatsRoutes(mux, db)
	RegisterUserReactionRoutes(mux, db)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		handlers.JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	return middleware.EnableCORS(mux)
}
