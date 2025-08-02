package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/HubertLipinski/go-rest-graphql-grpc/internal/config"
	"github.com/HubertLipinski/go-rest-graphql-grpc/internal/database"
	"github.com/HubertLipinski/go-rest-graphql-grpc/internal/graphql/generated"
	gql "github.com/HubertLipinski/go-rest-graphql-grpc/internal/graphql/resolver"
)

func main() {
	log.Print("Starting GraphQL")

	config.LoadEnv()

	connection, err := database.InitDBConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	executableSchema := generated.NewExecutableSchema(generated.Config{
		Resolvers: &gql.Resolver{DB: connection},
	})

	srv := handler.New(executableSchema)

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	srv.Use(extension.Introspection{})

	http.Handle("/", playground.Handler("GraphQL", "/query"))
	http.Handle("/query", srv)

	log.Println("GraphQL server listening on http://localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
