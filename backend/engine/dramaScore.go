package engine

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/tomaszSkrzyp/good-game/models"
)

type WinProbabilityData struct {
	WinProbability []struct {
		HomeWinPercentage float64 `json:"homeWinPercentage"`
	} `json:"winprobability"`
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
	IsHugeComeback bool
	IsElimination  bool
	IsGame7        bool
	IsPlayoff      bool
	IsPlayIn       bool
}

func FetchAndCalculateDrama(eventID string) (DramaContext, error) {
	url := fmt.Sprintf("http://site.api.espn.com/apis/site/v2/sports/basketball/nba/summary?event=%s", eventID)
	resp, err := http.Get(url)
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
		if series.Type == "playoff" || series.Type == "play-in" {
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

				if winA == 3 && winB == 3 {
					ctx.IsGame7 = true
					ctx.DramaScore += cfg.Game7Bonus
				} else if winA == 3 || winB == 3 {
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
		var probFlips int //track how many times the favorite changes
		minProb, maxProb := 1.0, 0.0

		for i := 0; i < len(probs); i++ {
			p := probs[i].HomeWinPercentage
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
				}
			}
		}

		totalSwing := maxProb - minProb
		ctx.IsHugeComeback = totalSwing > cfg.ComebackThreshold

		probFlipBonus := float64(probFlips) * cfg.ProbFlipWeight

		ctx.DramaScore += (totalVolatility * cfg.VolatilityWeight) +
			(totalSwing * cfg.SwingWeight) +
			probFlipBonus

		finalProb := probs[len(probs)-1].HomeWinPercentage
		homeWon := finalProb > 0.5
		awayWon := finalProb < 0.5

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
