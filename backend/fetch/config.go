package fetch

type MarginWeight struct {
	MaxMargin int `json:"maxMargin"`
	Points    int `json:"points"`
}

type GameQualityConfig struct {
	Margins             []MarginWeight `json:"margins"`
	HugeSwingBonus      int            `json:"hugeSwingBonus"`
	ClutchBonus         int            `json:"clutchBonus"`
	OvertimeBonus       int            `json:"overtimeBonus"`
	ShootoutBonus       int            `json:"shootoutBonus"`
	ShootoutThreshold   int            `json:"shootoutThreshold"`
	GrittyThreshold     int            `json:"grittyThreshold"`
	StarDuelBonus       int            `json:"starDuelBonus"`
	StarPointsThreshold int            `json:"starPointsThreshold"`
	BigGameBonus        int            `json:"bigGameBonus"`
}

var CurrentConfig = GameQualityConfig{
	Margins: []MarginWeight{
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
