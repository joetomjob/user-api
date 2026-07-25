package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeJsonError(w, "Not valid request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	u, err := h.service.Create(ctx, user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeJsonError(w, "username or email already exists", http.StatusConflict)
		} else if errors.Is(err, ErrNameRequired) || errors.Is(err, ErrEmailRequired) || errors.Is(err, ErrInvalidAge) {
			writeJsonError(w, "invalid input", http.StatusBadRequest)
		} else {
			writeJsonError(w, "failed to create user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(u)
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	user, err := h.service.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeJsonError(w, "No rows present", http.StatusNotFound)
		} else {
			writeJsonError(w, "failed to get user info", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeJsonError(w, "Not valid request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	u, err := h.service.Update(ctx, user, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeJsonError(w, "username or email already exists", http.StatusConflict)
		} else if errors.Is(err, ErrNotFound) {
			writeJsonError(w, "No record found", http.StatusNotFound)
		} else if errors.Is(err, ErrNameRequired) || errors.Is(err, ErrEmailRequired) || errors.Is(err, ErrInvalidAge) || errors.Is(err, ErrInvalidId) {
			writeJsonError(w, "invalid input", http.StatusBadRequest)
		} else {
			writeJsonError(w, "failed to update user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(u)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	err = h.service.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeJsonError(w, "No record found", http.StatusNotFound)
		} else if errors.Is(err, ErrInvalidId) {
			writeJsonError(w, "Invalid Id", http.StatusBadRequest)
		} else {
			writeJsonError(w, "failed to delete user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := map[string]string{"status": "OK"}
	json.NewEncoder(w).Encode(res)
}

func writeJsonError(w http.ResponseWriter, errormessage string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := map[string]string{"error": errormessage}
	json.NewEncoder(w).Encode(response)
}
