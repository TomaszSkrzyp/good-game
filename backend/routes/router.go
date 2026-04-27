package routes

import (
	"net/http"

	"github.com/tomaszSkrzyp/good-game/handlers"
	"github.com/tomaszSkrzyp/good-game/middleware"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) http.Handler {
	mainMux := http.NewServeMux()
	apiMux := http.NewServeMux()

	RegisterAuthRoutes(apiMux, db)
	RegisterGameRoutes(apiMux, db)
	RegisterTeamRoutes(apiMux, db)
	RegisterConferenceRoutes(apiMux, db)
	RegisterTeamStatsRoutes(apiMux, db)
	RegisterUserReactionRoutes(apiMux, db)
	RegisterConfigRoutes(apiMux, db)

	apiMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		handlers.JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mainMux.Handle("/api/", http.StripPrefix("/api", apiMux))

	return middleware.EnableCORS(mainMux)
}
