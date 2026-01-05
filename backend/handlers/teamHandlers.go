package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tomaszSkrzyp/good-game/models"
	"github.com/tomaszSkrzyp/good-game/services"
)

func GetTeamByID(w http.ResponseWriter, r *http.Request, service *services.TeamService) {
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

	team, err := service.GetByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}
	if team == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"team not found"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(team)
	w.Write(resp)
}

func FilterTeams(w http.ResponseWriter, r *http.Request, service *services.TeamService) {
	w.Header().Set("Content-Type", "application/json")

	q := r.URL.Query()
	name := q.Get("name")
	conferenceIDStr := q.Get("conferenceId")
	conferenceName := q.Get("conferenceName")

	var conferenceID uint64
	var err error
	if conferenceIDStr != "" {
		conferenceID, err = strconv.ParseUint(conferenceIDStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid conferenceId parameter"}`))
			return
		}
	}

	var teams []models.Team
	if teams, err = service.Filter(name, conferenceName, uint(conferenceID)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(teams)
	w.Write(resp)
}
