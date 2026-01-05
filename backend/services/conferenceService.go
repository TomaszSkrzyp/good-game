package services

import (
	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/models"
)

type ConferenceService struct {
	repo *db.ConferenceRepository
}

func NewConferenceService(repo *db.ConferenceRepository) *ConferenceService {
	return &ConferenceService{repo: repo}
}

func (s *ConferenceService) GetByID(id uint) (*models.Conference, error) {
	return s.repo.GetByID(id)
}
func (s *ConferenceService) List() ([]models.Conference, error) {
	return s.repo.List()
}
