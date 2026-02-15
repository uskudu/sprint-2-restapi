package main

import (
	"log"
	"net/http"
	"os"
	"sptringTwoRestAPI/internal/database"
	"sptringTwoRestAPI/internal/handlers"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	serverPort := os.Getenv("SERVER_PORT")

	log.Printf("starting server on port: %s", serverPort)

	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatalf("database didnt connect: %s", err)
	}

	log.Printf("database successfully connected: %s", databaseURL)
	defer db.Close()

	taskStore := database.NewTaskStore(db)
	handler := handlers.NewHandlers(taskStore)

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", methodHandler(handler.GetAllTasks, http.MethodGet))
	mux.HandleFunc("/tasks/create", methodHandler(handler.GetAllTasks, http.MethodPost))
	mux.HandleFunc("/tasks/", taskIDHandler(handler))

	loggedMux := loggingMiddleware(mux)
	serverAddr := ":" + serverPort

	err = http.ListenAndServe(serverAddr, loggedMux)
	if err != nil {
		log.Fatalf("server couldnt start: %s", err)
	} else {
		log.Printf("server started at %s", serverAddr)
	}
}
