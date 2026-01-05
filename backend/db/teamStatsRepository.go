package db

import (
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

func (r *TeamStatsRepository) Create(teamStats *models.TeamStats) error {
	return r.db.Create(teamStats).Error
}
func (r *TeamStatsRepository) GetByID(id uint) (*models.TeamStats, error) {
	var teamStats models.TeamStats
	if err := r.db.Preload("Team.Conference").First(&teamStats, id).Error; err != nil {
		return nil, err
	}
	return &teamStats, nil
}

func (r *TeamStatsRepository) Update(teamStats *models.TeamStats) error {
	return r.db.Save(teamStats).Error
}

func (r *TeamStatsRepository) Delete(id uint) error {
	res := r.db.Delete(&models.TeamStats{}, id)
	if res.RowsAffected == 0 {
		return errors.New("team stats not found")
	}
	return res.Error
}
func (r *TeamStatsRepository) Filter(teamID uint, season string) (*models.TeamStats, error) {
	var teamStats models.TeamStats
	if err := r.db.Where("team_id = ? AND season = ?", teamID, season).First(&teamStats).Error; err != nil {
		return nil, err
	}
	return &teamStats, nil
	
}