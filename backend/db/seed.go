package db

import (
	"fmt"
	"log"

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
				fmt.Println("Created role:", role.Name)
			} else {
				log.Fatalf("Failed to query roles: %v", err)
			}
		}
	}
}
func SeedAdminUser(db *gorm.DB) {
	var admin models.User
	if err := db.Where("user_name = ?", "admin").First(&admin).Error; err == gorm.ErrRecordNotFound {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		admin = models.User{
			UserName: "admin",
			Password: string(hashedPassword),
			RoleID:   1,
		}
		db.Create(&admin)
		fmt.Println("Admin user created with username: admin, password: admin")
	}
}
