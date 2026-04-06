package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aditya-sutar-45/finance-dashboard-api/internal/database"
	"github.com/aditya-sutar-45/finance-dashboard-api/models"
	"github.com/aditya-sutar-45/finance-dashboard-api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("error parsing json body: %v", err))
	}

	passwordHash, err := utils.HashPassword(params.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("error hashing password: %v", err))
		return
	}

	user, err := h.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:           uuid.New(),
		Name:         params.Name,
		Email:        params.Email,
		PasswordHash: passwordHash,
		Role:         params.Role,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("could not create user: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, models.DatabaseUserToUser(user))
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.DB.ListUsers(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("could not get users list: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.DatabaseUsersToUsers(users))
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")

	id, err := uuid.Parse(idString)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("error parsing uuid: %v", err))
		return
	}

	type parameters struct {
		Name         string `json:"name"`
		Email        string `json:"email"`
		PasswordHash string `json:"password_hash"`
		Role         string `json:"role"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("error parsing json body: %v", err))
		return
	}

	user, err := h.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:           id,
		Name:         params.Name,
		Email:        params.Email,
		PasswordHash: params.PasswordHash,
		Role:         params.Role,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("could not update user: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, models.DatabaseUserToUser(user))
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")

	id, err := uuid.Parse(idString)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("error parsing uuid: %v", err))
		return
	}

	err = h.DB.DeleteUser(r.Context(), id)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("error deleting user from db: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, struct{}{})
}

func (h *Handler) GetDeletedUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.DB.GetDeletedUsers(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error fetching users from DB:\n %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.DatabaseUsersToUsers(users))
}
