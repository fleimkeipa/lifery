package pkg

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fleimkeipa/lifery/model"

	"github.com/go-pg/pg/v10"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewPSQLClient() *pg.DB {
	host := "localhost"
	port := "5432"

	// Determine if we are on local or cluster
	stage := os.Getenv("STAGE")
	if stage == "prod" {
		host = os.Getenv("DB_HOST")
		port = os.Getenv("DB_PORT")
	}

	addr := ""
	if strings.Contains(host, ":") {
		addr = host
	} else {
		addr = fmt.Sprintf("%s:%s", host, port)
	}

	opts := pg.Options{
		Database: os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Addr:     addr,
	}

	if stage == model.StageProd && os.Getenv("DB_SSL") == "true" {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify:     true,
			SessionTicketsDisabled: true,
		}
	}

	db := pg.Connect(&opts)

	return db
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

	// Return the client and a cleanup function
	return client, func() {
		if err := client.Close(); err != nil {
			log.Printf("Error closing PostgreSQL client: %v", err)
		}
		if err := psqlClient.Terminate(ctx); err != nil {
			log.Printf("Error terminating PostgreSQL container: %v", err)
		}
	}
}
