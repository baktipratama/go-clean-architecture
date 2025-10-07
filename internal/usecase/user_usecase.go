package usecase

import (
	"context"
	"errors"

	"go-clean-code/internal/domain"
	"go-clean-code/internal/dto"
	"go-clean-code/internal/repository"

	"github.com/google/uuid"
)

// Usecase level errors for backward compatibility with tests
var (
	ErrInvalidInput = errors.New("invalid input")
	ErrEmailExists  = errors.New("email already exists")
	ErrUserNotFound = errors.New("user not found")
)



type UserUsecaseInterface interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUser(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context, limit, offset int) (*dto.ListUsersResponse, error)
}

type UserUsecase struct {
	userRepo repository.UserRepositoryInterface
}

func NewUserUsecase(userRepo repository.UserRepositoryInterface) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Use domain entity to create user with validation
	user, err := domain.NewUser(req.Name, req.Email)
	if err != nil {
		return nil, domain.NewValidationError("invalid user input", err)
	}

	// Check if email already exists
	existingUser, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil && !domain.IsNotFoundError(err) {
		return nil, err
	}
	if existingUser != nil {
		return nil, domain.NewConflictError("email already in use", domain.ErrEmailAlreadyUsed)
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (u *UserUsecase) GetUser(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (u *UserUsecase) UpdateUser(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Use domain entity methods for validation and updates
	if req.Name != "" {
		if err := user.UpdateName(req.Name); err != nil {
			return nil, domain.NewValidationError("invalid name", err)
		}
	}

	if req.Email != "" {
		// Check if email already exists for another user
		existingUser, err := u.userRepo.GetByEmail(ctx, req.Email)
		if err != nil && !domain.IsNotFoundError(err) {
			return nil, err
		}
		if existingUser != nil && existingUser.ID != id {
			return nil, domain.NewConflictError("email already in use by another user", domain.ErrEmailAlreadyUsed)
		}
		
		if err := user.UpdateEmail(req.Email); err != nil {
			return nil, domain.NewValidationError("invalid email", err)
		}
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return u.userRepo.Delete(ctx, id)
}

func (u *UserUsecase) ListUsers(ctx context.Context, limit, offset int) (*dto.ListUsersResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	users, err := u.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	userResponses := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = &dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
	}

	return &dto.ListUsersResponse{
		Users:  userResponses,
		Total:  len(userResponses),
		Limit:  limit,
		Offset: offset,
	}, nil
}