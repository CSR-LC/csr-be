package main

import (
	"context"
	"database/sql"
	"log"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/restapi"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/restapi/operations"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/handlers"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/middlewares"
	"github.com/go-openapi/loads"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

func main() {
	var loggerConfig = zap.NewProductionConfig()
	loggerConfig.Level.SetLevel(zap.DebugLevel)

	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatalln(err)
	}

	connectionString := "host=localhost user=csr password=csr dbname=csr sslmode=disable"
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	// Create an ent.Driver from `db`.
	drv := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(drv))

	ctx := context.Background()

	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		logger.Error("error loading swagger spec", zap.Error(err))
		return
	}

	userHandler := handlers.NewUser(
		client,
		logger,
	)

	api := operations.NewBeAPI(swaggerSpec)
	api.UseSwaggerUI()
	api.BearerAuth = middlewares.BearerAuthenticateFunc("key", logger)

	api.UsersPostUserHandler = userHandler.PostUserFunc()
	api.UsersGetCurrentUserHandler = userHandler.GetUserFunc()
	api.UsersPatchUserHandler = userHandler.PatchUserFunc()

	server := restapi.NewServer(api)
	listeners := []string{"http"}

	server.EnabledListeners = listeners
	server.Host = "127.0.0.1"
	server.Port = 8080

	if err := server.Serve(); err != nil {
		logger.Error("server fatal error", zap.Error(err))
		return
	}

	if err := server.Shutdown(); err != nil {
		logger.Error("error shutting down server", zap.Error(err))
		return
	}
}
