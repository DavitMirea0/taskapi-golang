package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.Recoverer)
	r.Use(requestLogger(h.log))

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Post("/auth/register", h.register)
	r.Post("/auth/login", h.login)

	r.Group(func(r chi.Router) {
		r.Use(h.authenticate)
		r.Route("/tasks", func(r chi.Router) {
			r.Post("/", h.createTask)
			r.Get("/", h.listTasks)
			r.Get("/{id}", h.getTask)
			r.Patch("/{id}", h.updateTask)
			r.Delete("/{id}", h.deleteTask)
		})
	})

	return r
}
