package resolver

import "github.com/HubertLipinski/go-rest-graphql-grpc/internal/database"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB *database.DBConnection
}
