package routes

import (
	"net/http"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/handlers"
	"github.com/tomaszSkrzyp/good-game/services"
	"gorm.io/gorm"
)

func RegisterUserReactionRoutes(mux *http.ServeMux, gormDB *gorm.DB) {
	repo := db.NewUserReactionRepository(gormDB)
	svc := services.NewUserReactionService(repo)

	mux.HandleFunc("/userReactions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// allow guests to see reactions (e.g., filtered by game)
			handlers.OptionalAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				handlers.FilterUserReactions(w, r, svc)
			})(w, r)
		case http.MethodPost:
			// strictly require login to create a reaction
			handlers.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				handlers.CreateUserReaction(w, r, svc)
			})(w, r)
		default:
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	mux.HandleFunc("/userReaction", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.OptionalAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				handlers.GetUserReactionByID(w, r, svc)
			})(w, r)
		case http.MethodPut:
			handlers.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				handlers.UpdateUserReaction(w, r, svc)
			})(w, r)
		case http.MethodDelete:
			handlers.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				handlers.DeleteUserReaction(w, r, svc)
			})(w, r)
		default:
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	// stats endpoints are public by default
	mux.HandleFunc("/userReactions/average/game", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		handlers.GetAverageReactionForGame(w, r, svc)
	})

	mux.HandleFunc("/userReactions/average/team", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		handlers.GetAverageReactionForTeam(w, r, svc)
	})
}
