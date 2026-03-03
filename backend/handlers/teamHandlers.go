package handlers

import (
	"net/http"
	"strconv"

	"github.com/tomaszSkrzyp/good-game/db"
)

func GetTeamByID(w http.ResponseWriter, r *http.Request, repo *db.TeamRepository) {
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

	team, err := repo.GetByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if team == nil {
		ErrorResponse(w, http.StatusNotFound, "team not found")
		return
	}

	JSONResponse(w, http.StatusOK, team)
}

func FilterTeams(w http.ResponseWriter, r *http.Request, repo *db.TeamRepository) {
	q := r.URL.Query()
	name := q.Get("name")
	conferenceName := q.Get("conferenceName")

	var conferenceID uint
	if conferenceIDStr := q.Get("conferenceId"); conferenceIDStr != "" {
		if id, err := strconv.ParseUint(conferenceIDStr, 10, 64); err == nil {
			conferenceID = uint(id)
		}
	}

	teams, err := repo.Filter(name, conferenceName, conferenceID)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, teams)
}
