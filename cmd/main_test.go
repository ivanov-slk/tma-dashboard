package main_test

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/ivanov-slk/tma-dashboard/adapters/httpserver"
	test "github.com/ivanov-slk/tma-dashboard/test"
	"github.com/ivanov-slk/tma-dashboard/test/specifications"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/testcontainers/testcontainers-go"
)

func TestDashboardServer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*120))
	defer cancel()

	networkName := "test-network"
	newNetwork, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{
			Name:           networkName,
			CheckDuplicate: true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		newNetwork.Remove(ctx)
	})

	natsContainer, natsCleanup, err := test.RunNATSContainer(t, ctx)
	if err != nil {
		t.Fatalf("could not initalize nats server: %s", err)
	}
	defer natsCleanup()

	// Publish a message to NATS
	nc, err := nats.Connect(natsContainer.URI)
	if err != nil {
		log.Fatalf("ERROR: failed to connect to nats server: %s", err)
	}
	defer nc.Close()

	js, _ := jetstream.New(nc)

	js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "TMA",
		Subjects: []string{"generated-data"},
	})

	js.Publish(ctx, "generated-data", []byte("{\"temperature\":15,\"humidity\":0.6,\"pressure\":1000,\"datetime\":\"2024-01-04T16:27:40Z\",\"id\":\"1\"}"))

	_, sutCleanup, err := test.RunSUTContainer(t, ctx, "1337", "http://nats-server:4222")
	if err != nil {
		t.Fatalf("could not initialize sut: %s", err)
	}
	defer sutCleanup()

	driver := httpserver.Driver{BaseURL: "http://localhost:1337", Client: &http.Client{Timeout: 120 * time.Second}}
	specifications.ExtractMetricsSpecification(t, &driver)
}
