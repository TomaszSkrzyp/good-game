package db

import (
	"context"
	"errors"

	"github.com/tomaszSkrzyp/good-game/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserReactionRepository struct {
	db *gorm.DB
}

func NewUserReactionRepository(db *gorm.DB) *UserReactionRepository {
	return &UserReactionRepository{db: db}
}

func (r *UserReactionRepository) UpdateOrCreate(ctx context.Context, ur *models.UserReaction) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "game_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"rating"}),
	}).Create(ur).Error
}

func (r *UserReactionRepository) GetByID(ctx context.Context, id uint) (*models.UserReaction, error) {
	var ur models.UserReaction
	if err := r.db.WithContext(ctx).Preload("User").Preload("Game").First(&ur, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ur, nil
}

func (r *UserReactionRepository) Update(ctx context.Context, ur *models.UserReaction) error {
	return r.db.WithContext(ctx).Save(ur).Error
}

func (r *UserReactionRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&models.UserReaction{}, id)
	if res.RowsAffected == 0 {
		return errors.New("user reaction not found")
	}
	return res.Error
}

func (r *UserReactionRepository) Filter(ctx context.Context, userID, gameID uint, rating *int, page, limit int) ([]models.UserReaction, error) {
	var list []models.UserReaction
	query := r.db.WithContext(ctx).Model(&models.UserReaction{}).
		Preload("User").
		Preload("Game.HomeTeam").
		Preload("Game.AwayTeam")

	if userID != 0 {
		query = query.Where("user_id = ?", userID)
	}
	if gameID != 0 {
		query = query.Where("game_id = ?", gameID)
	}
	if rating != nil {
		query = query.Where("rating = ?", *rating)
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 50
	}
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

type ReactionStats struct {
	Average float64
	Count   int64
}

func (r *UserReactionRepository) GetStatsForGame(ctx context.Context, gameID uint) (float64, int64, error) {
	var stats ReactionStats
	err := r.db.WithContext(ctx).Model(&models.UserReaction{}).
		Select("COALESCE(AVG(rating), 0) as average, COUNT(*) as count").
		Where("game_id = ?", gameID).
		Scan(&stats).Error

	return stats.Average, stats.Count, err
}

func (r *UserReactionRepository) GetStatsForTeam(ctx context.Context, teamID uint) (float64, int64, error) {
	var stats ReactionStats
	err := r.db.WithContext(ctx).Table("user_reactions").
		Select("COALESCE(AVG(rating), 0) as average, COUNT(user_reactions.id) as count").
		Joins("JOIN games ON games.id = user_reactions.game_id").
		Where("games.home_team_id = ? OR games.away_team_id = ?", teamID, teamID).
		Scan(&stats).Error

	return stats.Average, stats.Count, err
}
