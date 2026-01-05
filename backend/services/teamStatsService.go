package services

import (
	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/models"
)

type TeamStatsService struct {
	repo *db.TeamStatsRepository
}

func NewTeamStatsService(repo *db.TeamStatsRepository) *TeamStatsService {
	return &TeamStatsService{repo: repo}
}
func (s *TeamStatsService) GetByID(id uint) (*models.TeamStats, error) {
	return s.repo.GetByID(id)
}

func (s *TeamStatsService) Filter(teamID uint, season string) (*models.TeamStats, error) {
	return s.repo.Filter(teamID, season)
}
func (s *TeamStatsService) Create(teamStats *models.TeamStats) error {
	return s.repo.Create(teamStats)
}

func (s *TeamStatsService) Update(teamStats *models.TeamStats) error {
	return s.repo.Update(teamStats)
}
func (s *TeamStatsService) Delete(id uint) error {
	return s.repo.Delete(id)
}
