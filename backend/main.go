package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/fetch"
	"github.com/tomaszSkrzyp/good-game/routes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
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

	// check if we have enough arguments for cli commands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--fetch-season":
			log.Println("starting full season fetch...")
			if len(os.Args) > 2 && isValidSeason(os.Args[2]) {
				seasonYear, err := strconv.Atoi(os.Args[2])
				if err != nil {
					log.Fatalf("error: invalid season number: %s using default season 2026", os.Args[2])
				}
				if err := fetch.FetchFullSeason(gormDB, seasonYear); err != nil {
					log.Fatalf("fetch failed: %v for year: %d", err, seasonYear)
				}
				break
			}
			if err := fetch.FetchFullSeason(gormDB, 2026); err != nil {
				log.Fatalf("fetch failed: %v for default season 2026", err)
			}
			log.Println("fetch completed successfully.")
			os.Exit(0) // exit immediately

		case "--fetch-date":
			if len(os.Args) < 3 {
				log.Fatalf("error: --fetch-date requires a date argument (yyyy-mm-dd)")
			}
			date := os.Args[2]
			log.Printf("updating games for: %s", date)
			if err := fetch.FetchGamesByDate(gormDB, date); err != nil {
				log.Fatalf("update failed: %v", err)
			}
			log.Println("update completed successfully.")
			os.Exit(0) // exit immediately
		}
	}

	// if no fetch flags were matched, start the web server
	router := routes.NewRouter(gormDB)
	fmt.Printf("Server running on: http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

func isValidSeason(season string) bool {
	validSeasons := []string{"2010", "2011", "2012", "2013", "2014", "2015", "2016", "2017", "2018", "2019", "2020", "2021", "2022", "2023", "2024", "2025", "2026"}
	for _, valid := range validSeasons {
		if valid == season {
			return true
		}
	}
	return false
}
