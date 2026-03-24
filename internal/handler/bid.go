package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/FelixWinchester/ssubench/internal/domain"
	"github.com/FelixWinchester/ssubench/internal/middleware"
	"github.com/FelixWinchester/ssubench/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) createBid(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	executorID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	var input service.CreateBidInput
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	input.TaskID = taskID
	input.ExecutorID = executorID

	if err = h.validate.Struct(input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	bid, err := h.bid.Create(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrTaskNotPublished):
			writeError(w, http.StatusConflict, err.Error())
		case errors.Is(err, domain.ErrBidAlreadyExists):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, bid)
}

func (h *Handler) listBids(w http.ResponseWriter, r *http.Request) {
	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	limit, offset := getPagination(r)

	bids, err := h.bid.ListByTask(r.Context(), taskID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, bids)
}

func (h *Handler) acceptBid(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	customerID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	bidID, err := uuid.Parse(chi.URLParam(r, "bid_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid bid id")
		return
	}

	if err = h.bid.Accept(r.Context(), customerID, taskID, bidID); err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrBidNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, domain.ErrTaskNotPublished):
			writeError(w, http.StatusConflict, err.Error())
		case errors.Is(err, domain.ErrBidAlreadyAccepted):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "bid accepted"})
}
