package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tomaszSkrzyp/good-game/models"
	"github.com/tomaszSkrzyp/good-game/services"
)

func CreateGame(w http.ResponseWriter, r *http.Request, gs *services.GameService) {
	w.Header().Set("Content-Type", "application/json")

	var game models.Game
	if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid request"}`))
		return
	}
	if err := gs.Create(&game); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"could not create game"}`))
		return
	}
	w.WriteHeader(http.StatusCreated)
	resp, _ := json.Marshal(game)
	w.Write(resp)
}
func GetGameByID(w http.ResponseWriter, r *http.Request, gs *services.GameService) {
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

	game, err := gs.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"could not retrieve game"}`))
		return
	}
	if game == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"game not found"}`))
		return
	}
	resp, _ := json.Marshal(game)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
func UpdateGame(w http.ResponseWriter, r *http.Request, gs *services.GameService) {
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

	var game models.Game
	if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid request"}`))
		return
	}
	game.ID = id
	if err := gs.Update(&game); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"could not update game"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(game)
	w.Write(resp)
}
func FilterGame(w http.ResponseWriter, r *http.Request, gs *services.GameService) {
	w.Header().Set("Content-Type", "application/json")

	q := r.URL.Query()
	date := q.Get("date")

	var homeIDPtr *uint
	if v := q.Get("homeTeamId"); v != "" {
		if parsed, err := strconv.ParseUint(v, 10, 64); err == nil {
			tmp := uint(parsed)
			homeIDPtr = &tmp
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid homeTeamId"}`))
			return
		}
	}

	var awayIDPtr *uint
	if v := q.Get("awayTeamId"); v != "" {
		if parsed, err := strconv.ParseUint(v, 10, 64); err == nil {
			tmp := uint(parsed)
			awayIDPtr = &tmp
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid awayTeamId"}`))
			return
		}
	}

	var minRatingPtr *int
	if v := q.Get("minRating"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			tmp := parsed
			minRatingPtr = &tmp
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid minRating"}`))
			return
		}
	}

	var maxRatingPtr *int
	if v := q.Get("maxRating"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			tmp := parsed
			maxRatingPtr = &tmp
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid maxRating"}`))
			return
		}
	}

	sort := q.Get("sort")

	page := 1
	limit := 20
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

	games, err := gs.Filter(date, homeIDPtr, awayIDPtr, minRatingPtr, maxRatingPtr, sort, page, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"could not retrieve games"}`))
		return
	}
	resp, _ := json.Marshal(games)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func DeleteGame(w http.ResponseWriter, r *http.Request, gs *services.GameService) {
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

	if err := gs.Delete(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"could not delete game"}`))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
