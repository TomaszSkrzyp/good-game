package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/fetch"
	"github.com/tomaszSkrzyp/good-game/middleware"
	"github.com/tomaszSkrzyp/good-game/routes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	bootstrap := flag.Bool("bootstrap", false, "Wipe database and repopulate everything")
	updateAll := flag.Bool("update-all", false, "Fetch and update games without wiping users")
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatalf("error: database_url variable has not been set .env")
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if *bootstrap {
		log.Println("BOOTSTRAP: Wiping all data...")
		if err := db.HardReset(gormDB); err != nil {
			log.Fatalf("Reset failed: %v", err)
		}
	}

	if err := db.Initialize(ctx, gormDB); err != nil {
		log.Fatalf("Critical error during database initialization: %v", err)
	}

	if *bootstrap || *updateAll {
		season := 2026
		log.Printf("Starting data fetch for season %d", season)
		if err := fetch.FetchFullSeason(gormDB, season); err != nil {
			log.Fatalf("failed to fetch season: %v", err)
		}
		log.Println("Data sync complete")

		if flag.NFlag() > 0 {
			return
		}
	}

	// Start the web server
	router := routes.NewRouter(gormDB)
	log.Printf("Server running on: http://localhost:%s\n", port)

	if err := http.ListenAndServe(":"+port, middleware.Recovery(router)); err != nil {
		log.Fatal(err)
	}
}
