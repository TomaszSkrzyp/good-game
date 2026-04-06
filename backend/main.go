package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/fetch"
	"github.com/tomaszSkrzyp/good-game/middleware"
	"github.com/tomaszSkrzyp/good-game/routes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	updateAll := flag.Bool("update-all", false, "Fetch and update all games for the current season")
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

	db.Initialize(gormDB)

	if *updateAll {
		season := 2026
		log.Printf("Starting full season update for %d", season)
		if err := fetch.FetchFullSeason(gormDB, season); err != nil {
			log.Fatalf("failed to fetch season: %v", err)
		}
		log.Println("Season update complete")
		return
	}

	// Start the web server
	router := routes.NewRouter(gormDB)
	fmt.Printf("Server running on: http://localhost:%s\n", port)

	if err := http.ListenAndServe(":"+port, middleware.Recovery(router)); err != nil {
		log.Fatal(err)
	}
}
