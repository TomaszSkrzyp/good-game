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
}

func calculateFinalQuality(hScore, aScore int, hQs, aQs []int, homeLeaders, awayLeaders []ESPNLeaderCategory) models.GameQuality {
	var score float64 = 0
	quality := models.GameQuality{}

	// final Margin Logic
	margin := int(math.Abs(float64(hScore - aScore)))
	if margin <= 3 {
		score += 40
	} else if margin <= 7 {
		score += 25
	} else if margin <= 12 {
		score += 10
	}

	// Clutch 4th Quarter Analysis
	hRunning, aRunning := 0, 0
	for i := 0; i < 3 && i < len(hQs); i++ {
		hRunning += hQs[i]
		aRunning += aQs[i]
	}
	fourthStartMargin := int(math.Abs(float64(hRunning - aRunning)))
	if fourthStartMargin <= 8 || margin <= 3 {
		quality.IsClutch = true
		score += 20
	}

	// Individual Brilliance (Big Scoring & Big Game)
	playerVersatility := make(map[string]int)
	playerPoints := make(map[string]float64)
	allCategories := append(homeLeaders, awayLeaders...)

	for _, cat := range allCategories {
		for _, l := range cat.Leaders {
			pID := l.Athlete.DisplayName
			val := l.Value

			if cat.Name == "points" {
				playerPoints[pID] = val
				if val > 35 {
					quality.IsBigScoring = true
				}
				if val >= 50 {
					quality.IsBigGame = true
				}
			}

			if cat.Name == "points" || cat.Name == "rebounds" || cat.Name == "assists" {
				playerVersatility[pID]++
			}
		}
	}

	// Multi-category leadership check for "Big Game"
	for pID, count := range playerVersatility {
		if count >= 3 && playerPoints[pID] >= 35 {
			quality.IsBigGame = true
		}
	}

	// Apply boosts to the quality score
	if quality.IsBigScoring {
		score += 15
	}
	if quality.IsBigGame {
		score += 25
	}

	if score > 100 {
		score = 100
	}
	quality.QualityScore = uint(score)

	return quality
}
