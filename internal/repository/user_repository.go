package repository

import (
	"context"
	"database/sql"

	"go-clean-code/internal/domain"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)
}

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, name, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		if isUniqueConstraintError(err) {
			return domain.NewConflictError("user already exists", domain.ErrUserAlreadyExists)
		}
		return domain.NewInternalError("failed to create user", err)
	}

	return nil
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE id = $1`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NewNotFoundError("user not found", domain.ErrUserNotFound)
		}
		return nil, domain.NewInternalError("failed to get user by ID", err)
	}

	return user, nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE email = $1`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NewNotFoundError("user not found by email", domain.ErrUserNotFound)
		}
		return nil, domain.NewInternalError("failed to get user by email", err)
	}

	return user, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET name = $2, email = $3, updated_at = $4
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.UpdatedAt)
	if err != nil {
		if isUniqueConstraintError(err) {
			return domain.NewConflictError("email already in use", domain.ErrEmailAlreadyUsed)
		}
		return domain.NewInternalError("failed to update user", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.NewInternalError("failed to get rows affected", err)
	}

	if rowsAffected == 0 {
		return domain.NewNotFoundError("user not found for update", domain.ErrUserNotFound)
	}

	return nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return domain.NewInternalError("failed to delete user", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.NewInternalError("failed to get rows affected", err)
	}

	if rowsAffected == 0 {
		return domain.NewNotFoundError("user not found for deletion", domain.ErrUserNotFound)
	}

	return nil
}

func (r *UserRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, domain.NewInternalError("failed to list users", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, domain.NewInternalError("failed to scan user", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, domain.NewInternalError("error iterating rows", err)
	}

	return users, nil
}

func isUniqueConstraintError(err error) bool {
	// PostgreSQL unique constraint error check
	// This is a simplified check - in production you might want to use pq.Error
	return err != nil && (
		// Check for unique constraint violation (email)
		contains(err.Error(), "duplicate key value") ||
		contains(err.Error(), "unique constraint"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		 containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
