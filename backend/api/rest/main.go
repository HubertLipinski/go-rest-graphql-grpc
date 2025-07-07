package main

import (
	"github.com/HubertLipinski/go-rest-graphql-grpc/api/rest/middleware"
	"github.com/HubertLipinski/go-rest-graphql-grpc/internal/metrics"
	"log"
	"net/http"

	"github.com/HubertLipinski/go-rest-graphql-grpc/api/rest/handlers"
	"github.com/HubertLipinski/go-rest-graphql-grpc/internal/config"
	"github.com/HubertLipinski/go-rest-graphql-grpc/internal/database"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.Print("Starting REST API")

	config.LoadEnv()

	metrics.Init()

	connection, err := database.InitDBConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	// TODO: uncomment
	//err = seeders.SeedDB(connection)
	//if err != nil {
	//	log.Fatal(err)
	//}

	router := http.NewServeMux()

	router.HandleFunc("GET /tasks", handlers.GetAllTasks(connection))
	router.HandleFunc("POST /tasks", handlers.CreateTask(connection))
	// TODO: DELETE, PUT?

	router.HandleFunc("GET /task/{id}", handlers.GetTasksById(connection))
	router.HandleFunc("DELETE /task/{id}", handlers.DeleteTask(connection))

	// Prometheus metrics
	router.Handle("/metrics", promhttp.Handler())

	server := http.Server{
		Addr:    ":8080",
		Handler: middleware.PrometheusMiddleware(router),
	}

	log.Println("REST API listening on :8080")
	log.Fatal(server.ListenAndServe())

}
