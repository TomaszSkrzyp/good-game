package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tomaszSkrzyp/good-game/models"
	"github.com/tomaszSkrzyp/good-game/services"
)

func CreateTeamStats(w http.ResponseWriter, r *http.Request, tss *services.TeamStatsService) {
	w.Header().Set("Content-Type", "application/json")

	var teamStats models.TeamStats
	if err := json.NewDecoder(r.Body).Decode(&teamStats); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid request payload"}`))
		return
	}

	if err := tss.Create(&teamStats); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"failed to create team stats"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp, _ := json.Marshal(teamStats)
	w.Write(resp)
}
func GetTeamStatsByID(w http.ResponseWriter, r *http.Request, tss *services.TeamStatsService) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"missing id query param"}`))
		return
	}
	parsed, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid id"}`))
		return
	}
	id := uint(parsed)

	teamStats, err := tss.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"could not retrieve team stats"}`))
		return
	}
	if teamStats == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"team stats not found"}`))
		return
	}
	resp, _ := json.Marshal(teamStats)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
func UpdateTeamStats(w http.ResponseWriter, r *http.Request, tss *services.TeamStatsService) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"missing id query param"}`))
		return
	}
	parsed, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid id"}`))
		return
	}
	id := uint(parsed)

	var teamStats models.TeamStats
	if err := json.NewDecoder(r.Body).Decode(&teamStats); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid request payload"}`))
		return
	}
	teamStats.ID = id

	if err := tss.Update(&teamStats); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"failed to update team stats"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(teamStats)
	w.Write(resp)
}
func DeleteTeamStats(w http.ResponseWriter, r *http.Request, tss *services.TeamStatsService) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"missing id query param"}`))
		return
	}
	parsed, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid id"}`))
		return
	}
	id := uint(parsed)

	if err := tss.Delete(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"failed to delete team stats"}`))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func FilterTeamStats(w http.ResponseWriter, r *http.Request, tss *services.TeamStatsService) {
	w.Header().Set("Content-Type", "application/json")

	teamIDStr := r.URL.Query().Get("teamId")
	season := r.URL.Query().Get("season")
	if teamIDStr == "" || season == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"missing teamId or season query param"}`))
		return
	}
	parsed, err := strconv.ParseUint(teamIDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid teamId"}`))
		return
	}
	teamID := uint(parsed)

	teamStats, err := tss.Filter(teamID, season)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"could not filter team stats"}`))
		return
	}
	if teamStats == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"team stats not found"}`))
		return
	}
	resp, _ := json.Marshal(teamStats)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
