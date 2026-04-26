package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/engine"
	"github.com/tomaszSkrzyp/good-game/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatalf("error: DATABASE_URL environment variable not set")
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error: couldn't connect to pg: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.Initialize(ctx, gormDB); err != nil {
		log.Fatalf("Critical error during database initialization: %v", err)
	}
	log.Println("db connection ok")

	// run every 15 mins
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	runFetchCycle(gormDB)

	for range ticker.C {
		runFetchCycle(gormDB)
	}
}

func runFetchCycle(gormDB *gorm.DB) error {
	log.Println("starting fetch cycle")

	var unfinishedGames []struct {
		Date time.Time
	}

	// look for games that should have finished but haven't updated yet
	result := gormDB.
		Model(&models.Game{}).
		Where("status != 'STATUS_FINAL' AND game_time < ?", time.Now()).
		Distinct("DATE(game_time)").
		Select("DATE(game_time) as date").
		Scan(&unfinishedGames)

	if result.Error != nil {
		log.Printf("error fetching unfinished game dates: %v", result.Error)
		return result.Error
	}

	if len(unfinishedGames) > 0 {
		minDate := unfinishedGames[0].Date
		maxDate := unfinishedGames[0].Date
		for _, record := range unfinishedGames {
			if record.Date.Before(minDate) {
				minDate = record.Date
			}
			if record.Date.After(maxDate) {
				maxDate = record.Date
			}
		}

		// refresh past games to get final scores
		rangeStr := fmt.Sprintf("%s-%s", minDate.Format("20060102"), maxDate.Format("20060102"))

		log.Printf("updating %v unfinished games for range: %s to %s", len(unfinishedGames), minDate.Format("2006-01-02"), maxDate.Format("2006-01-02"))
		if err := engine.FetchGamesByDate(gormDB, rangeStr); err != nil {
			log.Printf("failed to refresh range %s: %v", rangeStr, err)
			return err
		}
	}

	// grab schedule for the next week
	start := time.Now()
	end := time.Now().AddDate(0, 0, 7)
	rangeStr := fmt.Sprintf("%s-%s", start.Format("20060102"), end.Format("20060102"))
	myrange := "20260420-20260421"
	log.Println("fetching for my playoff range")
	engine.FetchGamesByDate(gormDB, myrange)
	log.Printf("fetching upcoming games: %s to %s", start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err := engine.FetchGamesByDate(gormDB, rangeStr); err != nil {
		log.Printf("failed to fetch upcoming %s: %v", rangeStr, err)
		return err
	}

	log.Println("cycle finished")
	return nil
}
