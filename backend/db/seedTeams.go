package db

import (
	"github.com/tomaszSkrzyp/good-game/models"
	"gorm.io/gorm"
)

func BuildConferenceMap(db *gorm.DB) error {
	r := NewConferenceRepository(db)
	m, err := r.NameToIDMap()
	if err != nil {
		return err
	}
	ConferenceNameToIDMap = m
	return nil
}

var TeamAbbrToIDMap = make(map[string]uint)

func BuildTeamMap(db *gorm.DB) error {
	var teams []models.Team
	if err := db.Find(&teams).Error; err != nil {
		return err
	}

	// Czyścimy i wypełniamy mapę skrótami (np. "LAL" -> 1)
	for _, t := range teams {
		TeamAbbrToIDMap[t.Abbreviation] = t.ID
	}
	return nil
}

func SeedTeams(db *gorm.DB) {
	data := map[string][]struct{ Name, Abbr string }{
		"Eastern Conference": {
			{"Atlanta Hawks", "ATL"}, {"Boston Celtics", "BOS"}, {"Brooklyn Nets", "BKN"},
			{"Charlotte Hornets", "CHA"}, {"Chicago Bulls", "CHI"}, {"Cleveland Cavaliers", "CLE"},
			{"Detroit Pistons", "DET"}, {"Indiana Pacers", "IND"}, {"Miami Heat", "MIA"},
			{"Milwaukee Bucks", "MIL"}, {"New York Knicks", "NY"}, {"Orlando Magic", "ORL"},
			{"Philadelphia 76ers", "PHI"}, {"Toronto Raptors", "TOR"}, {"Washington Wizards", "WSH"},
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
			db.Where(models.Team{Name: t.Name}).FirstOrCreate(&models.Team{
				Name:         t.Name,
				Abbreviation: t.Abbr,
				ConferenceID: conf.ID,
			})
		}
	}
}
