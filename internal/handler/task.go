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

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	customerID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input service.CreateTaskInput
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	input.CustomerID = customerID

	if err = h.validate.Struct(input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.task.Create(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, task)
}

func (h *Handler) getTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	task, err := h.task.GetByID(r.Context(), taskID)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	limit, offset := getPagination(r)

	tasks, err := h.task.List(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}

func (h *Handler) publishTask(w http.ResponseWriter, r *http.Request) {
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

	if err = h.task.Publish(r.Context(), customerID, taskID); err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, domain.ErrTaskInvalidStatus):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "task published"})
}

func (h *Handler) cancelTask(w http.ResponseWriter, r *http.Request) {
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

	if err = h.task.Cancel(r.Context(), customerID, taskID); err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, domain.ErrTaskAlreadyCompleted):
			writeError(w, http.StatusConflict, err.Error())
		case errors.Is(err, domain.ErrTaskCancelled):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "task cancelled"})
}

func (h *Handler) completeTask(w http.ResponseWriter, r *http.Request) {
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

	if err = h.task.MarkCompleted(r.Context(), executorID, taskID); err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, domain.ErrTaskNotInProgress):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "task marked as completed"})
}

func (h *Handler) confirmTask(w http.ResponseWriter, r *http.Request) {
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

	if err = h.task.Confirm(r.Context(), customerID, taskID); err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, domain.ErrTaskInvalidStatus):
			writeError(w, http.StatusConflict, err.Error())
		case errors.Is(err, domain.ErrInsufficientBalance):
			writeError(w, http.StatusPaymentRequired, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "task confirmed, payment processed"})
}
