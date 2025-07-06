package seeders

import (
	"log"

	"github.com/HubertLipinski/go-rest-graphql-grpc/internal/database"
	"github.com/HubertLipinski/go-rest-graphql-grpc/internal/seeders/task"
)

func SeedDB(connection *database.DBConnection) error {
	log.Print("Seeding DB")

	err := task.SeedTasks(connection.Instance)
	if err != nil {
		return err
	}

	log.Print("Seeding DB completed")

	return nil
}
