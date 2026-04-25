package routes

import (
	"net/http"

	"github.com/tomaszSkrzyp/good-game/handlers"
	"github.com/tomaszSkrzyp/good-game/middleware"
)

func RegisterConfigRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /config", handlers.GetAlgorithmConfig)

	mux.HandleFunc("POST /config",
		middleware.AuthMiddleware(
			middleware.AdminOnly(
				handlers.UpdateConfig,
			),
		),
	)
}
