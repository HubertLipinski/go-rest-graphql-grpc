package main

import (
	"database/sql"
	"fmt"
	"github.com/HubertLipinski/go-rest-graphql-grpc/internal/seeders"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	log.Print("Hello from REST")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Błąd ładowania pliku .env: %v", err)
	}

	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbSSLMode := "disable"

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	err = seeders.SeedDB(db)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		return
	}
}
