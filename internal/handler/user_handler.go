package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go-clean-code/internal/dto"
	"go-clean-code/internal/entities"
	"go-clean-code/internal/usecase"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userUsecase usecase.UserUsecaseInterface
}

func NewUserHandler(userUsecase usecase.UserUsecaseInterface) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// handleError handles domain errors and maps them to appropriate HTTP responses
func (h *UserHandler) handleError(w http.ResponseWriter, err error) {
	switch err {
	case usecase.ErrInvalidInput:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case usecase.ErrEmailExists:
		http.Error(w, err.Error(), http.StatusConflict)
	case usecase.ErrUserNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		switch {
		case entities.IsValidationError(err):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case entities.IsNotFoundError(err):
			http.Error(w, err.Error(), http.StatusNotFound)
		case entities.IsConflictError(err):
			http.Error(w, err.Error(), http.StatusConflict)
		case entities.IsInternalError(err):
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.userUsecase.CreateUser(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userUsecase.GetUser(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.userUsecase.UpdateUser(r.Context(), id, req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.userUsecase.DeleteUser(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	users, err := h.userUsecase.ListUsers(r.Context(), limit, offset)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
