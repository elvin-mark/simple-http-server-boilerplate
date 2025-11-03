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

func (h *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers()
	if err != nil {
		utils.WriteJSONStatus(w, map[string]string{"error": "Failed to get users"}, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, users)
}

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
