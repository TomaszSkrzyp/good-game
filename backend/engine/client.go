package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// FetchFullSeason requests data in 7-day ranges to minimize API calls
func FetchFullSeason(gormDB *gorm.DB, year int) error {
	log.Printf("Initiating full season sync for year %d", year)

	// NBA Season roughly Oct to June
	start := time.Date(year-1, 10, 20, 0, 0, 0, 0, time.UTC)
	end := time.Date(year, 6, 20, 0, 0, 0, 0, time.UTC)

	for d := start; d.Before(end); d = d.AddDate(0, 0, 7) {
		rangeEnd := d.AddDate(0, 0, 7)
		if rangeEnd.After(end) {
			rangeEnd = end
		}

		dateParam := fmt.Sprintf("%s-%s", d.Format("20060102"), rangeEnd.Format("20060102"))
		log.Printf("Fetching range: %s", dateParam)

		if err := FetchGamesByDate(gormDB, dateParam); err != nil {
			log.Printf("Error fetching range %s: %v", dateParam, err)
		}

		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

func FetchGamesByDate(gormDB *gorm.DB, dates string) error {
	url := fmt.Sprintf("http://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard?dates=%s", dates)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	var apiRes ESPNResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiRes); err != nil {
		return fmt.Errorf("json decode error: %w", err)
	}
	log.Printf("Fetched %d events for date(s) %s", len(apiRes.Events), dates)
	for _, event := range apiRes.Events {
		saveESPNGame(gormDB, event)
	}
	return nil
}
