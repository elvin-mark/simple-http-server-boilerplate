package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	user "http-server/dto/user"
	"http-server/storage"
	"http-server/utils"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Initialize logger for tests
	utils.InitLogger("debug")
	// Run tests
	os.Exit(m.Run())
}

// MockUserRepository is a mock implementation of the UserRepository interface.
type MockUserRepository struct {
	GetUsersFunc   func() ([]user.User, error)
	GetUserFunc    func(id int) (*user.User, error)
	CreateUserFunc func(user *user.CreateUserRequest) (*user.User, error)
	DeleteUserFunc func(id int) error
}

func (m *MockUserRepository) GetUsers() ([]user.User, error) {
	if m.GetUsersFunc != nil {
		return m.GetUsersFunc()
	}
	return nil, errors.New("GetUsersFunc not implemented")
}

func (m *MockUserRepository) GetUser(id int) (*user.User, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(id)
	}
	return nil, errors.New("GetUserFunc not implemented")
}

func (m *MockUserRepository) CreateUser(req *user.CreateUserRequest) (*user.User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(req)
	}
	return nil, errors.New("CreateUserFunc not implemented")
}

func (m *MockUserRepository) DeleteUser(id int) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(id)
	}
	return errors.New("DeleteUserFunc not implemented")
}

func TestGetUsers(t *testing.T) {
	t.Run("should return users from cache when cache hit", func(t *testing.T) {
		// Arrange
		expectedUsers := []user.User{{ID: 1, Name: "Test User", Email: "test@example.com"}}
		usersJSON, _ := json.Marshal(expectedUsers)

		db, mock := redismock.NewClientMock()
		redisClient := &storage.RedisClient{Client: db}

		mock.ExpectGet("all_users").SetVal(string(usersJSON))

		repo := &MockUserRepository{}
		service := NewUserService(repo, redisClient)

		// Act
		users, err := service.GetUsers()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return users from db and set cache when cache miss", func(t *testing.T) {
		// Arrange
		expectedUsers := []user.User{{ID: 1, Name: "Test User", Email: "test@example.com"}}
		usersJSON, _ := json.Marshal(expectedUsers)

		db, mock := redismock.NewClientMock()
		redisClient := &storage.RedisClient{Client: db}

		mock.ExpectGet("all_users").RedisNil()
		mock.ExpectSet("all_users", usersJSON, 1*time.Minute).SetVal("OK")

		repo := &MockUserRepository{
			GetUsersFunc: func() ([]user.User, error) {
				return expectedUsers, nil
			},
		}
		service := NewUserService(repo, redisClient)

		// Act
		users, err := service.GetUsers()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when db fails and cache miss", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")

		db, mock := redismock.NewClientMock()
		redisClient := &storage.RedisClient{Client: db}

		mock.ExpectGet("all_users").RedisNil()

		repo := &MockUserRepository{
			GetUsersFunc: func() ([]user.User, error) {
				return nil, dbErr
			},
		}
		service := NewUserService(repo, redisClient)

		// Act
		users, err := service.GetUsers()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Equal(t, dbErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUser(t *testing.T) {
	userID := 1
	cacheKey := fmt.Sprintf("user:%d", userID)

	t.Run("should return user from cache when cache hit", func(t *testing.T) {
		// Arrange
		expectedUser := &user.User{ID: userID, Name: "Test User", Email: "test@example.com"}
		userJSON, _ := json.Marshal(expectedUser)

		db, mock := redismock.NewClientMock()
		redisClient := &storage.RedisClient{Client: db}

		mock.ExpectGet(cacheKey).SetVal(string(userJSON))

		repo := &MockUserRepository{}
		service := NewUserService(repo, redisClient)

		// Act
		u, err := service.GetUser(userID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, u)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return user from db and set cache when cache miss", func(t *testing.T) {
		// Arrange
		expectedUser := &user.User{ID: userID, Name: "Test User", Email: "test@example.com"}
		userJSON, _ := json.Marshal(expectedUser)

		db, mock := redismock.NewClientMock()
		redisClient := &storage.RedisClient{Client: db}

		mock.ExpectGet(cacheKey).RedisNil()
		mock.ExpectSet(cacheKey, userJSON, 1*time.Minute).SetVal("OK")

		repo := &MockUserRepository{
			GetUserFunc: func(id int) (*user.User, error) {
				assert.Equal(t, userID, id)
				return expectedUser, nil
			},
		}
		service := NewUserService(repo, redisClient)

		// Act
		u, err := service.GetUser(userID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, u)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when db fails and cache miss", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")

		db, mock := redismock.NewClientMock()
		redisClient := &storage.RedisClient{Client: db}

		mock.ExpectGet(cacheKey).RedisNil()

		repo := &MockUserRepository{
			GetUserFunc: func(id int) (*user.User, error) {
				return nil, dbErr
			},
		}
		service := NewUserService(repo, redisClient)

		// Act
		u, err := service.GetUser(userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, u)
		assert.Equal(t, dbErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreateUser(t *testing.T) {
	createUserReq := &user.CreateUserRequest{Name: "New User", Email: "new@example.com"}
	createdUser := &user.User{ID: 2, Name: "New User", Email: "new@example.com"}

	t.Run("should create user and invalidate cache", func(t *testing.T) {
		// Arrange
		db, mock := redismock.NewClientMock()
		redisClient := &storage.RedisClient{Client: db}

		mock.ExpectDel("all_users").SetVal(1)
		mock.ExpectDel(fmt.Sprintf("user:%d", createdUser.ID)).SetVal(1)

		repo := &MockUserRepository{
			CreateUserFunc: func(req *user.CreateUserRequest) (*user.User, error) {
				assert.Equal(t, createUserReq, req)
				return createdUser, nil
			},
		}
		service := NewUserService(repo, redisClient)

		// Act
		u, err := service.CreateUser(createUserReq)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, createdUser, u)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when db fails", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")

		db, mock := redismock.NewClientMock()
		redisClient := &storage.RedisClient{Client: db}

		repo := &MockUserRepository{
			CreateUserFunc: func(req *user.CreateUserRequest) (*user.User, error) {
				return nil, dbErr
			},
		}
		service := NewUserService(repo, redisClient)

		// Act
		u, err := service.CreateUser(createUserReq)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, u)
		assert.Equal(t, dbErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteUser(t *testing.T) {
	userID := 1

	t.Run("should delete user and invalidate cache", func(t *testing.T) {
		// Arrange
		db, mock := redismock.NewClientMock()
		redisClient := &storage.RedisClient{Client: db}

		mock.ExpectDel("all_users").SetVal(1)
		mock.ExpectDel(fmt.Sprintf("user:%d", userID)).SetVal(1)

		repo := &MockUserRepository{
			DeleteUserFunc: func(id int) error {
				assert.Equal(t, userID, id)
				return nil
			},
		}
		service := NewUserService(repo, redisClient)

		// Act
		err := service.DeleteUser(userID)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when db fails", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")

		db, mock := redismock.NewClientMock()
		redisClient := &storage.RedisClient{Client: db}

		repo := &MockUserRepository{
			DeleteUserFunc: func(id int) error {
				return dbErr
			},
		}
		service := NewUserService(repo, redisClient)

		// Act
		err := service.DeleteUser(userID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
