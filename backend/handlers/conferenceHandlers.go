package handlers

import (
	"net/http"
	"strconv"

	"github.com/tomaszSkrzyp/good-game/db"
)

func GetConferenceByID(w http.ResponseWriter, r *http.Request, repo *db.ConferenceRepository) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid id parameter")
		return
	}

	conference, err := repo.GetByID(r.Context(), uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if conference == nil {
		ErrorResponse(w, http.StatusNotFound, "conference not found")
		return
	}

	JSONResponse(w, http.StatusOK, conference)
}

func ListConferences(w http.ResponseWriter, r *http.Request, repo *db.ConferenceRepository) {
	conferences, err := repo.List(r.Context())
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, conferences)
}
