package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/middleware"
	"github.com/tomaszSkrzyp/good-game/models"
)

func CreateUserReaction(w http.ResponseWriter, r *http.Request, repo *db.UserReactionRepository) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var ur models.UserReaction
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ur.UserID = userID

	if ur.Rating < 1 || ur.Rating > 5 {
		ErrorResponse(w, http.StatusBadRequest, "rating must be between 1 and 5")
		return
	}

	if err := repo.UpdateOrCreate(&ur); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to create reaction")
		return
	}

	//fetch updated stats and return
	avg, count, err := repo.GetStatsForGame(ur.GameID)
	if err != nil {
		JSONResponse(w, http.StatusOK, ur)
		return
	}

	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"rating":      ur.Rating,
		"avgRating":   avg,
		"ratingCount": count,
	})
}

func GetUserReactionByID(w http.ResponseWriter, r *http.Request, repo *db.UserReactionRepository) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	parsed, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid id parameter")
		return
	}

	ur, err := repo.GetByID(uint(parsed))
	if err != nil || ur == nil {
		ErrorResponse(w, http.StatusNotFound, "reaction not found")
		return
	}

	userID, _ := r.Context().Value(middleware.UserIDKey).(uint)
	if userID != 0 && ur.UserID != userID {
		ErrorResponse(w, http.StatusForbidden, "forbidden")
		return
	}

	JSONResponse(w, http.StatusOK, ur)
}

func UpdateUserReaction(w http.ResponseWriter, r *http.Request, repo *db.UserReactionRepository) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	parsed, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid id parameter")
		return
	}

	existing, err := repo.GetByID(uint(parsed))
	if err != nil || existing == nil || existing.UserID != userID {
		ErrorResponse(w, http.StatusForbidden, "forbidden")
		return
	}

	var ur models.UserReaction
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ur.ID = uint(parsed)
	ur.UserID = userID

	if err := repo.Update(&ur); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to update reaction")
		return
	}

	JSONResponse(w, http.StatusOK, ur)
}

func DeleteUserReaction(w http.ResponseWriter, r *http.Request, repo *db.UserReactionRepository) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	parsed, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid id parameter")
		return
	}

	existing, err := repo.GetByID(uint(parsed))
	if err != nil || existing == nil || existing.UserID != userID {
		ErrorResponse(w, http.StatusForbidden, "forbidden")
		return
	}

	if err := repo.Delete(uint(parsed)); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to delete reaction")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func FilterUserReactions(w http.ResponseWriter, r *http.Request, repo *db.UserReactionRepository) {
	q := r.URL.Query()

	ctxID, _ := r.Context().Value(middleware.UserIDKey).(uint)

	var targetUserID uint

	if v := q.Get("userId"); v != "" {
		if p, err := strconv.ParseUint(v, 10, 64); err == nil {
			targetUserID = uint(p)
		}
	} else {
		targetUserID = ctxID
	}
	// security Check: If a user is trying to see a specific ID that isn't theirs, throw 403
	if targetUserID == 0 {
		ErrorResponse(w, http.StatusBadRequest, "no user identified")
		return
	}
	var gameID uint
	if v := q.Get("gameId"); v != "" {
		if p, err := strconv.ParseUint(v, 10, 64); err == nil {
			gameID = uint(p)
		}
	}

	var ratingPtr *int
	if v := q.Get("rating"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 1 && n <= 10 {
			ratingPtr = &n
		}
	}

	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 50
	}

	list, err := repo.Filter(targetUserID, gameID, ratingPtr, page, limit)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to fetch reactions")
		return
	}

	JSONResponse(w, http.StatusOK, list)
}

func GetAverageReactionForGame(w http.ResponseWriter, r *http.Request, repo *db.UserReactionRepository) {
	gameIDStr := r.URL.Query().Get("gameId")
	if gameIDStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing gameId parameter")
		return
	}

	parsed, err := strconv.ParseUint(gameIDStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid gameId parameter")
		return
	}

	avg, count, err := repo.GetStatsForGame(uint(parsed))
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to fetch stats")
		return
	}

	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"gameId":  uint(parsed),
		"average": avg,
		"count":   count,
	})
}

func GetAverageReactionForTeam(w http.ResponseWriter, r *http.Request, repo *db.UserReactionRepository) {
	teamIDStr := r.URL.Query().Get("teamId")
	if teamIDStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing teamId parameter")
		return
	}

	parsed, err := strconv.ParseUint(teamIDStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid teamId parameter")
		return
	}

	avg, count, err := repo.GetStatsForTeam(uint(parsed))
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to fetch stats")
		return
	}

	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"teamId":  uint(parsed),
		"average": avg,
		"count":   count,
	})
}
