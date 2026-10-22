package handler

import (
	"context"
	"log/slog"

	"taskapi/internal/domain"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (domain.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	ParseToken(raw string) (int64, error)
}

type TaskService interface {
	Create(ctx context.Context, userID int64, title string) (domain.Task, error)
	List(ctx context.Context, userID int64) ([]domain.Task, error)
	Get(ctx context.Context, userID, id int64) (domain.Task, error)
	Update(ctx context.Context, userID, id int64, title *string, done *bool) (domain.Task, error)
	Delete(ctx context.Context, userID, id int64) error
}

type Handler struct {
	log   *slog.Logger
	auth  AuthService
	tasks TaskService
}

func New(log *slog.Logger, auth AuthService, tasks TaskService) *Handler {
	return &Handler{log: log, auth: auth, tasks: tasks}
}
