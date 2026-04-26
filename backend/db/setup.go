package db

import (
	"context"
	"fmt"
	"log"

	"github.com/tomaszSkrzyp/good-game/models"
	"gorm.io/gorm"
)

func Initialize(ctx context.Context, db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Game{},
		&models.Team{},
		&models.Conference{},
		&models.TeamStats{},
		&models.UserReaction{},
		&models.ConfigRecord{},
	)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	SeedRoles(ctx, db)
	SeedAdminUser(ctx, db)
	SeedTeams(ctx, db)

	if err := SyncAlgorithmConfig(ctx, db); err != nil {
		log.Printf("Warning: Could not sync algorithm config: %v", err)
	}

	if err := BuildConferenceMap(ctx, db); err != nil {
		log.Printf("Warning: BuildConferenceMap failed: %v", err)
	}

	if err := BuildTeamMap(ctx, db); err != nil {
		log.Printf("Warning: BuildTeamMap failed: %v", err)
	}

	return nil
}
