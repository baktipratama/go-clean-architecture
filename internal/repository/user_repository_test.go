package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"go-clean-code/internal/entities"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestUserRepositoryImpl_Create(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	require.NoError(t, err)
	defer db.Close()

	repo := &UserRepositoryImpl{db: db.DB}
	ctx := context.Background()

	user := &entities.User{
		ID:        uuid.New(),
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("should create user successfully", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO users \(id, name, email, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5\)`).
			WithArgs(user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt).
			WillReturnResult(sqlxmock.NewResult(1, 1))

		err := repo.Create(ctx, user)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when email exists", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO users \(id, name, email, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5\)`).
			WithArgs(user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt).
			WillReturnError(&testError{msg: "duplicate key value violates unique constraint"})

		err := repo.Create(ctx, user)
		assert.True(t, entities.IsConflictError(err))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepositoryImpl_GetByID(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	require.NoError(t, err)
	defer db.Close()

	repo := &UserRepositoryImpl{db: db.DB}
	ctx := context.Background()

	userID := uuid.New()
	user := &entities.User{
		ID:        userID,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("should return user when exists", func(t *testing.T) {
		rows := sqlxmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
			AddRow(user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)

		mock.ExpectQuery(`SELECT id, name, email, created_at, updated_at FROM users WHERE id = \$1`).
			WithArgs(userID).
			WillReturnRows(rows)

		foundUser, err := repo.GetByID(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.Name, foundUser.Name)
		assert.Equal(t, user.Email, foundUser.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, name, email, created_at, updated_at FROM users WHERE id = \$1`).
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		foundUser, err := repo.GetByID(ctx, userID)
		assert.True(t, entities.IsNotFoundError(err))
		assert.Nil(t, foundUser)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepositoryImpl_Update(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	require.NoError(t, err)
	defer db.Close()

	repo := &UserRepositoryImpl{db: db.DB}
	ctx := context.Background()

	user := &entities.User{
		ID:        uuid.New(),
		Name:      "John Smith",
		Email:     "john.smith@example.com",
		UpdatedAt: time.Now(),
	}

	t.Run("should update user successfully", func(t *testing.T) {
		mock.ExpectExec(`UPDATE users SET name = \$2, email = \$3, updated_at = \$4 WHERE id = \$1`).
			WithArgs(user.ID, user.Name, user.Email, user.UpdatedAt).
			WillReturnResult(sqlxmock.NewResult(0, 1))

		err := repo.Update(ctx, user)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		mock.ExpectExec(`UPDATE users SET name = \$2, email = \$3, updated_at = \$4 WHERE id = \$1`).
			WithArgs(user.ID, user.Name, user.Email, user.UpdatedAt).
			WillReturnResult(sqlxmock.NewResult(0, 0))

		err := repo.Update(ctx, user)
		assert.True(t, entities.IsNotFoundError(err))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepositoryImpl_Delete(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	require.NoError(t, err)
	defer db.Close()

	repo := &UserRepositoryImpl{db: db.DB}
	ctx := context.Background()

	userID := uuid.New()

	t.Run("should delete user successfully", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
			WithArgs(userID).
			WillReturnResult(sqlxmock.NewResult(0, 1))

		err := repo.Delete(ctx, userID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
			WithArgs(userID).
			WillReturnResult(sqlxmock.NewResult(0, 0))

		err := repo.Delete(ctx, userID)
		assert.True(t, entities.IsNotFoundError(err))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestIsUniqueConstraintError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "should detect duplicate key error",
			err:      &testError{msg: "duplicate key value violates unique constraint"},
			expected: true,
		},
		{
			name:     "should detect unique constraint error",
			err:      &testError{msg: "unique constraint violation"},
			expected: true,
		},
		{
			name:     "should not detect regular error",
			err:      &testError{msg: "connection failed"},
			expected: false,
		},
		{
			name:     "should handle nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isUniqueConstraintError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContainsSubstring(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected bool
	}{
		{
			name:     "should find substring in middle",
			s:        "hello world test",
			substr:   "world",
			expected: true,
		},
		{
			name:     "should find substring at beginning",
			s:        "world hello",
			substr:   "world",
			expected: true,
		},
		{
			name:     "should find substring at end",
			s:        "hello world",
			substr:   "world",
			expected: true,
		},
		{
			name:     "should not find non-existent substring",
			s:        "hello test",
			substr:   "world",
			expected: false,
		},
		{
			name:     "should handle empty substring",
			s:        "hello",
			substr:   "",
			expected: true,
		},
		{
			name:     "should handle empty string",
			s:        "",
			substr:   "world",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsSubstring(tt.s, tt.substr)
			assert.Equal(t, tt.expected, result)
		})
	}
}
