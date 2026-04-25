package db

import (
	"context"
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

func (r *GameRepository) Create(ctx context.Context, game *models.Game) error {
	return r.db.WithContext(ctx).Create(game).Error
}

func (r *GameRepository) GetByID(ctx context.Context, id uint, userID uint) (*models.Game, error) {
	var game models.Game

	err := r.db.WithContext(ctx).Model(&models.Game{}).
		Select("games.*, "+
			"COALESCE(AVG(all_reactions.rating), 0) as avg_rating, "+
			"COUNT(all_reactions.id) as rating_count, "+
			"(SELECT rating FROM user_reactions WHERE game_id = games.id AND user_id = ?) as rating", userID). // FIXED HERE
		Joins("LEFT JOIN user_reactions AS all_reactions ON all_reactions.game_id = games.id").
		Group("games.id").
		Preload("HomeTeam").
		Preload("AwayTeam").
		First(&game, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &game, nil
}

func (r *GameRepository) Update(ctx context.Context, game *models.Game) error {
	return r.db.WithContext(ctx).Save(game).Error
}

func (r *GameRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&models.Game{}, id)
	if res.RowsAffected == 0 {
		return errors.New("game not found")
	}
	return res.Error
}

func (r *GameRepository) Filter(ctx context.Context, date string, homeID, awayID *uint, minRating, maxRating *int, sort string, page, limit int, userID uint) ([]models.Game, error) {
	var games []models.Game
	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).Model(&models.Game{}).
		Select("games.*, "+
			"COALESCE(AVG(all_reactions.rating), 0) as avg_rating, "+
			"COUNT(all_reactions.id) as rating_count, "+
			"(SELECT rating FROM user_reactions WHERE game_id = games.id AND user_id = ?) as rating", userID).
		Joins("LEFT JOIN user_reactions AS all_reactions ON all_reactions.game_id = games.id").
		Group("games.id").
		Preload("HomeTeam").
		Preload("AwayTeam")
	if date != "" {
		query = query.Where("DATE(game_time AT TIME ZONE 'America/New_York') = ?", date)
	}
	if homeID != nil {
		query = query.Where("home_team_id = ?", *homeID)
	}
	if awayID != nil {
		query = query.Where("away_team_id = ?", *awayID)
	}

	if minRating != nil {
		query = query.Having("AVG(all_reactions.rating) >= ?", *minRating)
	}
	if maxRating != nil {
		query = query.Having("AVG(all_reactions.rating) <= ?", *maxRating)
	}

	if sort != "" {
		query = query.Order(sort)
	} else {
		query = query.Order("game_time ASC")
	}

	err := query.Offset(offset).Limit(limit).Find(&games).Error
	return games, err
}
