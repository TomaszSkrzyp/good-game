package routes

import (
	"net/http"

	"github.com/tomaszSkrzyp/good-game/handlers"
	"github.com/tomaszSkrzyp/good-game/middleware"
)

func RegisterConfigRoutes(mux *http.ServeMux) {

	mux.HandleFunc("/config", middleware.OptionalAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetAlgorithmConfig(w, r)

		case http.MethodPost:
			// restricted: check if user is logged in AND is an admin
			val := r.Context().Value(middleware.UserIDKey)
			userRole := r.Context().Value(middleware.UserRoleKey)

			if val == nil || userRole != "admin" {
				handlers.ErrorResponse(w, http.StatusForbidden, "Only administrators can tune the algorithm")
				return
			}

			handlers.UpdateConfig(w, r)

		default:
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))
}
