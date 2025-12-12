package services

import (
	user "http-server/dto/user"
	"http-server/storage"
)

type UserService interface {
	GetUsers() ([]user.User, error)
	GetUser(id int) (*user.User, error)
	CreateUser(req *user.CreateUserRequest) (*user.User, error)
	DeleteUser(id int) error
}

// NewUserService creates a new UserService.
func NewUserService(repo storage.UserRepository, redisClient *storage.RedisClient) UserService {
	return &userServiceImpl{repo: repo, redisClient: redisClient}
}
