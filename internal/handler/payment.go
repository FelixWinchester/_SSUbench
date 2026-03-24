package handler

import (
	"net/http"

	"github.com/FelixWinchester/ssubench/internal/middleware"
	"github.com/google/uuid"
)

func (h *Handler) listPayments(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	limit, offset := getPagination(r)

	payments, err := h.payment.ListByUser(r.Context(), userID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, payments)
}
