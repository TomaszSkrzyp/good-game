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
		Type        string `json:"type"`        // "playoff" or "season"
		SeriesScore string `json:"seriesScore"` // "2-1" or "3-3"
	} `json:"seasonseries"`
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

	for _, series := range data.SeasonSeries {
		fmt.Println(series.Type)
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
					ctx.DramaScore += float64(cfg.Game7Bonus)
				} else if winA == 3 || winB == 3 {
					ctx.IsElimination = true
					ctx.DramaScore += float64(cfg.EliminationBonus)
				} else if series.Type == "play-in" {
					ctx.IsElimination = true
					ctx.DramaScore += float64(cfg.PlayInBonus)
				} else {
					ctx.DramaScore += float64(cfg.PlayoffBonus)
				}
			}
		}

		if series.Type == "season" {
			scores := strings.Split(series.SeriesScore, "-")
			if len(scores) == 2 {
				winA, _ := strconv.Atoi(scores[0])
				winB, _ := strconv.Atoi(scores[1])

				if winA == 1 || winB == 1 {
					ctx.DramaScore += float64(cfg.SeasonSeriesTiedBonus)
				}
			}
		}
	}

	probs := data.WinProbability
	if len(probs) >= 2 {
		var totalVolatility float64
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
				totalVolatility += math.Abs(p - probs[i-1].HomeWinPercentage)
			}
		}

		totalSwing := maxProb - minProb
		ctx.IsHugeComeback = totalSwing > cfg.ComebackThreshold
		ctx.DramaScore += (totalVolatility * cfg.VolatilityWeight) + (totalSwing * cfg.SwingWeight)
	}

	return ctx, nil
}
