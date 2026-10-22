package domain

import "time"

type Task struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"-"`
	Title     string    `db:"title" json:"title"`
	Done      bool      `db:"done" json:"done"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
