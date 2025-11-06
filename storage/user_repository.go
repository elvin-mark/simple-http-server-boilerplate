package storage

import (
	"context"
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

// userRepository is the PostgreSQL implementation of the UserRepository.
type userRepository struct {
	db *pgxpool.Pool
}

// UserRepository creates a new UserRepository.
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

// GetUsers retrieves all users from the database.
func (r *userRepository) GetUsers() ([]user.User, error) {
	rows, err := r.db.Query(context.Background(), "SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// GetUser retrieves a single user by ID from the database.
func (r *userRepository) GetUser(id int) (*user.User, error) {
	var u user.User
	err := r.db.QueryRow(context.Background(), "SELECT id, name, email FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// CreateUser inserts a new user into the database.
func (r *userRepository) CreateUser(req *user.CreateUserRequest) (*user.User, error) {
	var u user.User
	err := r.db.QueryRow(context.Background(), "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email", req.Name, req.Email).Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// DeleteUser deletes a user from the database.
func (r *userRepository) DeleteUser(id int) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	return err
}
