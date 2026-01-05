package services

import (
	"fmt"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/models"
	"gorm.io/gorm"
)

type GameService struct {
	repo *db.GameRepository
}

func NewGameService(repo *db.GameRepository) *GameService {
	return &GameService{repo: repo}
}
func (s *GameService) Create(game *models.Game) error {
	existingGame, err := s.repo.GetByID(game.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if existingGame != nil {
		return fmt.Errorf("game already exists")
	}
	return s.repo.Create(game)
}

func (s *GameService) GetByID(id uint) (*models.Game, error) {
	return s.repo.GetByID(id)
}

func (s *GameService) Update(game *models.Game) error {
	existingGame, err := s.repo.GetByID(game.ID)
	if err != nil {
		return err
	}
	if existingGame == nil {
		return fmt.Errorf("game does not exist")
	}
	return s.repo.Update(game)
}

func (s *GameService) Delete(id uint) error {
	return s.repo.Delete(id)
}
func (s *GameService) Filter(date string, homeID *uint, awayID *uint, minRating *int, maxRating *int, sort string, page int, limit int) ([]models.Game, error) {
	return s.repo.Filter(date, homeID, awayID, minRating, maxRating, sort, page, limit)
}
