package db

import (
	"context"
	"errors"
	"log"

	"github.com/tomaszSkrzyp/good-game/models"
	"gorm.io/gorm"
)

type ConfigRepository struct {
	db *gorm.DB
}

func NewConfigRepository(db *gorm.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}
func (r *ConfigRepository) SaveGameQuality(ctx context.Context, cfg models.GameQualityConfig) error {
	record := models.ConfigRecord{
		Key:  "game_quality",
		Data: cfg,
	}
	return r.db.WithContext(ctx).
		Where(models.ConfigRecord{Key: "game_quality"}).
		Assign(models.ConfigRecord{Data: cfg}).
		FirstOrCreate(&record).Error
}

func (r *ConfigRepository) GetGameQuality(ctx context.Context) (*models.GameQualityConfig, error) {
	var record models.ConfigRecord
	err := r.db.WithContext(ctx).Where("key = ?", "game_quality").First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record.Data, nil
}

func SyncAlgorithmConfig(ctx context.Context, db *gorm.DB) error {
	var record models.ConfigRecord

	err := db.WithContext(ctx).Where("key = ?", "game_quality").First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Config not found in DB. Seeding defaults...")

			record = models.ConfigRecord{
				Key:  "game_quality",
				Data: defaultGameConfig,
			}

			if err := db.Create(&record).Error; err != nil {
				return err
			}

			models.SetConfig(defaultGameConfig)
			return nil
		}
		return err
	}
	models.SetConfig(record.Data)
	log.Println("Global config hydrated from Database.")
	return nil
}

var defaultGameConfig = models.GameQualityConfig{
	Margins: []models.MarginWeight{
		{MaxMargin: 2, Points: 30},
		{MaxMargin: 5, Points: 20},
		{MaxMargin: 10, Points: 10},
	},

	HugeSwingBonus: 15,
	ClutchBonus:    10,
	OvertimeBonus:  15,
	ShootoutBonus:  10,
	GrittyBonus:    5,

	ShootoutThreshold:    240,
	GrittyThreshold:      185,
	StarDuelBonus:        10,
	StarPointsThreshold:  40,
	BigGameBonus:         15,
	BigScoringThreshold:  55,
	VersatilityThreshold: 3,

	EliminationBonus:      15,
	Game7Bonus:            30,
	PlayoffBonus:          10,
	SeasonSeriesTiedBonus: 5,

	VolatilityWeight:  3.5,
	SwingWeight:       15.0,
	ComebackThreshold: 0.85,
	MaxScore:          100.0,
}
