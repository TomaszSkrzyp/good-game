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
