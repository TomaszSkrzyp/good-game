package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tomaszSkrzyp/good-game/models"
	"github.com/tomaszSkrzyp/good-game/services"
)

func CreateGame(w http.ResponseWriter, r *http.Request, gs *services.GameService) {
	var game models.Game
	if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := gs.Create(&game); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "could not create game")
		return
	}

	JSONResponse(w, http.StatusCreated, game)
}

func GetGameByID(w http.ResponseWriter, r *http.Request, gs *services.GameService) {
	idStr := r.URL.Query().Get("id")
	parsed, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid or missing id")
		return
	}

	game, err := gs.GetByID(uint(parsed))
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "could not retrieve game")
		return
	}
	if game == nil {
		ErrorResponse(w, http.StatusNotFound, "game not found")
		return
	}

	JSONResponse(w, http.StatusOK, game)
}

func UpdateGame(w http.ResponseWriter, r *http.Request, gs *services.GameService) {
	idStr := r.URL.Query().Get("id")
	parsed, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid or missing id")
		return
	}

	var game models.Game
	if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	game.ID = uint(parsed)
	if err := gs.Update(&game); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "could not update game")
		return
	}

	JSONResponse(w, http.StatusOK, game)
}

func FilterGames(w http.ResponseWriter, r *http.Request, gs *services.GameService) {
	q := r.URL.Query()

	val := r.Context().Value(UserIDKey)
	var userID uint
	if val != nil {
		if id, ok := val.(uint); ok {
			userID = id
		}
	}

	parseUint := func(key string) *uint {
		if v := q.Get(key); v != "" {
			if p, err := strconv.ParseUint(v, 10, 64); err == nil {
				res := uint(p)
				return &res
			}
		}
		return nil
	}

	parseInt := func(key string) *int {
		if v := q.Get(key); v != "" {
			if p, err := strconv.Atoi(v); err == nil {
				return &p
			}
		}
		return nil
	}

	date := q.Get("date")
	homeID := parseUint("homeTeamId")
	awayID := parseUint("awayTeamId")
	minRating := parseInt("minRating")
	maxRating := parseInt("maxRating")
	sort := q.Get("sort")

	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	games, err := gs.Filter(date, homeID, awayID, minRating, maxRating, sort, page, limit, userID)

	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "could not retrieve games")
		return
	}

	JSONResponse(w, http.StatusOK, games)
}

func DeleteGame(w http.ResponseWriter, r *http.Request, gs *services.GameService) {
	idStr := r.URL.Query().Get("id")
	parsed, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid or missing id")
		return
	}

	if err := gs.Delete(uint(parsed)); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "could not delete game")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
