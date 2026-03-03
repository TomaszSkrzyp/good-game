package routes

import (
	"net/http"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/handlers"
	"github.com/tomaszSkrzyp/good-game/middleware"
	"gorm.io/gorm"
)

func RegisterUserReactionRoutes(mux *http.ServeMux, gormDB *gorm.DB) {
	repo := db.NewUserReactionRepository(gormDB)

	mux.HandleFunc("/userReactions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			middleware.OptionalAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				handlers.FilterUserReactions(w, r, repo)
			})(w, r)
		case http.MethodPost:
			middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				handlers.CreateUserReaction(w, r, repo)
			})(w, r)
		default:
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	mux.HandleFunc("/userReaction", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			middleware.OptionalAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				handlers.GetUserReactionByID(w, r, repo)
			})(w, r)
		case http.MethodPut:
			middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				handlers.UpdateUserReaction(w, r, repo)
			})(w, r)
		case http.MethodDelete:
			middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				handlers.DeleteUserReaction(w, r, repo)
			})(w, r)
		default:
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	mux.HandleFunc("/userReactions/average/game", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		handlers.GetAverageReactionForGame(w, r, repo)
	})

	mux.HandleFunc("/userReactions/average/team", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		handlers.GetAverageReactionForTeam(w, r, repo)
	})
}
