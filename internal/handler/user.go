package handler

import (
	"errors"
	"net/http"

	"github.com/FelixWinchester/ssubench/internal/domain"
	"github.com/FelixWinchester/ssubench/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) getMe(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.user.GetByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":         user.ID,
		"email":      user.Email,
		"role":       user.Role,
		"balance":    user.Balance,
		"is_blocked": user.IsBlocked,
		"created_at": user.CreatedAt,
	})
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	limit, offset := getPagination(r)

	users, err := h.user.List(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, users)
}

func (h *Handler) blockUser(w http.ResponseWriter, r *http.Request) {
	adminIDStr := middleware.GetUserID(r.Context())
	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if err = h.user.Block(r.Context(), adminID, userID); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "user blocked"})
}

func (h *Handler) unblockUser(w http.ResponseWriter, r *http.Request) {
	adminIDStr := middleware.GetUserID(r.Context())
	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if err = h.user.Unblock(r.Context(), adminID, userID); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "user unblocked"})
}
