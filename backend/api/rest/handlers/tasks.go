package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/HubertLipinski/go-rest-graphql-grpc/api/rest/response"
	"github.com/HubertLipinski/go-rest-graphql-grpc/internal/repository"
	"net/http"
	"strconv"

	"github.com/HubertLipinski/go-rest-graphql-grpc/internal/database"
)

func GetAllTasks(db *database.DBConnection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		dueDate := r.URL.Query().Get("due_date")

		if status != "" && !repository.IsValidStatus(status) {
			response.Error(w, "Invalid status filter. Allowed values: todo, done, in_progress", http.StatusBadRequest)
			return
		}

		tasks, err := repository.GetAllTasks(db, status, dueDate)

		if err != nil {
			response.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
			return
		}

		response.Success(w, tasks, http.StatusOK)
	}
}

func GetTasksById(db *database.DBConnection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			response.Error(w, "Invalid id format", http.StatusBadRequest)
		}

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
			response.Error(w, fmt.Sprintf("Invalid reqest body: %v", err), http.StatusBadRequest)
			return
		}

		if !repository.IsValidStatus(req.Status) {
			response.Error(w, "Invalid status", http.StatusBadRequest)
			return
		}

		task, err := repository.CreateTask(db, &req)
		if err != nil {
			response.Error(w, fmt.Sprintf("Server error %v", err), http.StatusInternalServerError)
			return
		}

		response.Success(w, task, http.StatusCreated)
	}
}

func DeleteTask(db *database.DBConnection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.Error(w, "Invalid id format", http.StatusBadRequest)
			return
		}

		err = repository.DeleteTaskById(db, id)
		if err != nil {
			response.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		response.Success(w, "Task deleted", http.StatusOK)
	}
}
