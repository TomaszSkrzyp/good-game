package db

import (
	"log"

	"github.com/tomaszSkrzyp/good-game/models"
	"gorm.io/gorm"
)

func Initialize(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Game{},
		&models.Team{},
		&models.Conference{},
		&models.TeamStats{},
		&models.UserReaction{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	SeedRoles(db)
	SeedAdminUser(db)
	SeedTeams(db)
	if err := BuildConferenceMap(db); err != nil {
		log.Printf("Warning: BuildConferenceMap failed: %v", err)
	}

	if err := BuildTeamMap(db); err != nil {
		log.Printf("Warning: BuildTeamMap failed: %v", err)
	}
}
