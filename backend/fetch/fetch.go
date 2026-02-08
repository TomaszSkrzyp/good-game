package fetch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/models"
	"gorm.io/gorm"
)

// ESPN API Schema
type ESPNResponse struct {
	Events []ESPNEvent `json:"events"`
}

type ESPNEvent struct {
	ID           string            `json:"id"`
	Date         string            `json:"date"`
	Status       ESPNStatus        `json:"status"`
	Competitions []ESPNCompetition `json:"competitions"`
}

type ESPNStatus struct {
	Type struct {
		Name        string `json:"name"`        // e.g., STATUS_FINAL
		Description string `json:"description"` // e.g., "Final"
	} `json:"type"`
}

type ESPNCompetition struct {
	Competitors []ESPNCompetitor `json:"competitors"`
}
type ESPNLineScore struct {
	Value float64 `json:"value"` // Score in that specific quarter
}

type ESPNCompetitor struct {
	HomeAway   string               `json:"homeAway"`
	Score      string               `json:"score"`
	Team       ESPNTeam             `json:"team"`
	Leaders    []ESPNLeaderCategory `json:"leaders"`
	LineScores []ESPNLineScore      `json:"linescores"`
}

type ESPNTeam struct {
	Abbreviation string `json:"abbreviation"`
}

type ESPNLeaderCategory struct {
	Name    string       `json:"name"` // "points", "rebounds", "assists"
	Leaders []ESPNLeader `json:"leaders"`
}

type ESPNLeader struct {
	Value   float64 `json:"value"`
	Athlete struct {
		DisplayName string `json:"displayName"`
	} `json:"athlete"`
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

		time.Sleep(500 * time.Millisecond) // Respectful throttling
	}
	return nil
}

func FetchGamesByDate(gormDB *gorm.DB, dates string) error {
	// dates can be single "YYYYMMDD" or range "YYYYMMDD-YYYYMMDD"
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

	for _, event := range apiRes.Events {
		saveESPNGame(gormDB, event)
	}
	return nil
}

func saveESPNGame(gormDB *gorm.DB, event ESPNEvent) {
	if len(event.Competitions) == 0 {
		return
	}
	comp := event.Competitions[0]
	var home, away ESPNCompetitor

	for _, c := range comp.Competitors {
		if c.HomeAway == "home" {
			home = c
		} else {
			away = c
		}
	}

	homeID := db.TeamAbbrToIDMap[home.Team.Abbreviation]
	awayID := db.TeamAbbrToIDMap[away.Team.Abbreviation]

	if homeID == 0 || awayID == 0 {
		log.Printf("Skip: Unknown team abbr %s or %s", home.Team.Abbreviation, away.Team.Abbreviation)
		return
	}

	parsedTime, err := parseGameTime(event.Date)
	if err != nil {
		log.Printf("Skip: Date parse error for event %s: %v", event.ID, err)
		return
	}
	if event.Status.Type.Name != "STATUS_FINAL" {
		log.Printf("Skipping unfinished game %s", event.ID)
		err = gormDB.Where(models.Game{
			HomeTeamID: homeID,
			AwayTeamID: awayID,
			GameTime:   parsedTime,
		}).Assign(models.Game{

			ESPNID: event.ID,
			Status: event.Status.Type.Name,
		}).FirstOrCreate(&models.Game{}).Error
		return
	}
	hScore, _ := strconv.Atoi(home.Score)
	aScore, _ := strconv.Atoi(away.Score)
	getQScore := func(scores []ESPNLineScore, index int) int {
		if index < len(scores) {
			return int(scores[index].Value)
		}
		return 0
	}
	hQs := make([]int, 4)
	aQs := make([]int, 4)
	for i := 0; i < 4; i++ {
		hQs[i] = getQScore(home.LineScores, i)
		aQs[i] = getQScore(away.LineScores, i)
	}
	gameQuality := calculateFinalQuality(hScore, aScore, hQs, aQs, home.Leaders, away.Leaders)

	// Transactional Upsert
	err = gormDB.Where(models.Game{
		HomeTeamID: homeID,
		AwayTeamID: awayID,
		GameTime:   parsedTime,
	}).Assign(models.Game{
		HomeTeamPoints: uint(hScore),
		AwayTeamPoints: uint(aScore),
		ESPNID:         event.ID,
		GameQuality:    gameQuality,
		Status:         event.Status.Type.Name,
	}).FirstOrCreate(&models.Game{}).Error

	if err != nil {
		log.Printf("Database error saving game %s: %v", event.ID, err)
	}
}

func parseGameTime(dateStr string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04Z",        // ESPN specific (no seconds)
		"2006-01-02T15:04:05.000Z", // Java/ESPN specific ms
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported date format: %s", dateStr)
}
