package handler

import "net/http"

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var in credentials
	if err := decode(w, r, &in); err != nil {
		h.fail(w, r, err)
		return
	}
	u, err := h.auth.Register(r.Context(), in.Email, in.Password)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": u.ID, "email": u.Email})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var in credentials
	if err := decode(w, r, &in); err != nil {
		h.fail(w, r, err)
		return
	}
	token, err := h.auth.Login(r.Context(), in.Email, in.Password)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}
