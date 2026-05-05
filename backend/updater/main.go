package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/tomaszSkrzyp/good-game/db"
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
