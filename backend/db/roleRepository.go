package db

import (
	"errors"

	"github.com/tomaszSkrzyp/good-game/models"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) GetByID(id uint) (*models.Role, error) {
	var role models.Role
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

func (r *RoleRepository) Delete(id uint) error {
	res := r.db.Delete(&models.Role{}, id)
	if res.RowsAffected == 0 {
		return errors.New("role not found")
	}
	return res.Error
}
