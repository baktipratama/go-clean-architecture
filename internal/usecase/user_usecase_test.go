package usecase

import (
	"context"
	"testing"
	"time"

	"go-clean-code/internal/dto"
	"go-clean-code/internal/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepositoryInterface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *repository.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*repository.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*repository.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *repository.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*repository.User, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*repository.User), args.Error(1)
}

func TestUserUsecase_CreateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("should create user successfully", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		req := dto.CreateUserRequest{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockRepo.On("GetByEmail", ctx, req.Email).Return((*repository.User)(nil), repository.ErrUserNotFound)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*repository.User")).Return(nil)

		result, err := usecase.CreateUser(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Name, result.Name)
		assert.Equal(t, req.Email, result.Email)
		assert.NotEqual(t, uuid.Nil, result.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error for invalid input", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		req := dto.CreateUserRequest{
			Name:  "",
			Email: "john@example.com",
		}

		result, err := usecase.CreateUser(ctx, req)

		assert.Equal(t, ErrInvalidInput, err)
		assert.Nil(t, result)
	})

	t.Run("should return error when email exists", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		req := dto.CreateUserRequest{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		existingUser := &repository.User{
			ID:    uuid.New(),
			Name:  "Existing User",
			Email: req.Email,
		}

		mockRepo.On("GetByEmail", ctx, req.Email).Return(existingUser, nil)

		result, err := usecase.CreateUser(ctx, req)

		assert.Equal(t, ErrEmailExists, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_GetUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	t.Run("should get user successfully", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		user := &repository.User{
			ID:    userID,
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockRepo.On("GetByID", ctx, userID).Return(user, nil)

		result, err := usecase.GetUser(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, user.Name, result.Name)
		assert.Equal(t, user.Email, result.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		mockRepo.On("GetByID", ctx, userID).Return((*repository.User)(nil), repository.ErrUserNotFound)

		result, err := usecase.GetUser(ctx, userID)

		assert.Equal(t, repository.ErrUserNotFound, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_UpdateUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	t.Run("should update user name successfully", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		existingUser := &repository.User{
			ID:        userID,
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		req := dto.UpdateUserRequest{
			Name: "John Smith",
		}

		mockRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*repository.User")).Return(nil)

		result, err := usecase.UpdateUser(ctx, userID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Name, result.Name)
		assert.Equal(t, existingUser.Email, result.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should update user email successfully", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		existingUser := &repository.User{
			ID:        userID,
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		req := dto.UpdateUserRequest{
			Email: "johnsmith@example.com",
		}

		mockRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
		mockRepo.On("GetByEmail", ctx, req.Email).Return((*repository.User)(nil), repository.ErrUserNotFound)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*repository.User")).Return(nil)

		result, err := usecase.UpdateUser(ctx, userID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, existingUser.Name, result.Name)
		assert.Equal(t, req.Email, result.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when email already exists for different user", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		existingUser := &repository.User{
			ID:        userID,
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		otherUser := &repository.User{
			ID:    uuid.New(),
			Name:  "Other User",
			Email: "other@example.com",
		}

		req := dto.UpdateUserRequest{
			Email: "other@example.com",
		}

		mockRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
		mockRepo.On("GetByEmail", ctx, req.Email).Return(otherUser, nil)

		result, err := usecase.UpdateUser(ctx, userID, req)

		assert.Equal(t, ErrEmailExists, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_DeleteUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	t.Run("should delete user successfully", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		mockRepo.On("Delete", ctx, userID).Return(nil)

		err := usecase.DeleteUser(ctx, userID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		mockRepo.On("Delete", ctx, userID).Return(repository.ErrUserNotFound)

		err := usecase.DeleteUser(ctx, userID)

		assert.Equal(t, repository.ErrUserNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_ListUsers(t *testing.T) {
	ctx := context.Background()

	users := []*repository.User{
		{
			ID:    uuid.New(),
			Name:  "John Doe",
			Email: "john@example.com",
		},
		{
			ID:    uuid.New(),
			Name:  "Jane Smith",
			Email: "jane@example.com",
		},
	}

	t.Run("should list users successfully", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		limit, offset := 10, 0

		mockRepo.On("List", ctx, limit, offset).Return(users, nil)

		result, err := usecase.ListUsers(ctx, limit, offset)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Users, 2)
		assert.Equal(t, 2, result.Total)
		assert.Equal(t, limit, result.Limit)
		assert.Equal(t, offset, result.Offset)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should use default values for invalid pagination", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		mockRepo.On("List", ctx, 10, 0).Return([]*repository.User{}, nil)

		result, err := usecase.ListUsers(ctx, 0, -1)

		assert.NoError(t, err)
		assert.Equal(t, 10, result.Limit)
		assert.Equal(t, 0, result.Offset)
		mockRepo.AssertExpectations(t)
	})
}