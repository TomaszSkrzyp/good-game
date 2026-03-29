package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/fetch"
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
		log.Fatalf("failed to connect to database: %v", err)
	}

	db.Initialize(gormDB)
	log.Println("Connected to database successfully")

	// Run fetch cycle every 15 minutes
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	// Run immediately on startup
	runFetchCycle(gormDB)

	for range ticker.C {
		runFetchCycle(gormDB)
	}
}
func runFetchCycle(gormDB *gorm.DB) error {
	log.Println("Starting fetch cycle...")

	// Get all unfinished games (status != "STATUS_FINAL") from any past date
	var unfinishedGames []struct {
		Date time.Time
	}

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

	// If there are unfinished games, fetch them as a range
	if len(unfinishedGames) > 0 {
		// Find min and max dates
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

		rangeStr := fmt.Sprintf("%s-%s", minDate.Format("20060102"), maxDate.Format("20060102"))
		log.Printf("Fetching games for range: %s to %s (updating unfinished)", minDate.Format("2006-01-02"), maxDate.Format("2006-01-02"))
		if err := fetch.FetchGamesByDate(gormDB, rangeStr); err != nil {
			log.Printf("error fetching games for range %s: %v", rangeStr, err)
			return err
		}
	}

	// Fetch a range including present and future games (e.g., today to next 7 days)
	start := time.Now()
	end := time.Now().AddDate(0, 0, 7)
	rangeStr := fmt.Sprintf("%s-%s", start.Format("20060102"), end.Format("20060102"))

	log.Printf("Fetching games for range: %s to %s (new/upcoming games)", start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err := fetch.FetchGamesByDate(gormDB, rangeStr); err != nil {
		log.Printf("error fetching games for range %s: %v", rangeStr, err)
		return err
	}

	log.Println("Fetch cycle completed")
	return nil
}
