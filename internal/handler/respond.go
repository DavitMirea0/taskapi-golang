package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"taskapi/internal/domain"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func errBody(msg string) map[string]string {
	return map[string]string{"error": msg}
}

func decode(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("%w: malformed json", domain.ErrValidation)
	}
	return nil
}

func (h *Handler) fail(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		writeJSON(w, http.StatusBadRequest, errBody(err.Error()))
	case errors.Is(err, domain.ErrInvalidCredentials):
		writeJSON(w, http.StatusUnauthorized, errBody("invalid credentials"))
	case errors.Is(err, domain.ErrNotFound):
		writeJSON(w, http.StatusNotFound, errBody("not found"))
	case errors.Is(err, domain.ErrEmailTaken):
		writeJSON(w, http.StatusConflict, errBody("email already taken"))
	default:
		h.log.Error("internal error", "err", err, "method", r.Method, "path", r.URL.Path)
		writeJSON(w, http.StatusInternalServerError, errBody("internal error"))
	}
}
