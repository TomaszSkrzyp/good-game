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
	var ur models.UserReaction
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid request"}`))
		return
	}
	if ur.Liked < 1 || ur.Liked > 10 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"liked must be between 1 and 10"}`))
		return
	}
	if err := urs.Create(&ur); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"could not create user reaction"}`))
		return
	}
	w.WriteHeader(http.StatusCreated)
	resp, _ := json.Marshal(ur)
	w.Write(resp)
}

func GetUserReactionByID(w http.ResponseWriter, r *http.Request, urs *services.UserReactionService) {
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

	ur, err := urs.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"could not retrieve user reaction"}`))
		return
	}
	if ur == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"user reaction not found"}`))
		return
	}
	resp, _ := json.Marshal(ur)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func UpdateUserReaction(w http.ResponseWriter, r *http.Request, urs *services.UserReactionService) {
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

	var ur models.UserReaction
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid request"}`))
		return
	}
	if ur.Liked < 1 || ur.Liked > 10 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"liked must be between 1 and 10"}`))
		return
	}
	ur.ID = id
	if err := urs.Update(&ur); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"could not update user reaction"}`))
		return
	}
	resp, _ := json.Marshal(ur)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func DeleteUserReaction(w http.ResponseWriter, r *http.Request, urs *services.UserReactionService) {
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
	if err := urs.Delete(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"could not delete user reaction"}`))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func FilterUserReactions(w http.ResponseWriter, r *http.Request, urs *services.UserReactionService) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query()

	var userID uint
	if v := q.Get("userId"); v != "" {
		if parsed, err := strconv.ParseUint(v, 10, 64); err == nil {
			userID = uint(parsed)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid userId"}`))
			return
		}
	}

	var gameID uint
	if v := q.Get("gameId"); v != "" {
		if parsed, err := strconv.ParseUint(v, 10, 64); err == nil {
			gameID = uint(parsed)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid gameId"}`))
			return
		}
	}

	var likedPtr *int
	if v := q.Get("liked"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			if n < 1 || n > 10 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"liked must be between 1 and 10"}`))
				return
			}
			likedPtr = &n
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid liked param"}`))
			return
		}
	}

	page := 1
	limit := 50
	if p := q.Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid page param"}`))
			return
		}
	}
	if l := q.Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid limit param"}`))
			return
		}
	}

	list, err := urs.Filter(userID, gameID, likedPtr, page, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"could not retrieve user reactions"}`))
		return
	}
	resp, _ := json.Marshal(list)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func GetAverageReactionForGame(w http.ResponseWriter, r *http.Request, urs *services.UserReactionService) {
	w.Header().Set("Content-Type", "application/json")
	gameIDStr := r.URL.Query().Get("gameId")
	if gameIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"missing gameId query param"}`))
		return
	}
	parsed, err := strconv.ParseUint(gameIDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid gameId"}`))
		return
	}
	avg, err := urs.GetAverageForGame(uint(parsed))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"could not compute average"}`))
		return
	}
	resp, _ := json.Marshal(map[string]interface{}{
		"gameId":  uint(parsed),
		"average": avg,
	})
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func GetAverageReactionForTeam(w http.ResponseWriter, r *http.Request, urs *services.UserReactionService) {
	w.Header().Set("Content-Type", "application/json")
	teamIDStr := r.URL.Query().Get("teamId")
	if teamIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"missing teamId query param"}`))
		return
	}
	parsed, err := strconv.ParseUint(teamIDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid teamId"}`))
		return
	}
	avg, err := urs.GetAverageForTeam(uint(parsed))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"could not compute average"}`))
		return
	}
	resp, _ := json.Marshal(map[string]interface{}{
		"teamId":  uint(parsed),
		"average": avg,
	})
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
