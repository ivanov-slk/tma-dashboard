package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// RunSUTContainer creates a container of the dashboard service
// from its Dockerfile and starts it. It is a test helper using a
// shared context.
func RunSUTContainer(t testing.TB, ctx context.Context, port string) (*testcontainers.Container, func(), error) {
	t.Helper()

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:       "../",
			PrintBuildLog: true,
		},
		Networks: []string{
			"test-network",
		},
		NetworkAliases: map[string][]string{
			"test-network": {"dashboard"},
		},
		ExposedPorts: []string{fmt.Sprintf("%s:%s", port, port)},
		WaitingFor:   wait.ForListeningPort(nat.Port(port)).WithStartupTimeout(120 * time.Second),
	}

	sut, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("could not start dashboard container: %s", err)
	}

	cleanupFunc := func() {
		if err := sut.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}

	return &sut, cleanupFunc, nil
}
