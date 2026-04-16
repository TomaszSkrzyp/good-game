package db

import (
	"log"
	"os"

	"github.com/tomaszSkrzyp/good-game/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	roles := []models.Role{
		{Name: "Admin"},
		{Name: "User"},
	}

	for _, role := range roles {
		var existing models.Role
		if err := db.Where("name = ?", role.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&role).Error; err != nil {
					log.Fatalf("Failed to seed role %s: %v", role.Name, err)
				}
			} else {
				log.Fatalf("Failed to query roles: %v", err)
			}
		}
	}
}

func SeedAdminUser(db *gorm.DB) {
	adminUser := os.Getenv("ADMIN_USERNAME")
	adminPass := os.Getenv("ADMIN_PASSWORD")

	if adminUser == "" || adminPass == "" {
		log.Println("Notice: ADMIN_USERNAME or ADMIN_PASSWORD not set. Skipping admin seed.")
		return
	}

	var admin models.User
	if err := db.Where("user_name = ?", adminUser).First(&admin).Error; err == gorm.ErrRecordNotFound {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash admin password: %v", err)
		}
		admin = models.User{
			UserName: adminUser,
			Password: string(hashedPassword),
			RoleID:   1,
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Fatalf("Failed to create admin user: %v", err)
		} else {
			log.Printf("Admin user '%s' seeded successfully.", adminUser)
		}
	}
}
