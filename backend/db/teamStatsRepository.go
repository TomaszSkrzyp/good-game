package db

import (
	"context"
	"errors"

	"github.com/tomaszSkrzyp/good-game/models"
	"gorm.io/gorm"
)

type TeamStatsRepository struct {
	db *gorm.DB
}

func NewTeamStatsRepository(db *gorm.DB) *TeamStatsRepository {
	return &TeamStatsRepository{db: db}
}

func (r *TeamStatsRepository) Create(ctx context.Context, teamStats *models.TeamStats) error {
	return r.db.WithContext(ctx).Create(teamStats).Error
}
func (r *TeamStatsRepository) GetByID(ctx context.Context, id uint) (*models.TeamStats, error) {
	var teamStats models.TeamStats
	if err := r.db.WithContext(ctx).Preload("Team.Conference").First(&teamStats, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &teamStats, nil
}

func (r *TeamStatsRepository) Update(ctx context.Context, teamStats *models.TeamStats) error {
	return r.db.WithContext(ctx).Model(&teamStats).Updates(teamStats).Error
}

func (r *TeamStatsRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&models.TeamStats{}, id)
	if res.RowsAffected == 0 {
		return errors.New("team stats not found")
	}
	return res.Error
}
func (r *TeamStatsRepository) Filter(ctx context.Context, teamID uint, season string) (*models.TeamStats, error) {
	var teamStats models.TeamStats
	if err := r.db.WithContext(ctx).Where("team_id = ? AND season = ?", teamID, season).First(&teamStats).Error; err != nil {
		return nil, err
	}
	return &teamStats, nil

}
