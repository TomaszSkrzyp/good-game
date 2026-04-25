package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/models"
)

func CreateTeamStats(w http.ResponseWriter, r *http.Request, repo *db.TeamStatsRepository) {
	var teamStats models.TeamStats
	if err := json.NewDecoder(r.Body).Decode(&teamStats); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := repo.Create(r.Context(), &teamStats); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to create team stats")
		return
	}

	JSONResponse(w, http.StatusCreated, teamStats)
}

func GetTeamStatsByID(w http.ResponseWriter, r *http.Request, repo *db.TeamStatsRepository) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing id query parameter")
		return
	}

	parsed, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid id parameter")
		return
	}

	teamStats, err := repo.GetByID(r.Context(), uint(parsed))
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "could not retrieve team stats")
		return
	}

	if teamStats == nil {
		ErrorResponse(w, http.StatusNotFound, "team stats not found")
		return
	}

	JSONResponse(w, http.StatusOK, teamStats)
}

func UpdateTeamStats(w http.ResponseWriter, r *http.Request, repo *db.TeamStatsRepository) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing id query parameter")
		return
	}

	parsed, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid id parameter")
		return
	}

	var teamStats models.TeamStats
	if err := json.NewDecoder(r.Body).Decode(&teamStats); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	teamStats.ID = uint(parsed)

	if err := repo.Update(r.Context(), &teamStats); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to update team stats")
		return
	}

	JSONResponse(w, http.StatusOK, teamStats)
}

func DeleteTeamStats(w http.ResponseWriter, r *http.Request, repo *db.TeamStatsRepository) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing id query parameter")
		return
	}

	parsed, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid id parameter")
		return
	}

	if err := repo.Delete(r.Context(), uint(parsed)); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to delete team stats")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func FilterTeamStats(w http.ResponseWriter, r *http.Request, repo *db.TeamStatsRepository) {
	teamIDStr := r.URL.Query().Get("teamId")
	season := r.URL.Query().Get("season")

	if teamIDStr == "" || season == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing teamId or season query parameter")
		return
	}

	parsed, err := strconv.ParseUint(teamIDStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid teamId parameter")
		return
	}

	teamStats, err := repo.Filter(r.Context(), uint(parsed), season)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "could not filter team stats")
		return
	}

	if teamStats == nil {
		ErrorResponse(w, http.StatusNotFound, "team stats not found")
		return
	}

	JSONResponse(w, http.StatusOK, teamStats)
}
