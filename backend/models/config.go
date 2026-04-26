package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"log"
	"sync"
)

type ConfigRecord struct {
	ID   uint              `gorm:"primaryKey"`
	Key  string            `gorm:"uniqueIndex"`
	Data GameQualityConfig `gorm:"type:jsonb"`
}

type GameQualityConfig struct {
	// core win probability score
	VolatilityWeight float64 `json:"volatilityWeight"`
	SwingWeight      float64 `json:"swingWeight"`

	// continuous margin calculation
	MaxMarginBonus  float64 `json:"maxMarginBonus"`
	MarginDecayRate float64 `json:"marginDecayRate"`

	// player performance curve
	StarPointsBase       float64 `json:"starPointsBase"`
	StarPointsMultiplier float64 `json:"starPointsMultiplier"`

	// event stakes
	OvertimeBonus         float64 `json:"overtimeBonus"`
	EliminationBonus      float64 `json:"eliminationBonus"`
	Game7Bonus            float64 `json:"game7Bonus"`
	SeasonSeriesTiedBonus float64 `json:"seasonSeriesTiedBonus"`
	StarDuelBonus         float64 `json:"starDuelBonus"`

	// ui tagging flags
	ShootoutThreshold int     `json:"shootoutThreshold"`
	GrittyThreshold   int     `json:"grittyThreshold"`
	ComebackThreshold float64 `json:"comebackThreshold"`
	MaxScore          float64 `json:"maxScore"`

	ProbFlipWeight    float64 `json:"ProbFlipWeight"`
	StolenGameMaxLead float64 `json:"stolenGameMaxLead"`
	StolenGameWeight  float64 `json:"stolenGameWeight"`
}

func (c *GameQualityConfig) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &c)
}

func (c GameQualityConfig) Value() (driver.Value, error) {
	return json.Marshal(c)
}

var (
	configMu      sync.RWMutex
	currentConfig GameQualityConfig
)

func GetConfig() GameQualityConfig {
	configMu.RLock()
	defer configMu.RUnlock()
	return currentConfig
}

func SetConfig(cfg GameQualityConfig) {
	configMu.Lock()
	defer configMu.Unlock()
	log.Println("setting config to ", cfg)
	currentConfig = cfg
}
