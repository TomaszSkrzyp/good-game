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

type MarginWeight struct {
	MaxMargin int `json:"maxMargin"`
	Points    int `json:"points"`
}

type GameQualityConfig struct {
	Margins              []MarginWeight `json:"margins"`
	HugeSwingBonus       int            `json:"hugeSwingBonus"`
	ClutchBonus          int            `json:"clutchBonus"`
	OvertimeBonus        int            `json:"overtimeBonus"`
	ShootoutBonus        int            `json:"shootoutBonus"`
	ShootoutThreshold    int            `json:"shootoutThreshold"`
	GrittyThreshold      int            `json:"grittyThreshold"`
	GrittyBonus          int            `json:"grittyBonus"`
	StarDuelBonus        int            `json:"starDuelBonus"`
	StarPointsThreshold  int            `json:"starPointsThreshold"`
	BigGameBonus         int            `json:"bigGameBonus"`
	BigScoringThreshold  int            `json:"bigScoringThreshold"`
	VersatilityThreshold int            `json:"versatilityThreshold"`

	// Playoff & Drama Config
	EliminationBonus  int     `json:"eliminationBonus"`
	Game7Bonus        int     `json:"game7Bonus"`
	PlayoffBonus      int     `json:"playoffBonus"`
	VolatilityWeight  float64 `json:"volatilityWeight"`
	SwingWeight       float64 `json:"swingWeight"`
	ComebackThreshold float64 `json:"comebackThreshold"`
	MaxScore          float64 `json:"maxScore"`

	PlayInBonus           int `json:"playInBonus"`
	SeasonSeriesTiedBonus int `json:"seasonSeriesTiedBonus"`
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
