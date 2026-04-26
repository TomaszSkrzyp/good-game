package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/models"
)

func GetAlgorithmConfig(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w, http.StatusOK, models.GetConfig())
}

func UpdateConfig(w http.ResponseWriter, r *http.Request, repo *db.ConfigRepository) {
	var newConfig models.GameQualityConfig
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid config data")
		return
	}

	if err := repo.SaveGameQuality(r.Context(), newConfig); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to save to database")
		return
	}

	models.SetConfig(newConfig)

	JSONResponse(w, http.StatusOK, map[string]string{"message": "config updated successfully"})
}
