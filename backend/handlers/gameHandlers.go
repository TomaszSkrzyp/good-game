package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/fetch"
	"github.com/tomaszSkrzyp/good-game/middleware"
	"github.com/tomaszSkrzyp/good-game/models"
)

func CreateGame(w http.ResponseWriter, r *http.Request, repo *db.GameRepository) {
	var game models.Game
	if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := repo.Create(r.Context(), &game); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "could not create game")
		return
	}

	JSONResponse(w, http.StatusCreated, game)
}

func GetGameByID(w http.ResponseWriter, r *http.Request, repo *db.GameRepository) {
	idStr := r.URL.Query().Get("id")
	parsed, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid or missing id")
		return
	}
	var userID uint
	if val := r.Context().Value(middleware.UserIDKey); val != nil {
		if id, ok := val.(uint); ok {
			userID = id
		}
	}

	game, err := repo.GetByID(r.Context(), uint(parsed), userID)
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

func UpdateGame(w http.ResponseWriter, r *http.Request, repo *db.GameRepository) {
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
	if err := repo.Update(r.Context(), &game); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "could not update game")
		return
	}

	JSONResponse(w, http.StatusOK, game)
}

func FilterGames(w http.ResponseWriter, r *http.Request, repo *db.GameRepository) {
	q := r.URL.Query()

	val := r.Context().Value(middleware.UserIDKey)
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

	games, err := repo.Filter(r.Context(), date, homeID, awayID, minRating, maxRating, sort, page, limit, userID)

	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "could not retrieve games")
		return
	}

	JSONResponse(w, http.StatusOK, games)
}

func DeleteGame(w http.ResponseWriter, r *http.Request, repo *db.GameRepository) {
	idStr := r.URL.Query().Get("id")
	parsed, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid or missing id")
		return
	}

	if err := repo.Delete(r.Context(), uint(parsed)); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "could not delete game")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetGameStats(w http.ResponseWriter, r *http.Request, repo *db.GameRepository) {
	gameID := r.URL.Query().Get("gameId")
	if gameID == "" {
		ErrorResponse(w, http.StatusBadRequest, "gameId query parameter required")
		return
	}

	gameIDParsed, err := strconv.ParseUint(gameID, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid gameId")
		return
	}

	var userID uint
	if val := r.Context().Value(middleware.UserIDKey); val != nil {
		if id, ok := val.(uint); ok {
			userID = id
		}
	}

	game, err := repo.GetByID(r.Context(), uint(gameIDParsed), userID)
	if err != nil || game == nil {
		ErrorResponse(w, http.StatusNotFound, "Game not found")
		return
	}

	if game.ESPNID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Game has no ESPN ID")
		return
	}

	stats, err := fetch.FetchGamePlayerStats(game.ESPNID, game.GameTime)
	if err != nil {
		log.Printf("Error fetching stats for ESPN ID %s: %v", game.ESPNID, err)
		ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch game stats")
		return
	}

	JSONResponse(w, http.StatusOK, stats)
}
