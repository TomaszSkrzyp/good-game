package engine

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/tomaszSkrzyp/good-game/models"
)

type WinProbabilityData struct {
	WinProbability []struct {
		HomeWinPercentage float64 `json:"homeWinPercentage"`
	} `json:"winprobability"`
	Header struct {
		Season struct {
			Type int `json:"type"`
		} `json:"season"`
	} `json:"header"`
	SeasonSeries []struct {
		Type        string `json:"type"`
		SeriesScore string `json:"seriesScore"`
	} `json:"seasonseries"`
	Boxscore struct {
		Teams []struct {
			HomeAway   string `json:"homeAway"`
			Statistics []struct {
				Name         string `json:"name"`
				DisplayValue string `json:"displayValue"`
			} `json:"statistics"`
		} `json:"teams"`
	} `json:"boxscore"`
}

type DramaContext struct {
	DramaScore     float64
	IsClutchEnding bool
	IsHugeComeback bool
	IsElimination  bool
	IsGame7        bool
	IsPlayoff      bool
	IsPlayIn       bool
}

func FetchAndCalculateDrama(eventID string) (DramaContext, error) {
	url := fmt.Sprintf("http://site.api.espn.com/apis/site/v2/sports/basketball/nba/summary?event=%s", eventID)
	resp, err := espnClient.Get(url)
	if err != nil {
		return DramaContext{}, err
	}
	defer resp.Body.Close()

	var data WinProbabilityData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return DramaContext{}, err
	}
	ctx := DramaContext{}
	cfg := models.GetConfig()

	homeLeadPct := -1.0
	awayLeadPct := -1.0

	for _, team := range data.Boxscore.Teams {
		isHome := team.HomeAway == "home"
		for _, stat := range team.Statistics {
			if stat.Name == "leadPercentage" {
				pct, _ := strconv.ParseFloat(stat.DisplayValue, 64)
				if isHome {
					homeLeadPct = pct
				} else {
					awayLeadPct = pct
				}
			}
		}
	}
	for _, series := range data.SeasonSeries {
		isPostseason := data.Header.Season.Type == 3
		if (series.Type == "playoff" || series.Type == "play-in") && isPostseason {
			if series.Type == "playoff" {
				ctx.IsPlayoff = true
			}
			if series.Type == "play-in" {
				ctx.IsPlayIn = true
			}

			scores := strings.Split(series.SeriesScore, "-")
			if len(scores) == 2 {
				winA, _ := strconv.Atoi(scores[0])
				winB, _ := strconv.Atoi(scores[1])

				// Game 7 and Elimination bonuses are additive
				if (winA == 3 && winB == 3) || (winA+winB == 7) {
					ctx.IsGame7 = true
					ctx.IsElimination = true
					ctx.DramaScore += cfg.Game7Bonus
					ctx.DramaScore += cfg.EliminationBonus
				} else if winA == 3 || winB == 3 || winA == 4 || winB == 4 {
					ctx.IsElimination = true
					ctx.DramaScore += cfg.EliminationBonus
				} else if series.Type == "play-in" {
					ctx.IsElimination = true
				}
			}
		}

		if series.Type == "season" {
			scores := strings.Split(series.SeriesScore, "-")
			if len(scores) == 2 {
				winA, _ := strconv.Atoi(scores[0])
				winB, _ := strconv.Atoi(scores[1])

				if winA == 1 || winB == 1 {
					ctx.DramaScore += cfg.SeasonSeriesTiedBonus

				}
			}
		}
	}
	probs := data.WinProbability
	if len(probs) >= 2 {
		var totalVolatility float64
		var probFlips int
		var lateFlips int
		var lateTensionTicks int
		var sumProb float64
		minProb, maxProb := 1.0, 0.0

		// Check the last 10% of the game for late-game flips and tension
		clutchThreshold := int(float64(len(probs)) * 0.9)
		finalTicks := len(probs) - clutchThreshold

		for i := 0; i < len(probs); i++ {
			p := probs[i].HomeWinPercentage
			sumProb += p

			if i >= clutchThreshold {
				if p >= 0.40 && p <= 0.60 {
					lateTensionTicks++
				}
			}

			if p < minProb {
				minProb = p
			}
			if p > maxProb {
				maxProb = p
			}

			if i > 0 {
				prevP := probs[i-1].HomeWinPercentage
				totalVolatility += math.Abs(p - prevP)

				if (prevP > 0.5 && p < 0.5) || (prevP < 0.5 && p > 0.5) {
					probFlips++
					if i >= clutchThreshold {
						lateFlips++
					}
				}
			}
		}

		avgHomeProb := sumProb / float64(len(probs))
		finalProb := probs[len(probs)-1].HomeWinPercentage

		homeWon := finalProb > 0.5
		awayWon := finalProb < 0.5

		var winnerAvgProb float64
		var winnerMinProb float64

		if homeWon {
			winnerAvgProb = avgHomeProb
			winnerMinProb = minProb
		} else {
			winnerAvgProb = 1.0 - avgHomeProb
			winnerMinProb = 1.0 - maxProb
		}

		comebackAmount := 1.0 - winnerMinProb

		isDeepDeficit := comebackAmount >= cfg.ComebackThreshold
		isSustainedDeficit := winnerAvgProb <= 0.35

		ctx.IsHugeComeback = isDeepDeficit && isSustainedDeficit

		totalSwing := maxProb - minProb
		probFlipBonus := float64(probFlips) * cfg.ProbFlipWeight

		ctx.DramaScore += (totalVolatility * cfg.VolatilityWeight) +
			(totalSwing * cfg.SwingWeight) +
			probFlipBonus
		isTenseEnding := finalTicks > 0 && (float64(lateTensionTicks)/float64(finalTicks)) >= 0.30
		ctx.IsClutchEnding = lateFlips > 0 || isTenseEnding

		var winnerLeadPct float64
		winnerLeadPct = -1.0

		if homeWon && homeLeadPct >= 0 {
			winnerLeadPct = homeLeadPct
		} else if awayWon && awayLeadPct >= 0 {
			winnerLeadPct = awayLeadPct
		}

		if winnerLeadPct >= 0 {
			stolenBonus := math.Max(cfg.StolenGameMaxLead-winnerLeadPct, 0) * cfg.StolenGameWeight

			if stolenBonus > 0 {
				ctx.DramaScore += stolenBonus

				if winnerLeadPct < (cfg.StolenGameMaxLead / 2) {
					ctx.IsHugeComeback = true
				}
			}
		}
	}
	return ctx, nil
}
