package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"taskapi/internal/domain"
)

const taskColumns = `id, user_id, title, done, created_at`

type TaskRepo struct {
	pool *pgxpool.Pool
}

func NewTaskRepo(pool *pgxpool.Pool) *TaskRepo {
	return &TaskRepo{pool: pool}
}

func (r *TaskRepo) Create(ctx context.Context, userID int64, title string) (domain.Task, error) {
	rows, err := r.pool.Query(ctx,
		`INSERT INTO tasks (user_id, title) VALUES ($1, $2) RETURNING `+taskColumns,
		userID, title)
	if err != nil {
		return domain.Task{}, err
	}
	return pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Task])
}

func (r *TaskRepo) List(ctx context.Context, userID int64) ([]domain.Task, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+taskColumns+` FROM tasks WHERE user_id = $1 ORDER BY id`,
		userID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Task])
}

func (r *TaskRepo) Get(ctx context.Context, userID, id int64) (domain.Task, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+taskColumns+` FROM tasks WHERE id = $1 AND user_id = $2`,
		id, userID)
	if err != nil {
		return domain.Task{}, err
	}
	t, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Task])
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Task{}, domain.ErrNotFound
	}
	return t, err
}

func (r *TaskRepo) Update(ctx context.Context, userID, id int64, title *string, done *bool) (domain.Task, error) {
	rows, err := r.pool.Query(ctx,
		`UPDATE tasks
		 SET title = COALESCE($1::text, title),
		     done  = COALESCE($2::boolean, done)
		 WHERE id = $3 AND user_id = $4
		 RETURNING `+taskColumns,
		title, done, id, userID)
	if err != nil {
		return domain.Task{}, err
	}
	t, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Task])
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Task{}, domain.ErrNotFound
	}
	return t, err
}

func (r *TaskRepo) Delete(ctx context.Context, userID, id int64) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM tasks WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
