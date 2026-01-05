package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tomaszSkrzyp/good-game/models"
	"github.com/tomaszSkrzyp/good-game/services"
)

func GetConferenceByID(w http.ResponseWriter, r *http.Request, service *services.ConferenceService) {
	w.Header().Set("Content-Type", "application/json")

	q := r.URL.Query()
	idStr := q.Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"missing id parameter"}`))
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid id parameter"}`))
		return
	}

	conference, err := service.GetByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}
	if conference == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"conference not found"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(conference)
	w.Write(resp)
}

func ListConferences(w http.ResponseWriter, r *http.Request, service *services.ConferenceService) {
	w.Header().Set("Content-Type", "application/json")

	var conferences []models.Conference
	var err error
	if conferences, err = service.List(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(conferences)
	w.Write(resp)
}
