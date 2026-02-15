package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sptringTwoRestAPI/internal/database"
	"sptringTwoRestAPI/internal/models"
	"strings"
)

type Handlers struct {
	store *database.TaskStore
}

func NewHandlers(store *database.TaskStore) *Handlers {
	return &Handlers{store: store}
}

func (h *Handlers) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err)
		return
	}
	RespondWithJSON(w, http.StatusOK, tasks)
}

func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := GetTaskID(r.URL.Path)
	if err != nil {
		RespondWithError(
			w, http.StatusBadRequest,
			fmt.Errorf("error converting path string to task id: %w", err),
		)
		return
	}

	tasks, err := h.store.GetByID(id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err)
		return
	}
	RespondWithJSON(w, http.StatusOK, tasks)
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTaskInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		RespondWithError(
			w, http.StatusBadRequest,
			fmt.Errorf("error decoding create task input: %w", err),
		)
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		RespondWithError(w, http.StatusBadRequest, ErrNoTaskTitle)
		return
	}

	task, err := h.store.Create(input)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err)
		return
	}
	RespondWithJSON(w, http.StatusCreated, task)
}

func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := GetTaskID(r.URL.Path)
	if err != nil {
		RespondWithError(
			w, http.StatusBadRequest,
			fmt.Errorf("error converting path string to task id: %w", err),
		)
		return
	}

	var input models.UpdateTask
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		RespondWithError(w, http.StatusBadRequest, ErrNoTaskTitle)
		return
	}

	task, err := h.store.Update(id, input)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrTaskNotFound):
			RespondWithError(w, http.StatusNotFound, err)
		default:
			RespondWithError(w, http.StatusInternalServerError, err)
		}
	}
	RespondWithJSON(w, http.StatusOK, task)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := GetTaskID(r.URL.Path)
	if err != nil {
		RespondWithError(
			w, http.StatusBadRequest,
			fmt.Errorf("error converting path string to task id: %w", err),
		)
		return
	}

	if err = h.store.Delete(id); err != nil {
		switch {
		case errors.Is(err, database.ErrTaskNotFound):
			RespondWithError(w, http.StatusNotFound, err)
		default:
			RespondWithError(w, http.StatusInternalServerError, err)
		}
		return
	}
	RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
