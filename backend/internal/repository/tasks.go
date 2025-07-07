package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/HubertLipinski/go-rest-graphql-grpc/internal/database"
)

type DateTimeFormat struct {
	time.Time
}

func (ct *DateTimeFormat) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1]
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

func (ct DateTimeFormat) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.Format(time.DateTime))
}

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	// TODO: project_id, assigned_id
	Description string         `json:"description"`
	Status      string         `json:"status"`
	Priority    string         `json:"priority"`
	DueDate     DateTimeFormat `json:"due_date"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

func IsValidStatus(status string) bool {
	allowedStatuses := map[string]bool{
		"todo":        true,
		"done":        true,
		"in_progress": true,
	}

	return allowedStatuses[status]
}

func GetAllTasks(db *database.DBConnection, status string, dueDate string) ([]*Task, error) {
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

	var tasks []*Task
	for rows.Next() {
		task, err := scanInoTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if len(tasks) == 0 {
		tasks = []*Task{}
	}

	return tasks, nil
}

func GetTaskById(db *database.DBConnection, taskId int) (*Task, error) {
	log.Print(taskId)
	rows, err := db.Instance.Query("SELECT * FROM tasks WHERE id = $1", taskId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanInoTask(rows)
	}

	return nil, fmt.Errorf("task %v not found", taskId)
}

type CreateTaskRequest struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	Priority    string         `json:"priority"`
	DueDate     DateTimeFormat `json:"due_date"`
}

func CreateTask(db *database.DBConnection, req *CreateTaskRequest) (*Task, error) {
	query := `
        INSERT INTO tasks (title, description, status, priority,  due_date)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING *
    `

	var task Task
	err := db.Instance.QueryRow(query, req.Title, req.Description, req.Status, req.Priority, req.DueDate.Time).
		Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.DueDate.Time,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func DeleteTaskById(db *database.DBConnection, taskId int) error {
	// TODO: Check permissions
	_, err := db.Instance.Exec("DELETE FROM tasks WHERE id = $1", taskId)

	if err != nil {
		return err
	}

	return nil
}

func scanInoTask(rows *sql.Rows) (*Task, error) {
	task := new(Task)

	err := rows.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate.Time,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	return task, err
}
