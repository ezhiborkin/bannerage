package handler

import (
	"banners/domain/models"
	"banners/internal/errorwriter"
	"banners/lib/logger/sl"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

type CreateUserResponse struct {
	Message string `json:"message"`
	Email   string `json:"email"`
}

type UserProvider interface {
	CreateUser(email, role, password string) error
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	const op = "handler.createUser"

	log := h.log.With(slog.String("op", op))

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Error("failed to decode request body", sl.Err(err))
		errorwriter.WriteError(w, "failed to decode request", http.StatusBadRequest)
		return
	}
	if errors.Is(err, io.EOF) {
		log.Error("request body is empty", sl.Err(err))
		errorwriter.WriteError(w, "empty request", http.StatusBadRequest)
		return
	}

	log.Info("request body decoded")

	if user.Email == "" {
		h.log.Error("email is empty")
		errorwriter.WriteError(w, "email is empty", http.StatusBadRequest)
		return
	}
	if user.Role == "" {
		h.log.Error("role is empty")
		errorwriter.WriteError(w, "role is empty", http.StatusBadRequest)

		return
	}
	if user.Password == "" {
		h.log.Error("password is empty")
		errorwriter.WriteError(w, "password is empty", http.StatusBadRequest)
		return
	}

	err = h.userProvider.CreateUser(user.Email, user.Role, user.Password)
	if err != nil {
		log.Error("failed to create user", sl.Err(err))
		errorwriter.WriteError(w, "failed to create user", http.StatusConflict)
		return
	}

	response := CreateUserResponse{
		Message: "Successfully created user.",
		Email:   user.Email,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Error("failed to marshal response", sl.Err(err))
		errorwriter.WriteError(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJSON)
}
