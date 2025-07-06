package handlers

import (
	"encoding/json"
	"github.com/HubertLipinski/go-rest-graphql-grpc/api/rest/response"
	"github.com/HubertLipinski/go-rest-graphql-grpc/internal/repository"
	"net/http"

	"github.com/HubertLipinski/go-rest-graphql-grpc/internal/database"
)

func GetAllTasks(db *database.DBConnection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		dueBefore := r.URL.Query().Get("due_before")

		if !repository.IsValidStatus(status) {
			response.Error(w, "Invalid status filter. Allowed values: todo, done, in_progress", http.StatusBadRequest)
			return
		}

		tasks, err := repository.GetAllTasks(db, status, dueBefore)

		if err != nil {
			response.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
			return
		}

		response.Success(w, tasks, http.StatusOK)
	}
}

func GetTasksById(db *database.DBConnection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		task, err := repository.GetTaskById(db, id)
		if err != nil {
			response.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		response.Success(w, task, http.StatusOK)
	}
}

func CreateTask(db *database.DBConnection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req repository.CreateTaskRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			response.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if !repository.IsValidStatus(req.Status) {
			response.Error(w, "Invalid status", http.StatusBadRequest)
			return
		}
	}
}
