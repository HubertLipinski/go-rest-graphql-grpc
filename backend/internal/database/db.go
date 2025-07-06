package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type DBCredentials struct {
	host, port, user, password, dbName, dbSSLMode string
}

func (dbc *DBCredentials) ConnectionStr() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbc.host, dbc.port, dbc.user, dbc.password, dbc.dbName, dbc.dbSSLMode,
	)
}

type DBConnection struct {
	Instance *sql.DB
}

func LoadCredentials() (*DBCredentials, error) {
	return &DBCredentials{
			host:      os.Getenv("POSTGRES_HOST"),
			port:      os.Getenv("POSTGRES_PORT"),
			user:      os.Getenv("POSTGRES_USER"),
			password:  os.Getenv("POSTGRES_PASSWORD"),
			dbName:    os.Getenv("POSTGRES_DB"),
			dbSSLMode: "disable",
		},
		nil
}

func InitDBConnection() (*DBConnection, error) {
	credentials, err := LoadCredentials()
	if err != nil {
		log.Fatalf("Error loading credentials: %v", err)
	}

	db, err := sql.Open("postgres", credentials.ConnectionStr())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DBConnection{Instance: db}, nil
}

func (dbc *DBConnection) Close() {
	err := dbc.Instance.Close()
	if err != nil {
		log.Fatal(err)
	}
}
