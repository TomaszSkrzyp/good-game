package db

import (
	"context"
	"errors"

	"github.com/tomaszSkrzyp/good-game/models"
	"gorm.io/gorm"
)

type ConferenceNameToID map[string]uint

// exported package-level map populated at startup
var ConferenceNameToIDMap ConferenceNameToID

type ConferenceRepository struct {
	db *gorm.DB
}

func (r *ConferenceRepository) NameToIDMap() (ConferenceNameToID, error) {
	var confs []models.Conference
	if err := r.db.Find(&confs).Error; err != nil {
		return nil, err
	}
	m := make(ConferenceNameToID, len(confs))
	for _, c := range confs {
		m[c.Name] = c.ID
	}
	return m, nil
}

func NewConferenceRepository(db *gorm.DB) *ConferenceRepository {
	return &ConferenceRepository{db: db}
}

func (r *ConferenceRepository) Create(ctx context.Context, conference *models.Conference) error {
	return r.db.WithContext(ctx).Create(conference).Error
}

func (r *ConferenceRepository) GetByID(ctx context.Context, id uint) (*models.Conference, error) {
	var conference models.Conference
	if err := r.db.WithContext(ctx).First(&conference, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &conference, nil
}
func (r *ConferenceRepository) List(ctx context.Context) ([]models.Conference, error) {
	var conferences []models.Conference
	err := r.db.WithContext(ctx).Find(&conferences).Error
	return conferences, err
}
func (r *ConferenceRepository) Update(ctx context.Context, conference *models.Conference) error {
	return r.db.WithContext(ctx).Model(&conference).Updates(conference).Error
}

func (r *ConferenceRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&models.Conference{}, id)
	if res.RowsAffected == 0 {
		return errors.New("conference not found")
	}
	return res.Error
}
