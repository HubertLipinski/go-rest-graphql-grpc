package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/HubertLipinski/go-rest-graphql-grpc/internal/database"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	// TODO: project_id, assigned_id
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func IsValidStatus(status string) bool {
	allowedStatuses := map[string]bool{
		"todo":        true,
		"done":        true,
		"in_progress": true,
	}

	return strings.TrimSpace(status) != "" && allowedStatuses[status]
}

func GetAllTasks(db *database.DBConnection, status string, dueDate string) ([]Task, error) {
	query := "SELECT * FROM tasks"

	var conditions []string
	var args []interface{}
	argId := 1

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argId))
		args = append(args, status)
		argId++
	}
	if dueDate != "" {
		conditions = append(conditions, fmt.Sprintf("due_date < $%d", argId))
		args = append(args, dueDate)
		argId++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := db.Instance.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func GetTaskById(db *database.DBConnection, taskId string) (*Task, error) {
	log.Print(taskId)
	query := "SELECT * FROM tasks WHERE id = $1"

	row := db.Instance.QueryRow(query, taskId)

	var t Task
	err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
}

func CreateTask(db *database.DBConnection, req *CreateTaskRequest) (*Task, error) {
	query := `
        INSERT INTO tasks (title, description, status, priority,  due_date)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

	var task Task
	err := db.Instance.QueryRow(query, req.Title, req.Description, req.Status, req.Priority, req.DueDate).
		Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}

	return &task, nil
}
