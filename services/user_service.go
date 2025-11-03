package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	user "http-server/dto/user"
	"http-server/storage"
	"http-server/utils"
)

// UserRepository defines the interface for user data storage.
// This will be implemented by the storage layer (e.g., in-memory or PostgreSQL).
type UserRepository interface {
	GetUsers() ([]user.User, error)
	GetUser(id int) (*user.User, error)
	CreateUser(user *user.CreateUserRequest) (*user.User, error)
	DeleteUser(id int) error
}

// UserService provides user-related business logic.
type UserService struct {
	repo        UserRepository
	redisClient *storage.RedisClient
}

// NewUserService creates a new UserService.
func NewUserService(repo UserRepository, redisClient *storage.RedisClient) *UserService {
	return &UserService{repo: repo, redisClient: redisClient}
}

// GetUsers returns all users.
func (s *UserService) GetUsers() ([]user.User, error) {
	ctx := context.Background()
	cacheKey := "all_users"

	// Try to get from cache
	val, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var users []user.User
		if err := json.Unmarshal([]byte(val), &users); err == nil {
			utils.Logger.Info("Cache hit for all users")
			return users, nil
		}
	}

	// Get from DB
	users, err := s.repo.GetUsers()
	if err != nil {
		return nil, err
	}

	// Store in cache
	data, err := json.Marshal(users)
	if err == nil {
		s.redisClient.Set(ctx, cacheKey, data, 1*time.Minute) // Cache for 1 minute
	}

	return users, nil
}

// GetUser returns a user by ID.
func (s *UserService) GetUser(id int) (*user.User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%d", id)

	// Try to get from cache
	val, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var u user.User
		if err := json.Unmarshal([]byte(val), &u); err == nil {
			utils.Logger.Info("Cache hit for user", "id", id)
			return &u, nil
		}
	}

	// Get from DB
	u, err := s.repo.GetUser(id)
	if err != nil {
		return nil, err
	}

	// Store in cache
	data, err := json.Marshal(u)
	if err == nil {
		s.redisClient.Set(ctx, cacheKey, data, 1*time.Minute) // Cache for 1 minute
	}

	return u, nil
}

// CreateUser creates a new user.
func (s *UserService) CreateUser(req *user.CreateUserRequest) (*user.User, error) {
	createdUser, err := s.repo.CreateUser(req)
	if err != nil {
		return nil, err
	}

	// Invalidate cache for all users and the specific user
	ctx := context.Background()
	s.redisClient.Del(ctx, "all_users")
	s.redisClient.Del(ctx, fmt.Sprintf("user:%d", createdUser.ID))

	return createdUser, nil
}

// DeleteUser deletes a user by ID.
func (s *UserService) DeleteUser(id int) error {
	if err := s.repo.DeleteUser(id); err != nil {
		return err
	}

	// Invalidate cache for all users and the specific user
	ctx := context.Background()
	s.redisClient.Del(ctx, "all_users")
	s.redisClient.Del(ctx, fmt.Sprintf("user:%d", id))

	return nil
}
