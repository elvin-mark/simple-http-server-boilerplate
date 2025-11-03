package handlers

import (
	"encoding/json"
	"http-server/dto/user"
	"http-server/services"
	"http-server/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ============== STRUCTS ==============

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// ============== METHODS ==============

// GetUsersHandler godoc
// @Summary Get all users
// @Description Get a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} user.User
// @Failure 500 {object} map[string]string
// @Router /users [get]
func (h *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers()
	if err != nil {
		utils.WriteJSONStatus(w, map[string]string{"error": "Failed to get users"}, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, users)
}

// GetUserHandler godoc
// @Summary Get a user by ID
// @Description Get a single user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} user.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteJSONStatus(w, map[string]string{"error": "Invalid user ID"}, http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUser(id)
	if err != nil {
		utils.WriteJSONStatus(w, map[string]string{"error": "User not found"}, http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, user)
}

// CreateUserHandler godoc
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body user.CreateUserRequest true "User object to be created"
// @Success 201 {object} user.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req user.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONStatus(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
		return
	}

	createdUser, err := h.service.CreateUser(&req)
	if err != nil {
		utils.WriteJSONStatus(w, map[string]string{"error": "Failed to create user"}, http.StatusInternalServerError)
		return
	}

	utils.Logger.Info("User created", "id", createdUser.ID)
	utils.WriteJSONStatus(w, createdUser, http.StatusCreated)
}

// DeleteUserHandler godoc
// @Summary Delete a user by ID
// @Description Delete a single user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteJSONStatus(w, map[string]string{"error": "Invalid user ID"}, http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		utils.WriteJSONStatus(w, map[string]string{"error": "Failed to delete user"}, http.StatusInternalServerError)
		return
	}

	utils.Logger.Info("User deleted", "id", id)
	w.WriteHeader(http.StatusNoContent)
}
