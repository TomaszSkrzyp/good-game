package db

import (
	"context"
	"log"
	"os"

	"github.com/tomaszSkrzyp/good-game/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedRoles(ctx context.Context, db *gorm.DB) {
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

func SeedAdminUser(ctx context.Context, db *gorm.DB) {
	adminUser := os.Getenv("ADMIN_USERNAME")
	adminPass := os.Getenv("ADMIN_PASSWORD")

	if adminUser == "" || adminPass == "" {
		log.Println("Notice: ADMIN_USERNAME or ADMIN_PASSWORD not set. Skipping admin seed.")
		return
	}

	var admin models.User
	if err := db.WithContext(ctx).Where("user_name = ?", adminUser).First(&admin).Error; err == gorm.ErrRecordNotFound {
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

func BuildConferenceMap(ctx context.Context, db *gorm.DB) error {
	r := NewConferenceRepository(db)
	m, err := r.NameToIDMap()
	if err != nil {
		return err
	}
	ConferenceNameToIDMap = m
	return nil
}

var TeamAbbrToIDMap = make(map[string]uint)

func BuildTeamMap(ctx context.Context, db *gorm.DB) error {
	var teams []models.Team
	if err := db.WithContext(ctx).Find(&teams).Error; err != nil {
		return err
	}

	for _, t := range teams {
		TeamAbbrToIDMap[t.Abbreviation] = t.ID
	}
	return nil
}

func SeedTeams(ctx context.Context, db *gorm.DB) {
	data := map[string][]struct{ Name, Abbr string }{
		"Eastern Conference": {
			{"Atlanta Hawks", "ATL"}, {"Boston Celtics", "BOS"}, {"Brooklyn Nets", "BKN"},
			{"Charlotte Hornets", "CHA"}, {"Chicago Bulls", "CHI"}, {"Cleveland Cavaliers", "CLE"},
			{"Detroit Pistons", "DET"}, {"Indiana Pacers", "IND"}, {"Miami Heat", "MIA"},
			{"Milwaukee Bucks", "MIL"}, {"New York Knicks", "NY"}, {"Orlando Magic", "ORL"},
			{"Philadelphia 76ers", "PHI"}, {"Toronto Raptors", "TOR"}, {"Washington Wizards", "WSH"},
			{"TBD", "TBD"},
		},
		"Western Conference": {
			{"Dallas Mavericks", "DAL"}, {"Denver Nuggets", "DEN"}, {"Golden State Warriors", "GS"},
			{"Houston Rockets", "HOU"}, {"Los Angeles Clippers", "LAC"}, {"Los Angeles Lakers", "LAL"},
			{"Memphis Grizzlies", "MEM"}, {"Minnesota Timberwolves", "MIN"}, {"New Orleans Pelicans", "NO"},
			{"Oklahoma City Thunder", "OKC"}, {"Phoenix Suns", "PHX"}, {"Portland Trail Blazers", "POR"},
			{"Sacramento Kings", "SAC"}, {"San Antonio Spurs", "SA"}, {"Utah Jazz", "UTAH"},
		},
	}

	for confName, teams := range data {
		var conf models.Conference
		db.FirstOrCreate(&conf, models.Conference{Name: confName})

		for _, t := range teams {
			db.WithContext(ctx).Where(models.Team{Name: t.Name}).FirstOrCreate(&models.Team{
				Name:         t.Name,
				Abbreviation: t.Abbr,
				ConferenceID: conf.ID,
			})
		}
	}
}

var defaultGameConfig = models.GameQualityConfig{
	Margins: []models.MarginWeight{
		{MaxMargin: 3, Points: 45},
		{MaxMargin: 7, Points: 30},
		{MaxMargin: 12, Points: 15},
	},
	HugeSwingBonus:      25,
	ClutchBonus:         20,
	OvertimeBonus:       15,
	ShootoutBonus:       15,
	ShootoutThreshold:   235,
	GrittyThreshold:     200,
	StarDuelBonus:       20,
	StarPointsThreshold: 35,
	BigGameBonus:        15,
}
