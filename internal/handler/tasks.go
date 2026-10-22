package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"taskapi/internal/domain"
)

func idParam(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("%w: bad id", domain.ErrValidation)
	}
	return id, nil
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Title string `json:"title"`
	}
	if err := decode(w, r, &in); err != nil {
		h.fail(w, r, err)
		return
	}
	t, err := h.tasks.Create(r.Context(), userIDFrom(r.Context()), in.Title)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, t)
}

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.tasks.List(r.Context(), userIDFrom(r.Context()))
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *Handler) getTask(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	t, err := h.tasks.Get(r.Context(), userIDFrom(r.Context()), id)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var in struct {
		Title *string `json:"title"`
		Done  *bool   `json:"done"`
	}
	if err := decode(w, r, &in); err != nil {
		h.fail(w, r, err)
		return
	}
	t, err := h.tasks.Update(r.Context(), userIDFrom(r.Context()), id, in.Title, in.Done)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	if err := h.tasks.Delete(r.Context(), userIDFrom(r.Context()), id); err != nil {
		h.fail(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
