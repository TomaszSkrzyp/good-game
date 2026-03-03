package fetch

import (
	"math"

	"github.com/tomaszSkrzyp/good-game/models"
)

type Gamequality struct {
	Quality      uint
	IsBigScoring bool
	IsBigGame    bool
	IsClutch     bool
	IsStarDuel   bool
	IsHugeSwing  bool
	IsShootout   bool
	IsGritty     bool
}

func calculateFinalQuality(hScore, aScore int, hQs, aQs []int, homeLeaders, awayLeaders []ESPNLeaderCategory) models.GameQuality {
	var score float64 = 0
	quality := models.GameQuality{}

	// Final margin calculation
	margin := int(math.Abs(float64(hScore - aScore)))
	if margin <= 3 {
		score += 45
	} else if margin <= 7 {
		score += 30
	} else if margin <= 12 {
		score += 15
	}

	// Score after 3 quarters to detect huge leads
	hRunning, aRunning := 0, 0
	for i := 0; i < 3 && i < len(hQs); i++ {
		hRunning += hQs[i]
		aRunning += aQs[i]
	}
	thirdQtrMargin := int(math.Abs(float64(hRunning - aRunning)))

	// Huge swing: big lead after 3Q that ended in a close game
	if thirdQtrMargin >= 15 && margin <= 7 {
		quality.IsHugeSwing = true
		score += 25
	}

	// Clutch and Overtime logic
	if thirdQtrMargin <= 8 || margin <= 3 {
		quality.IsClutch = true
		score += 20
	}
	if len(hQs) > 4 {
		score += 15
	}

	// Game style flags
	totalPoints := hScore + aScore
	if totalPoints >= 235 {
		quality.IsShootout = true
		score += 15
	} else if totalPoints <= 200 && len(hQs) == 4 {
		quality.IsGritty = true
		score += 10
	}

	playerPoints := make(map[string]float64)
	playerVersatility := make(map[string]int)
	homeHasStar := false
	awayHasStar := false

	// Internal helper for leader processing
	processLeaders := func(leaders []ESPNLeaderCategory, isHome bool) {
		for _, cat := range leaders {
			for _, l := range cat.Leaders {
				pID := l.Athlete.DisplayName
				if cat.Name == "points" {
					playerPoints[pID] = l.Value
					if l.Value >= 35 {
						quality.IsBigScoring = true
						if isHome {
							homeHasStar = true
						} else {
							awayHasStar = true
						}
					}
					if l.Value >= 50 {
						quality.IsBigGame = true
					}
				}
				if cat.Name == "points" || cat.Name == "rebounds" || cat.Name == "assists" {
					playerVersatility[pID]++
				}
			}
		}
	}

	processLeaders(homeLeaders, true)
	processLeaders(awayLeaders, false)

	// Star Duel: high scorers on both sides
	if homeHasStar && awayHasStar {
		quality.IsStarDuel = true
		score += 20
	}

	// Big Game: based on multi-category performance
	for pID, count := range playerVersatility {
		if count >= 3 && playerPoints[pID] >= 30 {
			quality.IsBigGame = true
			score += 15
			break
		}
	}

	if score > 100 {
		score = 100
	}
	quality.QualityScore = uint(score)

	return quality
}
