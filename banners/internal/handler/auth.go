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

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type AuthProvider interface {
	LoginUser(email, password string) (string, error)
}

func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request) {
	const op = "handler.loginUser"

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
		errorwriter.WriteError(w, "request body is empty", http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" {
		log.Error("email or password is empty")
		errorwriter.WriteError(w, "email or password is empty", http.StatusBadRequest)
		return
	}

	log.Info("request body decoded")

	tokenString, err := h.authProvider.LoginUser(user.Email, user.Password)
	if err != nil {
		log.Error("failed to login user", sl.Err(err))
		errorwriter.WriteError(w, "failed to login user", http.StatusBadRequest)
		return
	}

	response := LoginResponse{
		Message: "Successfully logged in.",
		Token:   tokenString,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		errorwriter.WriteError(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+tokenString)
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
