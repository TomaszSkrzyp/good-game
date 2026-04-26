package engine

// ESPN API Schema
type ESPNResponse struct {
	Events []ESPNEvent `json:"events"`
}

type ESPNEvent struct {
	ID           string            `json:"id"`
	Date         string            `json:"date"`
	Status       ESPNStatus        `json:"status"`
	Competitions []ESPNCompetition `json:"competitions"`
}

type ESPNStatus struct {
	Type struct {
		Name        string `json:"name"`        // e.g., STATUS_FINAL
		Description string `json:"description"` // e.g., "Final"
	} `json:"type"`
}

type ESPNCompetition struct {
	Competitors []ESPNCompetitor `json:"competitors"`
}
type ESPNLineScore struct {
	Value float64 `json:"value"` // Score in that specific quarter
}

type ESPNCompetitor struct {
	HomeAway   string               `json:"homeAway"`
	Score      string               `json:"score"`
	Team       ESPNTeam             `json:"team"`
	Leaders    []ESPNLeaderCategory `json:"leaders"`
	LineScores []ESPNLineScore      `json:"linescores"`
}

type ESPNTeam struct {
	Abbreviation string `json:"abbreviation"`
}

type ESPNLeaderCategory struct {
	Name    string       `json:"name"` // "points", "rebounds", "assists"
	Leaders []ESPNLeader `json:"leaders"`
}

type ESPNLeader struct {
	Value   float64 `json:"value"`
	Athlete struct {
		DisplayName string `json:"displayName"`
	} `json:"athlete"`
}

type GamePlayerStatsDTO struct {
	HomeTopScorer    string  `json:"homeTopScorer"`
	HomeTopScorerPts float64 `json:"homeTopScorerPts"`
	HomeTopAssister  string  `json:"homeTopAssister"`
	HomeTopAssists   float64 `json:"homeTopAssists"`
	HomeTopRebounder string  `json:"homeTopRebounder"`
	HomeTopRebounds  float64 `json:"homeTopRebounds"`
	AwayTopScorer    string  `json:"awayTopScorer"`
	AwayTopScorerPts float64 `json:"awayTopScorerPts"`
	AwayTopAssister  string  `json:"awayTopAssister"`
	AwayTopAssists   float64 `json:"awayTopAssists"`
	AwayTopRebounder string  `json:"awayTopRebounder"`
	AwayTopRebounds  float64 `json:"awayTopRebounds"`
}
