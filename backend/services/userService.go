package services

import (
	"fmt"
	"time"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo *db.UserRepository
}

func NewUserService(repo *db.UserRepository) *UserService {
	return &UserService{repo: repo}
}
func (s *UserService) Register(username, password string, roleID uint) (*models.User, error) {

	existingUser, err := s.repo.GetByUserName(username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("account already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := &models.User{
		UserName: username,
		Password: string(hashedPassword),
		RoleID:   roleID,
	}

	if err := s.repo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *UserService) Authenticate(username, password string) (*models.User, error) {
	existingUser, err := s.repo.GetByUserName(username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existingUser == nil {
		return nil, fmt.Errorf("account doesn't exists")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("credentials incorrect")
	}
	existingUser.LastLoginAt = time.Now()
	if err := s.repo.Update(existingUser); err != nil {
		return nil, fmt.Errorf("failed to update last login: %w", err)
	}
	existingUser.Password = ""
	return existingUser, nil

}
