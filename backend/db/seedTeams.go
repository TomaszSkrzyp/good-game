package db

import (
	"fmt"

	"github.com/tomaszSkrzyp/good-game/models"
	"gorm.io/gorm"
)

func SeedConferences(db *gorm.DB) {
	var easternConference models.Conference
	easternConference = models.Conference{
		Name: "Eastern Conference",
	}
	if err := db.Create(&easternConference).Error; err != nil {
		fmt.Println("Failed to seed Eastern Conference:", err)
	} else {
		fmt.Println("Seeded Eastern Conference")
	}

	var westernConference models.Conference
	westernConference = models.Conference{
		Name: "Western Conference",
	}
	if err := db.Create(&westernConference).Error; err != nil {
		fmt.Println("Failed to seed Western Conference:", err)
	} else {
		fmt.Println("Seeded Western Conference")
	}
}

func BuildConferenceMap(db *gorm.DB) error {
	r := NewConferenceRepository(db)
	m, err := r.NameToIDMap()
	if err != nil {
		return err
	}
	ConferenceNameToIDMap = m
	return nil
}

func SeedTeams(db *gorm.DB) {
	conferenceTeams := map[string][]string{
		"Eastern Conference": {
			"Boston Celtics", "Brooklyn Nets", "New York Knicks", "Philadelphia 76ers", "Toronto Raptors",
			"Chicago Bulls", "Cleveland Cavaliers", "Detroit Pistons", "Indiana Pacers", "Milwaukee Bucks",
			"Atlanta Hawks", "Charlotte Hornets", "Miami Heat", "Orlando Magic", "Washington Wizards",
		},
		"Western Conference": {
			"Denver Nuggets", "Minnesota Timberwolves", "Oklahoma City Thunder", "Portland Trail Blazers", "Utah Jazz",
			"Golden State Warriors", "Los Angeles Clippers", "Los Angeles Lakers", "Phoenix Suns", "Sacramento Kings",
			"Dallas Mavericks", "Houston Rockets", "Memphis Grizzlies", "New Orleans Pelicans", "San Antonio Spurs",
		},
	}

	for confName, teams := range conferenceTeams {
		var conf models.Conference
		if err := db.Where("name = ?", confName).First(&conf).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// create conference if missing
				conf = models.Conference{Name: confName}
				if err := db.Create(&conf).Error; err != nil {
					fmt.Println("Failed to create conference:", confName, err)
					continue
				}
			} else {
				fmt.Println("Conference query error:", confName, err)
				continue
			}
		}

		for _, tname := range teams {
			var existing models.Team
			if err := db.Where("name = ? AND conference_id = ?", tname, conf.ID).First(&existing).Error; err == nil {
				continue
			}
			team := models.Team{
				Name:         tname,
				ConferenceID: conf.ID,
			}
			if err := db.Create(&team).Error; err != nil {
				fmt.Println("Failed to create team:", tname, err)
			} else {
				fmt.Println("Created team:", tname, "in", confName)
			}
		}
	}
}

func BuildTeamMap(db *gorm.DB) error {
	var teams []models.Team
	if err := db.Find(&teams).Error; err != nil {
		return err
	}
	TeamNameToIDMap = make(TeamNameToID, len(teams))
	for _, t := range teams {
		TeamNameToIDMap[t.Name] = t.ID
	}
	return nil
}
