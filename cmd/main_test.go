package main_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ivanov-slk/tma-dashboard/adapters/httpserver"
	test "github.com/ivanov-slk/tma-dashboard/test"
	"github.com/ivanov-slk/tma-dashboard/test/specifications"
)

func TestDashboardServer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*120))
	defer cancel()

	_, sutCleanup, err := test.RunSUTContainer(t, ctx, "1337")
	if err != nil {
		t.Fatalf("could not initialize sut: %s", err)
	}
	defer sutCleanup()

	// Run the dashboard as a container here; later add NATS to source data from.
	driver := httpserver.Driver{BaseURL: "http://localhost:1337", Client: &http.Client{Timeout: 3 * time.Second}}
	specifications.ExtractMetricsSpecification(t, &driver)
}
