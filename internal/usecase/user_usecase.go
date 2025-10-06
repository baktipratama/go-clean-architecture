package usecase

import (
	"context"
	"errors"
	"time"

	"go-clean-code/internal/dto"
	"go-clean-code/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrEmailExists  = errors.New("email already exists")
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
	if req.Name == "" || req.Email == "" {
		return nil, ErrInvalidInput
	}

	existingUser, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailExists
	}

	now := time.Now()
	user := &repository.User{
		ID:        uuid.New(),
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: now,
		UpdatedAt: now,
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

	if req.Name != "" {
		user.Name = req.Name
		user.UpdatedAt = time.Now()
	}

	if req.Email != "" {
		existingUser, err := u.userRepo.GetByEmail(ctx, req.Email)
		if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
			return nil, err
		}
		if existingUser != nil && existingUser.ID != id {
			return nil, ErrEmailExists
		}
		user.Email = req.Email
		user.UpdatedAt = time.Now()
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