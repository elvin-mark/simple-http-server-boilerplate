package storage

import (
	user "http-server/dto/user"

	"github.com/jackc/pgx/v4/pgxpool"
)

// UserRepository defines the interface for user data storage.
type UserRepository interface {
	GetUsers() ([]user.User, error)
	GetUser(id int) (*user.User, error)
	CreateUser(user *user.CreateUserRequest) (*user.User, error)
	DeleteUser(id int) error
}

// UserRepository creates a new UserRepository.
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepositoryImpl{db: db}
}
