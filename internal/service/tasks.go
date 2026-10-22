package service

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"taskapi/internal/domain"
)

type TaskStore interface {
	Create(ctx context.Context, userID int64, title string) (domain.Task, error)
	List(ctx context.Context, userID int64) ([]domain.Task, error)
	Get(ctx context.Context, userID, id int64) (domain.Task, error)
	Update(ctx context.Context, userID, id int64, title *string, done *bool) (domain.Task, error)
	Delete(ctx context.Context, userID, id int64) error
}

type TaskService struct {
	repo TaskStore
}

func NewTaskService(repo TaskStore) *TaskService {
	return &TaskService{repo: repo}
}

func validateTitle(title string) (string, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return "", fmt.Errorf("%w: title is required", domain.ErrValidation)
	}
	if utf8.RuneCountInString(title) > 200 {
		return "", fmt.Errorf("%w: title is too long (max 200)", domain.ErrValidation)
	}
	return title, nil
}

func (s *TaskService) Create(ctx context.Context, userID int64, title string) (domain.Task, error) {
	title, err := validateTitle(title)
	if err != nil {
		return domain.Task{}, err
	}
	return s.repo.Create(ctx, userID, title)
}

func (s *TaskService) List(ctx context.Context, userID int64) ([]domain.Task, error) {
	tasks, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, err
	}
	if tasks == nil {
		tasks = []domain.Task{}
	}
	return tasks, nil
}

func (s *TaskService) Get(ctx context.Context, userID, id int64) (domain.Task, error) {
	return s.repo.Get(ctx, userID, id)
}

func (s *TaskService) Update(ctx context.Context, userID, id int64, title *string, done *bool) (domain.Task, error) {
	if title == nil && done == nil {
		return domain.Task{}, fmt.Errorf("%w: nothing to update", domain.ErrValidation)
	}
	if title != nil {
		t, err := validateTitle(*title)
		if err != nil {
			return domain.Task{}, err
		}
		title = &t
	}
	return s.repo.Update(ctx, userID, id, title, done)
}

func (s *TaskService) Delete(ctx context.Context, userID, id int64) error {
	return s.repo.Delete(ctx, userID, id)
}
