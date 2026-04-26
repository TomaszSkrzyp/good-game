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
			val := r.Context().Value(middleware.UserIDKey)
			if val == nil {
				handlers.ErrorResponse(w, http.StatusUnauthorized, "You must be logged in to create games")
				return
			}
			handlers.CreateGame(w, r, repo)
		default:
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))

	mux.HandleFunc("/game", middleware.OptionalAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetGameByID(w, r, repo)
		case http.MethodPut, http.MethodDelete:
			val := r.Context().Value(middleware.UserIDKey)
			if val == nil {
				handlers.ErrorResponse(w, http.StatusUnauthorized, "Administrative action requires login")
				return
			}

			if r.Method == http.MethodPut {
				handlers.UpdateGame(w, r, repo)
			} else {
				handlers.DeleteGame(w, r, repo)
			}
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

	mux.HandleFunc("/teams/quality", middleware.OptionalAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetTeamQualityStats(w, r, repo)
		default:
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))
}
