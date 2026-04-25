package db

import (
	"context"
	"errors"

	"github.com/tomaszSkrzyp/good-game/models"
	"gorm.io/gorm"
)

type TeamNameToID map[string]uint

// exported package-level map populated at startup
var TeamNameToIDMap TeamNameToID

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(ctx context.Context, team *models.Team) error {
	return r.db.WithContext(ctx).Create(team).Error
}

func (r *TeamRepository) GetByID(ctx context.Context, id uint) (*models.Team, error) {
	var team models.Team
	if err := r.db.WithContext(ctx).Preload("Conference").First(&team, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &team, nil
}
func (r *TeamRepository) Update(ctx context.Context, team *models.Team) error {
	return r.db.WithContext(ctx).Save(team).Error
}

func (r *TeamRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&models.Team{}, id)
	if res.RowsAffected == 0 {
		return errors.New("team not found")
	}
	return res.Error
}
func (r *TeamRepository) Filter(ctx context.Context, name, conferenceName string, conferenceID uint) ([]models.Team, error) {
	var teams []models.Team
	query := r.db.WithContext(ctx).Model(&models.Team{}).Preload("Conference")

	if name != "" {
		query = query.Where("teams.name = ?", name)
	}
	if conferenceID != 0 {
		query = query.Where("teams.conference_id = ?", conferenceID)
	}
	if conferenceName != "" {
		query = query.Joins("JOIN conferences ON conferences.id = teams.conference_id").
			Where("conferences.name = ?", conferenceName)
	}

	err := query.Find(&teams).Error
	return teams, err
}
