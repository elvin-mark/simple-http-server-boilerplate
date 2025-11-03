package services

import (
	user "http-server/dto/user"
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
	repo UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetUsers returns all users.
func (s *UserService) GetUsers() ([]user.User, error) {
	return s.repo.GetUsers()
}

// GetUser returns a user by ID.
func (s *UserService) GetUser(id int) (*user.User, error) {
	return s.repo.GetUser(id)
}

// CreateUser creates a new user.
func (s *UserService) CreateUser(req *user.CreateUserRequest) (*user.User, error) {
	// Here you could add more complex business logic, e.g.,
	// - check if email is unique
	// - hash password if there was one
	// - send a welcome email
	return s.repo.CreateUser(req)
}

// DeleteUser deletes a user by ID.
func (s *UserService) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}
