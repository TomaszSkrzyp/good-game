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

type GamePlayerStats struct {
	HomeTopScorer    string  `json:"homeTopScorer"`
	HomeTopScorerPts float64 `json:"homeTopScorerPts"`
	HomeTopAssister  string  `json:"homeTopAssister"`
	HomeTopAssists   float64 `json:"homeTopAssists"`
	HomeTopRebounder string  `json:"homeTopRebounder"`
	HomeTopRebounds  float64 `json:"homeTopRebounds"`
	AwayTopScorer    string  `json:"awayTopScorer"`
	AwayTopScorerPts float64 `json:"awayTopScorerPts"`
	AwayTopAssister  string  `json:"awayTopAssister"`
	AwayTopAssists   float64 `json:"awayTopAssists"`
	AwayTopRebounder string  `json:"awayTopRebounder"`
	AwayTopRebounds  float64 `json:"awayTopRebounds"`
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

func FetchGamePlayerStats(espnID string, gameDate time.Time) (*GamePlayerStats, error) {
	// Fetch a 2-day range (1 day before and the game date) to account for timezone differences
	startDate := gameDate.AddDate(0, 0, -1)
	dateParam := fmt.Sprintf("%s-%s", startDate.Format("20060102"), gameDate.Format("20060102"))

	log.Printf("Fetching stats for ESPN ID %s on date range %s", espnID, dateParam)

	url := fmt.Sprintf("http://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard?dates=%s", dateParam)
	log.Printf("Request URL: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ESPN API returned status %d", resp.StatusCode)
	}

	var data ESPNResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("json decode error: %w", err)
	}

	log.Printf("Found %d events for date range %s", len(data.Events), dateParam)

	// Find the game by ESPN ID
	var targetEvent ESPNEvent
	for _, event := range data.Events {
		log.Printf("Checking event ID: %s vs %s", event.ID, espnID)
		if event.ID == espnID {
			targetEvent = event
			break
		}
	}

	if targetEvent.ID == "" {
		return nil, fmt.Errorf("game not found for espn id %s in date range %s", espnID, dateParam)
	}

	if len(targetEvent.Competitions) == 0 {
		return nil, fmt.Errorf("no competition data found")
	}

	comp := targetEvent.Competitions[0]
	var home, away ESPNCompetitor

	for _, c := range comp.Competitors {
		if c.HomeAway == "home" {
			home = c
		} else {
			away = c
		}
	}

	stats := &GamePlayerStats{
		HomeTopScorer:    getTopPlayer(home.Leaders, "points"),
		HomeTopScorerPts: getTopPlayerAmount(home.Leaders, "points"),
		HomeTopAssister:  getTopPlayer(home.Leaders, "assists"),
		HomeTopAssists:   getTopPlayerAmount(home.Leaders, "assists"),
		HomeTopRebounder: getTopPlayer(home.Leaders, "rebounds"),
		HomeTopRebounds:  getTopPlayerAmount(home.Leaders, "rebounds"),
		AwayTopScorer:    getTopPlayer(away.Leaders, "points"),
		AwayTopScorerPts: getTopPlayerAmount(away.Leaders, "points"),
		AwayTopAssister:  getTopPlayer(away.Leaders, "assists"),
		AwayTopAssists:   getTopPlayerAmount(away.Leaders, "assists"),
		AwayTopRebounder: getTopPlayer(away.Leaders, "rebounds"),
		AwayTopRebounds:  getTopPlayerAmount(away.Leaders, "rebounds"),
	}

	return stats, nil
}

func getTopPlayer(leaders []ESPNLeaderCategory, statName string) string {
	for _, category := range leaders {
		if category.Name == statName && len(category.Leaders) > 0 {
			return category.Leaders[0].Athlete.DisplayName
		}
	}
	return ""
}

func getTopPlayerAmount(leaders []ESPNLeaderCategory, statName string) float64 {
	for _, category := range leaders {
		if category.Name == statName && len(category.Leaders) > 0 {
			return category.Leaders[0].Value
		}
	}
	return 0
}

func saveESPNGame(gormDB *gorm.DB, event ESPNEvent) {

	if len(event.Competitions) == 0 {
		return
	}
	if event.Status.Type.Name == "STATUS_POSTPONED" {
		fmt.Println("game has been pospotend hence its skipped")
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
		fmt.Printf(`Skipping unfinished game %s with status %s of %s and teams %s vs %s	`, event.ID, event.Status.Type.Name, event.Date, home.Team.Abbreviation, away.Team.Abbreviation)

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
	hQs := make([]int, len(home.LineScores))
	aQs := make([]int, len(away.LineScores))
	for i := 0; i < len(home.LineScores); i++ {
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
