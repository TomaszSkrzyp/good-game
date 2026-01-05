package db

import (
	"database/sql"
	"errors"

	"github.com/tomaszSkrzyp/good-game/models"
	"gorm.io/gorm"
)

type UserReactionRepository struct {
	db *gorm.DB
}

func NewUserReactionRepository(db *gorm.DB) *UserReactionRepository {
	return &UserReactionRepository{db: db}
}

func (r *UserReactionRepository) Create(ur *models.UserReaction) error {
	return r.db.Create(ur).Error
}

func (r *UserReactionRepository) GetByID(id uint) (*models.UserReaction, error) {
	var ur models.UserReaction
	if err := r.db.Preload("User").Preload("Game").First(&ur, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ur, nil
}

func (r *UserReactionRepository) Update(ur *models.UserReaction) error {
	return r.db.Save(ur).Error
}

func (r *UserReactionRepository) Delete(id uint) error {
	res := r.db.Delete(&models.UserReaction{}, id)
	if res.RowsAffected == 0 {
		return errors.New("user reaction not found")
	}
	return res.Error
}

// Filter by optional userID, gameID, liked. zero values are ignored.
func (r *UserReactionRepository) Filter(userID, gameID uint, liked *int, page, limit int) ([]models.UserReaction, error) {
	var list []models.UserReaction
	query := r.db.Model(&models.UserReaction{}).Preload("User").Preload("Game")

	if userID != 0 {
		query = query.Where("user_id = ?", userID)
	}
	if gameID != 0 {
		query = query.Where("game_id = ?", gameID)
	}
	if liked != nil {
		query = query.Where("liked = ?", *liked)
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

func (r *UserReactionRepository) GetAverageForGame(gameID uint) (float64, error) {
	var avg sql.NullFloat64
	row := r.db.Model(&models.UserReaction{}).
		Where("game_id = ?", gameID).
		Select("AVG(liked)").Row()
	if err := row.Scan(&avg); err != nil {
		// if no rows, return 0
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	if !avg.Valid {
		return 0, nil
	}
	return avg.Float64, nil
}

func (r *UserReactionRepository) GetAverageForTeam(teamID uint) (float64, error) {
	var avg sql.NullFloat64
	row := r.db.Model(&models.UserReaction{}).
		Joins("JOIN games ON games.id = user_reactions.game_id").
		Where("games.home_team_id = ? OR games.away_team_id = ?", teamID, teamID).
		Select("AVG(liked)").Row()

	if err := row.Scan(&avg); err != nil {
		// no rows or other error
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	if !avg.Valid {
		return 0, nil
	}
	return avg.Float64, nil
}
