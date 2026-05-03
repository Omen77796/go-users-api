package services

import (
	"context"
	"encoding/json"
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/omen77796/go-users-api/internal/logger"
	"github.com/omen77796/go-users-api/internal/models"
	"github.com/omen77796/go-users-api/internal/repository"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type UserService struct {
	repo *repository.UserRepository
	rdb  *redis.Client
}

func NewUserService(r *repository.UserRepository, rdb *redis.Client) *UserService {
	return &UserService{
		repo: r,
		rdb:  rdb,
	}
}

// Create

func (s *UserService) Create(user *models.User) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	if user.Name == "" {
		return errors.New("name is required")
	}

	if user.Email == "" {
		return errors.New("email is required")
	}

	if _, err := mail.ParseAddress(user.Email); err != nil {
		return errors.New("invalid email format")
	}

	err := s.repo.Create(user)
	if err != nil {
		return err
	}

	// 🔥 invalidar cache
	s.rdb.Del(context.Background(), "users")

	return nil
}

// GET ALL con cache

func (s *UserService) GetAll() ([]models.User, error) {
	ctx := context.Background()

	// 🔥 1. Intentar cache
	cached, err := s.rdb.Get(ctx, "users").Result()
	if err == nil {
		var users []models.User
		if err := json.Unmarshal([]byte(cached), &users); err == nil {
			logger.Log.Info("cache hit: users")
			return users, nil
		}
		logger.Log.Warn("failed to unmarshal cache, fallback to DB", zap.Error(err))
	}

	// 🔥 2. DB fallback
	logger.Log.Info("cache miss: querying database")
	users, err := s.repo.GetAll()
	if err != nil {
		logger.Log.Error("failed to get users from repository", zap.Error(err))
		return nil, err
	}

	// 🔥 3. Guardar en cache
	jsonData, err := json.Marshal(users)
	if err == nil {
		s.rdb.Set(ctx, "users", jsonData, time.Minute)
	} else {
		logger.Log.Warn("failed to marshal users for cache", zap.Error(err))
	}

	return users, nil
}

// GET by ID

func (s *UserService) GetByID(id int) (*models.User, error) {
	return s.repo.GetByID(id)
}

// DELETE

func (s *UserService) Delete(id int) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	s.rdb.Del(context.Background(), "users")

	return nil
}
