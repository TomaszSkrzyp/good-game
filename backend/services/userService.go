package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tomaszSkrzyp/good-game/db"
	"github.com/tomaszSkrzyp/good-game/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrEmailExists  = errors.New("email already registered")
	ErrInvalidAuth  = errors.New("invalid username or password")
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	repo *db.UserRepository
}

func NewUserService(repo *db.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, username, password, email string, roleID uint) (*models.User, error) {
	exists, _ := s.repo.GetByUserName(ctx, username)
	if exists != nil {
		return nil, ErrUserExists
	}

	exists, _ = s.repo.GetByEmail(ctx, email)
	if exists != nil {
		return nil, ErrEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing error: %w", err)
	}

	newUser := &models.User{
		UserName: username,
		Password: string(hashedPassword),
		Email:    email,
		RoleID:   roleID,
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *UserService) Authenticate(ctx context.Context, username, password string) (*models.User, error) {
	user, err := s.repo.GetByUserName(ctx, username)
	if err != nil || user == nil {
		return nil, ErrInvalidAuth
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidAuth
	}

	user.LastLoginAt = time.Now()
	_ = s.repo.Update(ctx, user)

	user.Password = ""
	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) SetHideScores(ctx context.Context, userID uint, hideScores bool) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return ErrUserNotFound
	}

	user.HideScores = hideScores
	return s.repo.Update(ctx, user)
}
