package pkg

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fleimkeipa/lifery/model"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewPSQLClient() *pg.DB {
	// Determine if we are on local or cluster
	addr := "localhost:5432"
	if stage := os.Getenv("STAGE"); stage == "prod" {
		addr = "postgres:5432"
	}

	opts := pg.Options{
		Database: "case",
		User:     "postgres",
		Password: "password",
		Addr:     addr,
	}
	db := pg.Connect(&opts)

	if err := createSchema(db); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	return db
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*model.Event)(nil),
		(*model.Era)(nil),
		(*model.User)(nil),
		(*model.Connect)(nil),
	}

	for _, model := range models {
		opts := &orm.CreateTableOptions{
			IfNotExists: true,
		}

		if err := db.Model(model).CreateTable(opts); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

// GetTestInstance starts a PostgreSQL container for testing and returns a connected pg.DB client along with a cleanup function.
func GetTestInstance(ctx context.Context) (*pg.DB, func()) {
	const psqlVersion = "17.0"
	const port = "5432"

	req := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("postgres:%s", psqlVersion),
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", port)},
		WaitingFor:   wait.ForListeningPort(port), // Wait until the port is ready
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "test_db",
		},
		Cmd: []string{"postgres", "-c", "fsync=off"}, // Disable fsync for performance in tests
	}
	psqlClient, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("an error occurred while starting postgres container! error details: %v", err)
	}

	psqlPort, err := psqlClient.MappedPort(ctx, port)
	if err != nil {
		log.Fatalf("an error occurred while getting postgres port! error details: %v", err)
	}

	after, _ := strings.CutPrefix(psqlPort.Port(), "/")

	dbAddr := fmt.Sprintf("localhost:%s", after)
	opts := pg.Options{
		User:     "postgres",
		Password: "password",
		Database: "test_db",
		Addr:     dbAddr,
	}
	client := pg.Connect(&opts)

	if err := createTestTables(client); err != nil {
		log.Fatalf("Failed to create test schema: %v", err)
	}

	// Return the client and a cleanup function
	return client, func() {
		client.Close()
		psqlClient.Terminate(ctx)
	}
}

// createTestTables creates temporary test tables for the provided models.
func createTestTables(db *pg.DB) error {
	models := []interface{}{
		(*model.Event)(nil),
	}

	for _, model := range models {
		opts := orm.CreateTableOptions{
			Temp:        true, // Creates a temporary table for testing purposes.
			IfNotExists: true,
		}
		err := db.
			Model(model).
			CreateTable(&opts)
		if err != nil {
			return err
		}
	}

	return nil
}
