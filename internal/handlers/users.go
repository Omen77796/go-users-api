package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/omen77796/go-users-api/internal/models"
	"github.com/omen77796/go-users-api/internal/services"
	"github.com/omen77796/go-users-api/internal/utils"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "error getting users")
		return
	}

	utils.JSON(w, http.StatusOK, users)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	u, err := decodeAndValidateUser(r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Create(&u)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, u)
}

func decodeAndValidateUser(r *http.Request) (models.User, error) {
	var u models.User

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&u); err != nil {
		return models.User{}, fmt.Errorf("invalid JSON: %w", err)
	}

	var extra json.RawMessage
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return models.User{}, errors.New("only one JSON object allowed")
	}

	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(u.Email)

	if u.Name == "" {
		return models.User{}, errors.New("name is required")
	}

	if u.Email == "" {
		return models.User{}, errors.New("email is required")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return models.User{}, errors.New("invalid email format")
	}

	return u, nil
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUserID(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "user not found")
		return
	}

	utils.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseUserID(chi.URLParam(r, "id"))
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "user not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseUserID(rawID string) (int, error) {
	id, err := strconv.Atoi(rawID)
	if err != nil || id <= 0 {
		return 0, errors.New("id must be a positive integer")
	}
	return id, nil

}
