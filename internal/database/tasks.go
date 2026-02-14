package database

import (
	"database/sql"
	"errors"
	"fmt"
	"sptringTwoRestAPI/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type TaskStore struct {
	db *sqlx.DB
}

func NewTaskStore(db *sqlx.DB) *TaskStore {
	return &TaskStore{db: db}
}

func (s *TaskStore) GetAll() ([]models.Task, error) {
	var tasks []models.Task

	query := `
select id, title, description, completed, created_at, updated_at 
from tasks 
order by created_at desc;`

	err := s.db.Select(&tasks, query)
	if err != nil {
		return nil, fmt.Errorf("error getting all tasks: %w", err)
	}
	return tasks, nil
}

func (s *TaskStore) GetByID(id int) (*models.Task, error) {
	var task models.Task

	query := `
select id, title, description, completed, created_at, updated_at 
from tasks 
where id = $1;`

	err := s.db.Get(&task, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("task with id %d not found: %w", id, err)
	}
	if err != nil {
		return nil, fmt.Errorf("error getting task by id: %w", err)
	}
	return &task, nil
}

func (s *TaskStore) Create(input models.CreateTaskInput) (*models.Task, error) {
	var task models.Task

	query := `
insert into tasks (title, description, completed, created_at, updated_at) values 
                                                                              ($1, $2, $3, $4, $5)
returning id, title, description, completed, created_at, updated_at;`
	now := time.Now()

	err := s.db.QueryRowx(query, input.Title, input.Description, input.Description, now, now).StructScan(&task)
	if err != nil {
		return nil, fmt.Errorf("error creating task: %w", err)
	}
	return &task, nil
}

func (s *TaskStore) Update(id int, input models.UpdateTask) (*models.Task, error) {
	task, err := s.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("error updating task (not found): %w", err)
	}
	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Completed != nil {
		task.Completed = *input.Completed
	}
	task.UpdatedAt = time.Now()

	query := `
update tasks 
    set title = $1, description = $2, completed = $3, updated_at = $4) 
    where id = $5
returning id, title, description, completed, created_at, updated_at;`

	var updatedTask models.Task

	err = s.db.QueryRowx(query, task.Title, task.Description, task.Description, task.UpdatedAt, id).StructScan(&updatedTask)
	if err != nil {
		return nil, fmt.Errorf("error updating task: %w", err)
	}
	return &updatedTask, nil
}
