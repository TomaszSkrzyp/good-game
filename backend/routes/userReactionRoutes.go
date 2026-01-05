package routes

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/handlers"
	"github.com/tomaszSkrzyp/good-game/services"
)

func RegisterUserReactionRoutes(gormDB *gorm.DB) {
	repo := db.NewUserReactionRepository(gormDB)
	svc := services.NewUserReactionService(repo)

	http.HandleFunc("/userReactions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.FilterUserReactions(w, r, svc)
		case http.MethodPost:
			handlers.CreateUserReaction(w, r, svc)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/userReaction", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetUserReactionByID(w, r, svc)
		case http.MethodPut:
			handlers.UpdateUserReaction(w, r, svc)
		case http.MethodDelete:
			handlers.DeleteUserReaction(w, r, svc)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// average endpoint
	http.HandleFunc("/userReactions/average/game", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handlers.GetAverageReactionForGame(w, r, svc)
	})

	http.HandleFunc("/userReactions/average/team", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handlers.GetAverageReactionForTeam(w, r, svc)
	})
}
