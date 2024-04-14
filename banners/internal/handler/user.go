package handler

import (
	"banners/domain/models"
	"banners/lib/logger/sl"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

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
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}
	if errors.Is(err, io.EOF) {
		log.Error("request body is empty", sl.Err(err))
		http.Error(w, "empty request", http.StatusBadRequest)
		return
	}

	log.Info("request body decoded")

	if user.Email == "" {
		h.log.Error("email is empty")
		http.Error(w, "email is empty", http.StatusBadRequest)
		return
	}
	if user.Role == "" {
		h.log.Error("role is empty")
		http.Error(w, "role is empty", http.StatusBadRequest)
		return
	}
	if user.Password == "" {
		h.log.Error("password is empty")
		http.Error(w, "password is empty", http.StatusBadRequest)
		return
	}

	err = h.userProvider.CreateUser(user.Email, user.Role, user.Password)
	if err != nil {
		log.Error("failed to create user", sl.Err(err))
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Created user with email - " + user.Email))
}
