package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/models"
	"github.com/tomaszSkrzyp/good-game/routes"
)

func initDB() *gorm.DB {
	// SQLite file will be created in the project folder
	dbConn, err := gorm.Open(sqlite.Open("goodgame.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to SQLite: %v", err)
	}

	// Auto-migrate your models
	err = dbConn.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Game{},
		&models.Team{},
		&models.Conference{},
		&models.TeamStats{},
		&models.UserReaction{},
	)
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	fmt.Println("SQLite database connected and migrated")
	return dbConn
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	gormDB := initDB()
	db.SeedRoles(gormDB)
	db.SeedAdminUser(gormDB)
	err := db.BuildConferenceMap(gormDB)
	if err != nil {
		log.Fatalf("Failed to build conference map: %v", err)
	}
	db.SeedTeams(gormDB)
	err = db.BuildTeamMap(gormDB)
	if err != nil {
		log.Fatalf("Failed to build team map: %v", err)
	}
	// register routes by concern
	http.HandleFunc("/health", healthHandler)

	routes.RegisterConferenceRoutes(gormDB)
	routes.RegisterTeamRoutes(gormDB)
	routes.RegisterAuthRoutes(gormDB)
	routes.RegisterGameRoutes(gormDB)
	routes.RegisterTeamStatsRoutes(gormDB)
	routes.RegisterUserReactionRoutes(gormDB)

	fmt.Printf("Server running on port: %s\n", port)
	err = http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
