package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/models"
	"gorm.io/gorm"
)

type BDLTeam struct {
	Abbreviation string `json:"abbreviation"`
}

type BDLGame struct {
	ID               int     `json:"id"`
	Date             string  `json:"date"`
	HomeTeam         BDLTeam `json:"home_team"`
	VisitorTeam      BDLTeam `json:"visitor_team"`
	HomeTeamScore    int     `json:"home_team_score"`
	VisitorTeamScore int     `json:"visitor_team_score"`
}

type BDLResponse struct {
	Data []BDLGame `json:"data"`
	Meta BDLMeta   `json:"meta"`
}

type BDLMeta struct {
	NextCursor int `json:"next_cursor"`
}

func FetchFullSeason(gormDB *gorm.DB, season int) error {
	apiKey := os.Getenv("BDL_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("BDL_API_KEY environment variable is not set")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	cursor := 0

	for {
		log.Printf("Fetching data (cursor: %d) for season %d...", cursor, season)

		url := fmt.Sprintf("https://api.balldontlie.io/v1/games?seasons[]=%d&per_page=100", season)
		if cursor > 0 {
			url = fmt.Sprintf("%s&cursor=%d", url, cursor)
		}

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", apiKey)

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode == 429 {
			log.Println("Rate limit reached. Sleeping for 60 seconds...")
			resp.Body.Close()
			time.Sleep(60 * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return fmt.Errorf("API error: status %d. Body: %s", resp.StatusCode, string(body))
		}

		var apiRes BDLResponse
		err = json.NewDecoder(resp.Body).Decode(&apiRes)
		resp.Body.Close()

		if err != nil {
			return fmt.Errorf("JSON decoding failed: %w", err)
		}

		if len(apiRes.Data) == 0 {
			break
		}

		for _, g := range apiRes.Data {
			saveGame(gormDB, g)
		}

		if apiRes.Meta.NextCursor == 0 {
			break
		}
		cursor = apiRes.Meta.NextCursor

		time.Sleep(2 * time.Second)
	}

	log.Println("Season fetch completed.")
	return nil
}

func saveGame(gormDB *gorm.DB, g BDLGame) {
	homeID := db.TeamAbbrToIDMap[g.HomeTeam.Abbreviation]
	awayID := db.TeamAbbrToIDMap[g.VisitorTeam.Abbreviation]

	if homeID == 0 || awayID == 0 {
		return
	}

	parsedDate, err := time.Parse("2006-01-02", g.Date)
	if err != nil {
		return
	}

	gormDB.Where(models.Game{
		GameTime:   parsedDate,
		HomeTeamID: homeID,
		AwayTeamID: awayID,
	}).Assign(models.Game{
		HomeTeamPoints: uint(g.HomeTeamScore),
		AwayTeamPoints: uint(g.VisitorTeamScore),
	}).FirstOrCreate(&models.Game{})
}
func FetchGamesByDate(gormDB *gorm.DB, date string) error {
	apiKey := os.Getenv("BDL_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("BDL_API_KEY environment variable is not set")
	}

	client := &http.Client{Timeout: 30 * time.Second}

	url := fmt.Sprintf("https://api.balldontlie.io/v1/games?dates[]=%s", date)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: status %d. Body: %s", resp.StatusCode, string(body))
	}

	var apiRes BDLResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiRes); err != nil {
		return err
	}

	for _, g := range apiRes.Data {
		saveGame(gormDB, g)
	}

	log.Printf("Successfully updated %d games for date %s", len(apiRes.Data), date)
	return nil
}
