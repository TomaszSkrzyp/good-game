package services

import (
	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/models"
)

type UserReactionService struct {
	repo *db.UserReactionRepository
}

func NewUserReactionService(r *db.UserReactionRepository) *UserReactionService {
	return &UserReactionService{repo: r}
}

func (s *UserReactionService) Create(ur *models.UserReaction) error {
	return s.repo.Create(ur)
}
func (s *UserReactionService) GetByID(id uint) (*models.UserReaction, error) {
	return s.repo.GetByID(id)
}
func (s *UserReactionService) Update(ur *models.UserReaction) error {
	return s.repo.Update(ur)
}
func (s *UserReactionService) Delete(id uint) error {
	return s.repo.Delete(id)
}
func (s *UserReactionService) Filter(userID, gameID uint, liked *int, page, limit int) ([]models.UserReaction, error) {
	return s.repo.Filter(userID, gameID, liked, page, limit)
}
func (s *UserReactionService) GetAverageForGame(gameID uint) (float64, error) {
	return s.repo.GetAverageForGame(gameID)
}
func (s *UserReactionService) GetAverageForTeam(teamID uint) (float64, error) {
	return s.repo.GetAverageForTeam(teamID)
}
