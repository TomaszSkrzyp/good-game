package db

import (
	"fmt"
	"log"

	"github.com/tomaszSkrzyp/good-game/models"
	"gorm.io/gorm"
)

func HardReset(db *gorm.DB) error {
	log.Println("Starting full database wipe...")

	tables := []interface{}{
		&models.UserReaction{},
		&models.TeamStats{},
		&models.Game{},
		&models.Team{},
		&models.Conference{},
		&models.ConfigRecord{},
		&models.User{},
		&models.Role{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			return fmt.Errorf("failed to drop table: %w", err)
		}
	}

	log.Println("Database wiped. Re-initializing...")
	return nil
}
