package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestRedisClient(t *testing.T) {
	// Create a Redis container with a random port and wait for it to start
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	ctx := context.Background()
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start Redis container: %v", err)
	}
	defer redisC.Terminate(ctx)

	// Get the Redis container's host and port
	redisHost, err := redisC.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get Redis container's host: %v", err)
	}
	redisPort, err := redisC.MappedPort(ctx, "6379/tcp")
	if err != nil {
		t.Fatalf("Failed to get Redis container's port: %v", err)
	}

	// Create a Redis client and perform some operations
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort.Port()),
	})
	defer client.Close()

	err = client.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	val, err := client.Get(ctx, "key").Result()
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}

	if val != "value" {
		t.Errorf("Expected value %q, but got %q", "value", val)
	}
}
