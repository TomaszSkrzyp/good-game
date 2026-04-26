package engine

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/models"
	"gorm.io/gorm"
)

func getTopPlayer(leaders []ESPNLeaderCategory, statName string) string {
	for _, category := range leaders {
		if category.Name == statName && len(category.Leaders) > 0 {
			if category.Leaders[0].Athlete.DisplayName == "" {
				continue
			}
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
		log.Printf("Game %s has been postponed; skipping", event.ID)
		return
	}

	comp := event.Competitions[0]
	var home, away ESPNCompetitor

	for _, c := range comp.Competitors {
		if c.HomeAway == "home" {
			home = c
		} else if c.HomeAway == "away" {
			away = c
		}
	}

	getAbbrID := func(abbr string) uint {
		if strings.Contains(abbr, "/") {
			return db.TeamAbbrToIDMap["TBD"]
		}
		return db.TeamAbbrToIDMap[abbr]
	}

	homeID := getAbbrID(home.Team.Abbreviation)
	awayID := getAbbrID(away.Team.Abbreviation)

	if homeID == 0 || awayID == 0 {
		log.Printf("Skip: Unknown team abbr %s or %s", home.Team.Abbreviation, away.Team.Abbreviation)
		return
	}

	parsedTime, err := parseGameTime(event.Date)
	if err != nil {
		log.Printf("Skip: Date parse error for event %s: %v", event.ID, err)
		return
	}

	gameStats := GamePlayerStatsDTO{
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

	if event.Status.Type.Name != "STATUS_FINAL" {
		log.Printf("Syncing unfinished game %s (%s vs %s)", event.ID, home.Team.Abbreviation, away.Team.Abbreviation)

		err = gormDB.Where(models.Game{ESPNID: event.ID}).Assign(models.Game{
			HomeTeamID:       homeID,
			AwayTeamID:       awayID,
			GameTime:         parsedTime,
			Status:           event.Status.Type.Name,
			HomeTopScorer:    gameStats.HomeTopScorer,
			HomeTopScorerPts: gameStats.HomeTopScorerPts,
			HomeTopAssister:  gameStats.HomeTopAssister,
			HomeTopAssists:   gameStats.HomeTopAssists,
			HomeTopRebounder: gameStats.HomeTopRebounder,
			HomeTopRebounds:  gameStats.HomeTopRebounds,
			AwayTopScorer:    gameStats.AwayTopScorer,
			AwayTopScorerPts: gameStats.AwayTopScorerPts,
			AwayTopAssister:  gameStats.AwayTopAssister,
			AwayTopAssists:   gameStats.AwayTopAssists,
			AwayTopRebounder: gameStats.AwayTopRebounder,
			AwayTopRebounds:  gameStats.AwayTopRebounds,
		}).FirstOrCreate(&models.Game{}).Error

		if err != nil {
			log.Printf("Database error saving unfinished game %s: %v", event.ID, err)
		}
		return
	}

	hScore, _ := strconv.Atoi(home.Score)
	aScore, _ := strconv.Atoi(away.Score)

	getLineScores := func(scores []ESPNLineScore) []int {
		res := make([]int, len(scores))
		for i, s := range scores {
			res[i] = int(s.Value)
		}
		return res
	}
	hQs := getLineScores(home.LineScores)
	aQs := getLineScores(away.LineScores)
	dramaCtx, _ := FetchAndCalculateDrama(event.ID)
	gameQuality := CalculateFinalQuality(hScore, aScore, hQs, aQs, home.Leaders, away.Leaders, dramaCtx)
	log.Printf("drama score: %v, quality score: %v for id: %s", dramaCtx.DramaScore, gameQuality, event.ID)
	err = gormDB.Where(models.Game{
		ESPNID: event.ID,
	}).Assign(models.Game{
		HomeTeamID:       homeID,
		AwayTeamID:       awayID,
		GameTime:         parsedTime,
		HomeTeamPoints:   uint(hScore),
		AwayTeamPoints:   uint(aScore),
		GameQuality:      gameQuality,
		Status:           event.Status.Type.Name,
		HomeTopScorer:    gameStats.HomeTopScorer,
		HomeTopScorerPts: gameStats.HomeTopScorerPts,
		HomeTopAssister:  gameStats.HomeTopAssister,
		HomeTopAssists:   gameStats.HomeTopAssists,
		HomeTopRebounder: gameStats.HomeTopRebounder,
		HomeTopRebounds:  gameStats.HomeTopRebounds,
		AwayTopScorer:    gameStats.AwayTopScorer,
		AwayTopScorerPts: gameStats.AwayTopScorerPts,
		AwayTopAssister:  gameStats.AwayTopAssister,
		AwayTopAssists:   gameStats.AwayTopAssists,
		AwayTopRebounder: gameStats.AwayTopRebounder,
		AwayTopRebounds:  gameStats.AwayTopRebounds,
	}).FirstOrCreate(&models.Game{}).Error

	if err != nil {
		log.Printf("Database error saving finalized game %s: %v", event.ID, err)
	}
}

func parseGameTime(dateStr string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04Z",
		"2006-01-02T15:04:05.000Z",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported date format: %s", dateStr)
}
