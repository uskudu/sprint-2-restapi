package models

import "time"

type Task struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Competed    bool      `json:"competed" db:"competed"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type UpdateTask struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}
