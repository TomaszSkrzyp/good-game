package db

import (
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

func (r *ConferenceRepository) Create(conference *models.Conference) error {
	return r.db.Create(conference).Error
}

func (r *ConferenceRepository) GetByID(id uint) (*models.Conference, error) {
	var conference models.Conference
	if err := r.db.First(&conference, id).Error; err != nil {
		return nil, err
	}
	return &conference, nil
}
func (r *ConferenceRepository) List() ([]models.Conference, error) {
	var conferences []models.Conference
	err := r.db.Find(&conferences).Error
	return conferences, err
}
func (r *ConferenceRepository) Update(conference *models.Conference) error {
	return r.db.Save(conference).Error
}

func (r *ConferenceRepository) Delete(id uint) error {
	res := r.db.Delete(&models.Conference{}, id)
	if res.RowsAffected == 0 {
		return errors.New("conference not found")
	}
	return res.Error
}
