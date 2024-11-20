package pkg

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewRedisClient() *redis.Client {
	addr := "localhost:6379"
	if stage := os.Getenv("STAGE"); stage == "prod" {
		addr = "redis:6379"
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return redisClient
}

// GetCacheTestInstance starts a Redis container for testing and returns a connected redis client along with a cleanup function.
func GetCacheTestInstance(ctx context.Context) (*redis.Client, func()) {
	const redisVersion = "7.4"
	const port = "6379"

	req := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("redis:%s", redisVersion),
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", port)},
		WaitingFor:   wait.ForListeningPort(port), // Wait until the port is ready
		Env:          map[string]string{},
		Cmd:          []string{"redis-server", "--appendonly", "no"}, // Disable appendonly for performance in tests
	}
	containerClient, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("an error occurred while starting redis container! error details: %v", err)
	}

	containerPort, err := containerClient.MappedPort(ctx, port)
	if err != nil {
		log.Fatalf("an error occurred while getting redis port! error details: %v", err)
	}

	after, _ := strings.CutPrefix(containerPort.Port(), "/")

	opts := redis.Options{
		Addr: fmt.Sprintf("redis:%s", after),
	}
	redisClient := redis.NewClient(&opts)

	// Return the client and a cleanup function
	return redisClient, func() {
		redisClient.Close()
		containerClient.Terminate(ctx)
	}
}
