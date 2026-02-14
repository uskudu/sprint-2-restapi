package database

import (
	"database/sql"
	"errors"
	"fmt"
	"sptringTwoRestAPI/internal/models"

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
