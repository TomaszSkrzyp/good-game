package services

import (
	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/models"
)

type TeamService struct {
	repo *db.TeamRepository
}

func NewTeamService(repo *db.TeamRepository) *TeamService {
	return &TeamService{repo: repo}
}
func (s *TeamService) GetByID(id uint) (*models.Team, error) {
	return s.repo.GetByID(id)
}

func (s *TeamService) Filter(name, conferenceName string, conferenceID uint) ([]models.Team, error) {
	return s.repo.Filter(name, conferenceName, conferenceID)
}
