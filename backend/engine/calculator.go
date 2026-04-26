package engine

import (
	"math"

	"github.com/tomaszSkrzyp/good-game/models"
)

// CalculateFinalQuality now takes drama context fetched from the summary API
func CalculateFinalQuality(
	hScore, aScore int,
	hQs, aQs []int,
	homeLeaders, awayLeaders []ESPNLeaderCategory,
	dramaContext DramaContext,
) models.GameQuality {

	var score float64 = dramaContext.DramaScore
	quality := models.GameQuality{
		IsHugeSwing:   dramaContext.IsHugeComeback,
		IsGame7:       dramaContext.IsGame7,
		IsElimination: dramaContext.IsElimination,
		IsPlayoff:     dramaContext.IsPlayoff,
		IsPlayIn:      dramaContext.IsPlayIn,
	}

	cfg := models.GetConfig()

	margin := int(math.Abs(float64(hScore - aScore)))
	for _, m := range cfg.Margins {
		if margin <= m.MaxMargin {
			score += float64(m.Points)
			break
		}
	}

	hRunning, aRunning := 0, 0
	for i := 0; i < 3 && i < len(hQs); i++ {
		hRunning += hQs[i]
		aRunning += aQs[i]
	}
	thirdQtrMargin := int(math.Abs(float64(hRunning - aRunning)))

	if thirdQtrMargin >= 15 && margin <= 7 {
		quality.IsHugeSwing = true
		score += float64(cfg.HugeSwingBonus)
	}

	if thirdQtrMargin <= 8 || margin <= 3 {
		quality.IsClutch = true
		score += float64(cfg.ClutchBonus)
	}

	if len(hQs) > 4 {
		quality.IsOvertime = true
		score += float64(cfg.OvertimeBonus)
	}

	totalPoints := hScore + aScore
	if totalPoints >= cfg.ShootoutThreshold {
		quality.IsShootout = true
		score += float64(cfg.ShootoutBonus)
	} else if totalPoints <= cfg.GrittyThreshold && !quality.IsOvertime {
		quality.IsGritty = true
		score += 10
	}

	playerPoints := make(map[string]float64)
	playerVersatility := make(map[string]int)
	homeHasStar, awayHasStar := false, false

	processLeaders := func(leaders []ESPNLeaderCategory, isHome bool) {
		for _, cat := range leaders {
			for _, l := range cat.Leaders {
				pID := l.Athlete.DisplayName
				if cat.Name == "points" {
					playerPoints[pID] = l.Value
					if l.Value >= float64(cfg.StarPointsThreshold) {
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

	if homeHasStar && awayHasStar {
		quality.IsStarDuel = true
		score += float64(cfg.StarDuelBonus)
	}

	for pID, count := range playerVersatility {
		if count >= 3 && playerPoints[pID] >= 30 {
			quality.IsBigGame = true
			score += float64(cfg.BigGameBonus)
			break
		}
	}

	if score > 100 {
		score = 100
	}
	quality.QualityScore = uint(score)

	return quality
}
