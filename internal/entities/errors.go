package entities

import (
	"errors"
	"fmt"
)

// Domain errors
var (
	ErrInvalidName       = errors.New("invalid name: name cannot be empty")
	ErrInvalidEmail      = errors.New("invalid email: email must be valid format")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrEmailAlreadyUsed  = errors.New("email is already in use")
)

// DomainError represents a domain-specific error with additional context
type DomainError struct {
	Type    ErrorType
	Message string
	Cause   error
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Cause
}

// ErrorType represents different types of domain errors
type ErrorType string

const (
	ValidationError ErrorType = "VALIDATION_ERROR"
	NotFoundError   ErrorType = "NOT_FOUND_ERROR"
	ConflictError   ErrorType = "CONFLICT_ERROR"
	InternalError   ErrorType = "INTERNAL_ERROR"
)

// NewValidationError creates a new validation error
func NewValidationError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    ValidationError,
		Message: message,
		Cause:   cause,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    NotFoundError,
		Message: message,
		Cause:   cause,
	}
}

// NewConflictError creates a new conflict error
func NewConflictError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    ConflictError,
		Message: message,
		Cause:   cause,
	}
}

// NewInternalError creates a new internal error
func NewInternalError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    InternalError,
		Message: message,
		Cause:   cause,
	}
}

// IsValidationError checks if error is a validation error
func IsValidationError(err error) bool {
	return isErrorType(err, ValidationError)
}

// IsNotFoundError checks if error is a not found error
func IsNotFoundError(err error) bool {
	return isErrorType(err, NotFoundError)
}

// IsConflictError checks if error is a conflict error
func IsConflictError(err error) bool {
	return isErrorType(err, ConflictError)
}

// IsInternalError checks if error is an internal error
func IsInternalError(err error) bool {
	return isErrorType(err, InternalError)
}

func isErrorType(err error, errorType ErrorType) bool {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.Type == errorType
	}
	return false
}
