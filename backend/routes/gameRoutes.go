package routes

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/handlers"
	"github.com/tomaszSkrzyp/good-game/services"
)

func RegisterGameRoutes(gormDB *gorm.DB) {
	gameRepo := db.NewGameRepository(gormDB)
	gameService := services.NewGameService(gameRepo)

	http.HandleFunc("/games", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.FilterGame(w, r, gameService)
		case http.MethodPost:
			handlers.CreateGame(w, r, gameService)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetGameByID(w, r, gameService)
		case http.MethodPut:
			handlers.UpdateGame(w, r, gameService)
		case http.MethodDelete:
			handlers.DeleteGame(w, r, gameService)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
