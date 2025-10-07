package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-clean-code/internal/dto"
	"go-clean-code/internal/usecase"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserUsecase is a mock implementation of UserUsecaseInterface
type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserUsecase) GetUser(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserUsecase) UpdateUser(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserUsecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserUsecase) ListUsers(ctx context.Context, limit, offset int) (*dto.ListUsersResponse, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ListUsersResponse), args.Error(1)
}

func TestUserHandler_CreateUser(t *testing.T) {
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	t.Run("should create user successfully", func(t *testing.T) {
		req := dto.CreateUserRequest{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		expectedResponse := &dto.UserResponse{
			ID:    uuid.New(),
			Name:  req.Name,
			Email: req.Email,
		}

		mockUsecase.On("CreateUser", mock.Anything, req).Return(expectedResponse, nil)

		reqBody, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		handler.CreateUser(recorder, request)

		assert.Equal(t, http.StatusCreated, recorder.Code)

		var response dto.UserResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.Name, response.Name)
		assert.Equal(t, expectedResponse.Email, response.Email)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid JSON", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte("invalid json")))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		handler.CreateUser(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("should return bad request for invalid input", func(t *testing.T) {
		req := dto.CreateUserRequest{
			Name:  "",
			Email: "john@example.com",
		}

		mockUsecase.On("CreateUser", mock.Anything, req).Return((*dto.UserResponse)(nil), usecase.ErrInvalidInput)

		reqBody, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		handler.CreateUser(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("should return conflict when email exists", func(t *testing.T) {
		mockUsecase := new(MockUserUsecase)
		handler := NewUserHandler(mockUsecase)

		req := dto.CreateUserRequest{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockUsecase.On("CreateUser", mock.Anything, req).Return((*dto.UserResponse)(nil), usecase.ErrEmailExists)

		reqBody, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		handler.CreateUser(recorder, request)

		assert.Equal(t, http.StatusConflict, recorder.Code)
		mockUsecase.AssertExpectations(t)
	})
}

func TestUserHandler_GetUser(t *testing.T) {
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	userID := uuid.New()

	t.Run("should get user successfully", func(t *testing.T) {
		expectedResponse := &dto.UserResponse{
			ID:    userID,
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockUsecase.On("GetUser", mock.Anything, userID).Return(expectedResponse, nil)

		request := httptest.NewRequest(http.MethodGet, "/users/"+userID.String(), nil)
		request = mux.SetURLVars(request, map[string]string{"id": userID.String()})
		recorder := httptest.NewRecorder()

		handler.GetUser(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response dto.UserResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Name, response.Name)
		assert.Equal(t, expectedResponse.Email, response.Email)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid UUID", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/users/invalid-uuid", nil)
		request = mux.SetURLVars(request, map[string]string{"id": "invalid-uuid"})
		recorder := httptest.NewRecorder()

		handler.GetUser(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("should return not found when user doesn't exist", func(t *testing.T) {
		mockUsecase := new(MockUserUsecase)
		handler := NewUserHandler(mockUsecase)

		mockUsecase.On("GetUser", mock.Anything, userID).Return((*dto.UserResponse)(nil), usecase.ErrUserNotFound)

		request := httptest.NewRequest(http.MethodGet, "/users/"+userID.String(), nil)
		request = mux.SetURLVars(request, map[string]string{"id": userID.String()})
		recorder := httptest.NewRecorder()

		handler.GetUser(recorder, request)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		mockUsecase.AssertExpectations(t)
	})
}

func TestUserHandler_UpdateUser(t *testing.T) {
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	userID := uuid.New()

	t.Run("should update user successfully", func(t *testing.T) {
		req := dto.UpdateUserRequest{
			Name: "John Smith",
		}

		expectedResponse := &dto.UserResponse{
			ID:    userID,
			Name:  req.Name,
			Email: "john@example.com",
		}

		mockUsecase.On("UpdateUser", mock.Anything, userID, req).Return(expectedResponse, nil)

		reqBody, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPut, "/users/"+userID.String(), bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")
		request = mux.SetURLVars(request, map[string]string{"id": userID.String()})
		recorder := httptest.NewRecorder()

		handler.UpdateUser(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response dto.UserResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.Name, response.Name)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid UUID", func(t *testing.T) {
		req := dto.UpdateUserRequest{Name: "John Smith"}
		reqBody, _ := json.Marshal(req)

		request := httptest.NewRequest(http.MethodPut, "/users/invalid-uuid", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")
		request = mux.SetURLVars(request, map[string]string{"id": "invalid-uuid"})
		recorder := httptest.NewRecorder()

		handler.UpdateUser(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}

func TestUserHandler_DeleteUser(t *testing.T) {
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	userID := uuid.New()

	t.Run("should delete user successfully", func(t *testing.T) {
		mockUsecase.On("DeleteUser", mock.Anything, userID).Return(nil)

		request := httptest.NewRequest(http.MethodDelete, "/users/"+userID.String(), nil)
		request = mux.SetURLVars(request, map[string]string{"id": userID.String()})
		recorder := httptest.NewRecorder()

		handler.DeleteUser(recorder, request)

		assert.Equal(t, http.StatusNoContent, recorder.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid UUID", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodDelete, "/users/invalid-uuid", nil)
		request = mux.SetURLVars(request, map[string]string{"id": "invalid-uuid"})
		recorder := httptest.NewRecorder()

		handler.DeleteUser(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}

func TestUserHandler_ListUsers(t *testing.T) {
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	t.Run("should list users successfully", func(t *testing.T) {
		expectedResponse := &dto.ListUsersResponse{
			Users: []*dto.UserResponse{
				{
					ID:    uuid.New(),
					Name:  "John Doe",
					Email: "john@example.com",
				},
			},
			Total:  1,
			Limit:  10,
			Offset: 0,
		}

		mockUsecase.On("ListUsers", mock.Anything, 0, 0).Return(expectedResponse, nil)

		request := httptest.NewRequest(http.MethodGet, "/users", nil)
		recorder := httptest.NewRecorder()

		handler.ListUsers(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response dto.ListUsersResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.Total, response.Total)
		assert.Len(t, response.Users, 1)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("should handle query parameters", func(t *testing.T) {
		expectedResponse := &dto.ListUsersResponse{
			Users:  []*dto.UserResponse{},
			Total:  0,
			Limit:  5,
			Offset: 10,
		}

		mockUsecase.On("ListUsers", mock.Anything, 5, 10).Return(expectedResponse, nil)

		request := httptest.NewRequest(http.MethodGet, "/users?limit=5&offset=10", nil)
		recorder := httptest.NewRecorder()

		handler.ListUsers(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		mockUsecase.AssertExpectations(t)
	})
}