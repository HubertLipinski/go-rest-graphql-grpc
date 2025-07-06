package seeders

import (
	"database/sql"
	"log"
	
	"github.com/HubertLipinski/go-rest-graphql-grpc/internal/seeders/task"
)

func SeedDB(db *sql.DB) error {

	log.Print("Seeding DB")

	err := task.SeedTasks(db)
	if err != nil {
		return err
	}

	log.Print("Seeding DB completed")

	return nil
}
