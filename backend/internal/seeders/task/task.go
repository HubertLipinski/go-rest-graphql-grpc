package task

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/HubertLipinski/go-rest-graphql-grpc/internal/repository"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func TruncateTable(db *sql.DB) error {
	log.Print("Truncating tasks table")
	_, err := db.Exec(`TRUNCATE TABLE tasks RESTART IDENTITY`)
	if err != nil {
		return err
	}
	log.Print("Tasks table truncated")

	return nil
}

func SeedTasks(db *sql.DB) error {
	log.Print("Task seeder execution started")

	err := TruncateTable(db)
	if err != nil {
		return err
	}

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	jsonPath := filepath.Join(dir, "dummy_data.json")

	log.Print("Reading data from file: ", jsonPath)

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var tasks []repository.Task

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return err
	}

	log.Printf("Inserting %d tasks...", len(tasks))
	for _, t := range tasks {
		dueDate, err := time.Parse(time.DateTime, t.DueDate)
		if err != nil {
			return err
		}

		_, err = db.Exec(`INSERT INTO tasks (title, description, status, priority, due_date) VALUES ($1, $2, $3, $4, $5)`,
			t.Title, t.Description, t.Status, t.Priority, dueDate)
		if err != nil {
			return fmt.Errorf("error inserting task '%s': %w", t.Title, err)
		}

	}

	log.Print("Task seeder execution completed")

	return nil
}
