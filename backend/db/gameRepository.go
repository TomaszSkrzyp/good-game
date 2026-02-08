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
func (r *GameRepository) Filter(date string, homeID *uint, awayID *uint, minRating *int, maxRating *int, sort string, page int, limit int, userID uint) ([]models.Game, error) {
	var games []models.Game

	// base select
	selectQuery := "games.*, " +
		"(SELECT COALESCE(AVG(liked), 0) FROM user_reactions WHERE game_id = games.id) as avg_rating, " +
		"(SELECT COUNT(*) FROM user_reactions WHERE user_reactions.game_id = games.id) as rating_count"

	// dynamic select based on auth status
	var queryArgs []interface{}
	if userID > 0 {
		selectQuery += ", (SELECT liked FROM user_reactions WHERE game_id = games.id AND user_id = ? LIMIT 1) as rating"
		queryArgs = append(queryArgs, userID)
	} else {
		selectQuery += ", NULL as rating"
	}

	// initialize query with correct number of arguments
	query := r.db.Model(&models.Game{}).Select(selectQuery, queryArgs...)

	// filters
	if date != "" {
		query = query.Where("date(game_time) = ?", date)
	}

	// basic game filters
	if date != "" {
		query = query.Where("date(game_time) = ?", date)
	}
	if homeID != nil {
		query = query.Where("home_team_id = ?", *homeID)
	}
	if awayID != nil {
		query = query.Where("away_team_id = ?", *awayID)
	}

	// rating filters only for authenticated users
	if userID > 0 {
		if minRating != nil {
			query = query.Where("EXISTS (SELECT 1 FROM user_reactions WHERE game_id = games.id AND user_id = ? AND liked >= ?)", userID, *minRating)
		}
		if maxRating != nil {
			query = query.Where("EXISTS (SELECT 1 FROM user_reactions WHERE game_id = games.id AND user_id = ? AND liked <= ?)", userID, *maxRating)
		}
	}

	if sort != "" {
		query = query.Order(sort)
	}

	// fetch with preloading and pagination
	err := query.
		Preload("HomeTeam").
		Preload("AwayTeam").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&games).
		Error

	return games, err
}
