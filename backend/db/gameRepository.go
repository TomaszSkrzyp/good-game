package db

import (
	"errors"

	"github.com/tomaszSkrzyp/good-game/models"
	"gorm.io/gorm"
)

type GameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) *GameRepository {
	return &GameRepository{db: db}
}

func (r *GameRepository) Create(game *models.Game) error {
	return r.db.Create(game).Error
}

func (r *GameRepository) GetByID(id uint) (*models.Game, error) {
	var game models.Game
	if err := r.db.First(&game, id).Error; err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *GameRepository) Update(game *models.Game) error {
	return r.db.Save(game).Error
}

func (r *GameRepository) Delete(id uint) error {
	res := r.db.Delete(&models.Game{}, id)
	if res.RowsAffected == 0 {
		return errors.New("game not found")
	}
	return res.Error
}

func (r *GameRepository) Filter(date string, homeID *uint, awayID *uint, minRating *int, maxRating *int, sort string, page int, limit int) ([]models.Game, error) {
	var games []models.Game
	query := r.db.Model(&models.Game{})

	if date != "" {
		query = query.Where("date = ?", date)
	}
	if homeID != nil {
		query = query.Where("home_team_id = ?", *homeID)
	}
	if awayID != nil {
		query = query.Where("away_team_id = ?", *awayID)
	}
	if minRating != nil {
		query = query.Where("rating >= ?", *minRating)
	}
	if maxRating != nil {
		query = query.Where("rating <= ?", *maxRating)
	}

	if sort != "" {
		query = query.Order(sort)
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&games).Error; err != nil {
		return nil, err
	}
	return games, nil
}
