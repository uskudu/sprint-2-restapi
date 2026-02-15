package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func RespondWithError(w http.ResponseWriter, statusCode int, err error) {
	RespondWithJSON(w, statusCode, map[string]error{"error": err})
}

func GetTaskID(path string) (int, error) {
	pathParts := strings.Split(strings.TrimPrefix(path, "/tasks/"), "/")

	id, err := strconv.Atoi(pathParts[0])
	if err != nil {
		return -1, err
	}
	return id, nil
}
