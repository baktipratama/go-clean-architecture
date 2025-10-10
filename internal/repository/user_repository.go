package repository

import (
	"context"
	"database/sql"

	"go-clean-code/internal/entities"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)
}

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (id, name, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		if isUniqueConstraintError(err) {
			return entities.NewConflictError("user already exists", entities.ErrUserAlreadyExists)
		}
		return entities.NewInternalError("failed to create user", err)
	}

	return nil
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE id = $1`

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.NewNotFoundError("user not found", entities.ErrUserNotFound)
		}
		return nil, entities.NewInternalError("failed to get user by ID", err)
	}

	return user, nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE email = $1`

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.NewNotFoundError("user not found by email", entities.ErrUserNotFound)
		}
		return nil, entities.NewInternalError("failed to get user by email", err)
	}

	return user, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users
		SET name = $2, email = $3, updated_at = $4
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.UpdatedAt)
	if err != nil {
		if isUniqueConstraintError(err) {
			return entities.NewConflictError("email already in use", entities.ErrEmailAlreadyUsed)
		}
		return entities.NewInternalError("failed to update user", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return entities.NewInternalError("failed to get rows affected", err)
	}

	if rowsAffected == 0 {
		return entities.NewNotFoundError("user not found for update", entities.ErrUserNotFound)
	}

	return nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return entities.NewInternalError("failed to delete user", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return entities.NewInternalError("failed to get rows affected", err)
	}

	if rowsAffected == 0 {
		return entities.NewNotFoundError("user not found for deletion", entities.ErrUserNotFound)
	}

	return nil
}

func (r *UserRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, entities.NewInternalError("failed to list users", err)
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, entities.NewInternalError("failed to scan user", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, entities.NewInternalError("error iterating rows", err)
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
