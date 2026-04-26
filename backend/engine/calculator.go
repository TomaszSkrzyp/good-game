package engine

import (
	"math"

	"github.com/tomaszSkrzyp/good-game/models"
)

func CalculateFinalQuality(
	hScore, aScore int,
	hQs, aQs []int,
	homeLeaders, awayLeaders []ESPNLeaderCategory,
	dramaContext DramaContext,
) models.GameQuality {

	var score = dramaContext.DramaScore
	quality := models.GameQuality{
		IsHugeSwing:   dramaContext.IsHugeComeback,
		IsGame7:       dramaContext.IsGame7,
		IsElimination: dramaContext.IsElimination,
		IsPlayoff:     dramaContext.IsPlayoff,
		IsPlayIn:      dramaContext.IsPlayIn,
	}

	cfg := models.GetConfig()

	// continuous margin score, 0 if blowout
	margin := math.Abs(float64(hScore - aScore))
	score += math.Max(0, cfg.MaxMarginBonus-(margin*cfg.MarginDecayRate))

	hRunning, aRunning := 0, 0
	for i := 0; i < 3 && i < len(hQs); i++ {
		hRunning += hQs[i]
		aRunning += aQs[i]
	}
	thirdQtrMargin := int(math.Abs(float64(hRunning - aRunning)))

	if thirdQtrMargin <= 8 || margin <= 3 {
		quality.IsClutch = true
	}

	if len(hQs) > 4 {
		quality.IsOvertime = true
		score += cfg.OvertimeBonus
	}

	totalPoints := hScore + aScore
	if totalPoints >= cfg.ShootoutThreshold {
		quality.IsShootout = true
	} else if totalPoints <= cfg.GrittyThreshold && !quality.IsOvertime {
		quality.IsGritty = true
	}

	playerPoints := make(map[string]float64)
	playerVersatility := make(map[string]int)
	homeHasStar, awayHasStar := false, false
	highestScorer := 0.0

	processLeaders := func(leaders []ESPNLeaderCategory, isHome bool) {
		for _, cat := range leaders {
			for _, l := range cat.Leaders {
				pID := l.Athlete.DisplayName
				if cat.Name == "points" {
					playerPoints[pID] = l.Value
					if l.Value > highestScorer {
						highestScorer = l.Value
					}
					if l.Value >= cfg.StarPointsBase {
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

	// linear scale for individual brilliance
	score += math.Max(0, (highestScorer-cfg.StarPointsBase)*cfg.StarPointsMultiplier)

	if homeHasStar && awayHasStar {
		quality.IsStarDuel = true
		score += cfg.StarDuelBonus
	}

	for pID, count := range playerVersatility {
		if count >= 3 && playerPoints[pID] >= 30 {
			quality.IsBigGame = true
			break
		}
	}

	// keep score within logical bounds
	if score > cfg.MaxScore {
		score = cfg.MaxScore
	}
	quality.QualityScore = uint(math.Round(score))

	return quality
}
