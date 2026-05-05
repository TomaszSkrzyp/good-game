package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tomaszSkrzyp/good-game/models"
	"gorm.io/gorm"
)

var espnClient = &http.Client{
	Timeout: 10 * time.Second,
}

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

	resp, err := espnClient.Get(url)
	if err != nil {
		return fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	var apiRes ESPNResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiRes); err != nil {
		return fmt.Errorf("json decode error: %w", err)
	}
	log.Printf("Fetched %d events for date(s) %s", len(apiRes.Events), dates)
	activeIDs := make([]string, 0, len(apiRes.Events))
	for _, event := range apiRes.Events {
		activeIDs = append(activeIDs, event.ID)
		log.Print(apiRes.Events)
		saveESPNGame(gormDB, event)
	}

	if err := cleanUpGhostGames(gormDB, activeIDs, dates); err != nil {
		log.Printf("Warning: ghost game cleanup failed: %v", err)
	}
	return nil
}
func cleanUpGhostGames(db *gorm.DB, activeIDs []string, dates string) error {
	var start, end time.Time

	if len(dates) == 17 {
		start, _ = time.Parse("20060102", dates[:8])
		end, _ = time.Parse("20060102", dates[9:])
	} else if len(dates) == 8 {
		start, _ = time.Parse("20060102", dates)
		end = start
	} else {
		return nil
	}

	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	now := time.Now()
	result := db.Where("game_time BETWEEN ? AND ?", start, end).
		Where("id NOT IN ?", activeIDs).
		Where("status != ?", "STATUS_FINAL").
		Where("game_time < ?", now.Add(1*time.Hour)).
		Delete(&models.Game{})
	if result.RowsAffected > 0 {
		log.Printf("Canceled %d ghost games (series ended) for range %s", result.RowsAffected, dates)
	}

	return result.Error
}
