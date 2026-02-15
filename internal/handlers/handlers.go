package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sptringTwoRestAPI/internal/database"
	"sptringTwoRestAPI/internal/models"
	"sptringTwoRestAPI/internal/utils"
	"strconv"
	"strings"
)

type Handlers struct {
	store *database.TaskStore
}

func NewHandlers(store *database.TaskStore) *Handlers {
	return &Handlers{store: store}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

// todo err error is message string at 1:14:14
func respondWithError(w http.ResponseWriter, statusCode int, err error) {
	respondWithJSON(w, statusCode, map[string]error{"error": err})
}

func (h *Handlers) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll()
	if err != nil {
		respondWithError(
			w, http.StatusInternalServerError,
			fmt.Errorf("error getting all tasks: %w", err),
		)
		return
	}
	respondWithJSON(w, http.StatusOK, tasks)
}

func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")

	id, err := strconv.Atoi(pathParts[0])
	if err != nil {
		respondWithError(
			w, http.StatusBadRequest,
			fmt.Errorf("error converting path string to task id: %w", err),
		)
		return
	}

	tasks, err := h.store.GetByID(id)
	if err != nil {
		respondWithError(
			w, http.StatusBadRequest,
			fmt.Errorf("error getting a task: %w", err),
		)
		return
	}
	respondWithJSON(w, http.StatusOK, tasks)
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTaskInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(
			w, http.StatusBadRequest,
			fmt.Errorf("error decoding create task input: %w", err),
		)
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		respondWithError(
			w, http.StatusBadRequest,
			fmt.Errorf("error creating task: %w", utils.ErrNoTaskTitle),
		)
		return
	}

	task, err := h.store.Create(input)
	if err != nil {
		respondWithError(
			w, http.StatusInternalServerError,
			fmt.Errorf("error creating a task: %w", err),
		)
		return
	}
	respondWithJSON(w, http.StatusCreated, task)
}
