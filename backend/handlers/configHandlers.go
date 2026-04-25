package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tomaszSkrzyp/good-game/fetch"
)

func GetAlgorithmConfig(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w, http.StatusOK, fetch.CurrentConfig)
}

func UpdateConfig(w http.ResponseWriter, r *http.Request) {
	var newConfig fetch.GameQualityConfig
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid config data")
		return
	}

	// update the global variable in the fetch package
	fetch.CurrentConfig = newConfig

	JSONResponse(w, http.StatusOK, map[string]string{"message": "config updated"})
}
