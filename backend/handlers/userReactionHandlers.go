package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tomaszSkrzyp/good-game/models"
	"github.com/tomaszSkrzyp/good-game/services"
)

func CreateUserReaction(w http.ResponseWriter, r *http.Request, urs *services.UserReactionService) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(UserIDKey).(uint)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var ur models.UserReaction
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ur.UserID = userID

	if ur.Liked < 1 || ur.Liked > 10 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Wywołanie metody UpdateOrCreate (Upsert)
	if err := urs.UpdateOrCreate(&ur); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ur)
}

func GetUserReactionByID(w http.ResponseWriter, r *http.Request, urs *services.UserReactionService) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Query().Get("id")
	parsed, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ur, err := urs.GetByID(uint(parsed))
	if err != nil || ur == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userID, _ := r.Context().Value(UserIDKey).(uint)
	if ur.UserID != userID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(ur)
}

func UpdateUserReaction(w http.ResponseWriter, r *http.Request, urs *services.UserReactionService) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(UserIDKey).(uint)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	idStr := r.URL.Query().Get("id")
	parsed, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	existing, err := urs.GetByID(uint(parsed))
	if err != nil || existing == nil || existing.UserID != userID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var ur models.UserReaction
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ur.ID = uint(parsed)
	ur.UserID = userID

	if err := urs.Update(&ur); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ur)
}

func DeleteUserReaction(w http.ResponseWriter, r *http.Request, urs *services.UserReactionService) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(UserIDKey).(uint)
	idStr := r.URL.Query().Get("id")
	parsed, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	existing, _ := urs.GetByID(uint(parsed))
	if existing == nil || existing.UserID != userID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err := urs.Delete(uint(parsed)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func FilterUserReactions(w http.ResponseWriter, r *http.Request, urs *services.UserReactionService) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query()

	userID, _ := r.Context().Value(UserIDKey).(uint)

	if v := q.Get("userId"); v != "" {
		if p, err := strconv.ParseUint(v, 10, 64); err == nil {
			userID = uint(p)
		}
	}

	var gameID uint
	if v := q.Get("gameId"); v != "" {
		if p, err := strconv.ParseUint(v, 10, 64); err == nil {
			gameID = uint(p)
		}
	}

	var likedPtr *int
	if v := q.Get("liked"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 1 && n <= 10 {
			likedPtr = &n
		}
	}

	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit < 1 {
		limit = 50
	}

	list, err := urs.Filter(userID, gameID, likedPtr, page, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(list)
}

func GetAverageReactionForGame(w http.ResponseWriter, r *http.Request, urs *services.UserReactionService) {
	w.Header().Set("Content-Type", "application/json")
	gameIDStr := r.URL.Query().Get("gameId")
	parsed, err := strconv.ParseUint(gameIDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	avg, count, err := urs.GetAverageAndCountForGame(uint(parsed))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"gameId":  uint(parsed),
		"average": avg,
		"count":   count,
	})
}

func GetAverageReactionForTeam(w http.ResponseWriter, r *http.Request, urs *services.UserReactionService) {
	w.Header().Set("Content-Type", "application/json")
	teamIDStr := r.URL.Query().Get("teamId")
	parsed, err := strconv.ParseUint(teamIDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	avg, count, err := urs.GetAverageAndCountForTeam(uint(parsed))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"teamId":  uint(parsed),
		"average": avg,
		"count":   count,
	})
}
