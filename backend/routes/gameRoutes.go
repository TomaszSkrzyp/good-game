package routes

import (
	"net/http"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/handlers"
	"github.com/tomaszSkrzyp/good-game/middleware"
	"gorm.io/gorm"
)

func RegisterGameRoutes(mux *http.ServeMux, gormDB *gorm.DB) {
	repo := db.NewGameRepository(gormDB)

	mux.HandleFunc("/games", middleware.OptionalAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.FilterGames(w, r, repo)
		case http.MethodPost:
			handlers.CreateGame(w, r, repo)
		default:
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))

	mux.HandleFunc("/game", middleware.OptionalAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetGameByID(w, r, repo)
		case http.MethodPut:
			handlers.UpdateGame(w, r, repo)
		case http.MethodDelete:
			handlers.DeleteGame(w, r, repo)
		default:
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))

	mux.HandleFunc("/game/stats", middleware.OptionalAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetGameStats(w, r, repo)
		default:
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))
}
