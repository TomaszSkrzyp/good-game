package routes

import (
	"net/http"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/handlers"
	"github.com/tomaszSkrzyp/good-game/services"
	"gorm.io/gorm"
)

func RegisterGameRoutes(mux *http.ServeMux, gormDB *gorm.DB) {
	repo := db.NewGameRepository(gormDB)
	svc := services.NewGameService(repo)

	mux.HandleFunc("/games", handlers.OptionalAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.FilterGames(w, r, svc)
		case http.MethodPost:
			handlers.CreateGame(w, r, svc)
		default:
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))

	mux.HandleFunc("/game", handlers.OptionalAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetGameByID(w, r, svc)
		case http.MethodPut:
			handlers.UpdateGame(w, r, svc)
		case http.MethodDelete:
			handlers.DeleteGame(w, r, svc)
		default:
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))
}
