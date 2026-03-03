package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/fetch"
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
	// Calculate the start date (2 days ago) and end date (today)
	start := time.Now().AddDate(0, 0, -9)
	end := time.Now()

	startDate := start.Format("20060102")
	endDate := end.Format("20060102")
	rangeStr := fmt.Sprintf("%s-%s", startDate, endDate)

	log.Printf("Updating games for the last two days: %s to %s", startDate, endDate)

	err := fetch.FetchGamesByDate(gormDB, rangeStr)
	if err != nil {
		log.Printf("Error fetching games for range %s: %v", rangeStr, err)
		return err
	}

	return nil
}
